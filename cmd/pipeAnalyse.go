package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/cmds"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/code"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/sast"
	"github.com/PublicareDevelopers/pipeline-hero/sdk/slack"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"sync"
)

const BAD = 50.0
const MEDIUM = 75.0

var testSetup string
var envVariables map[string]string

// pipeAnalyseCmd represents the pipeAnalyse command
var pipeAnalyseCmd = &cobra.Command{
	Use:   "analyse",
	Short: "running a pipeline-hero analyse for go; will exit with an error if things are not as expected",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for key, value := range envVariables {
			err := os.Setenv(key, value)
			if err != nil {
				color.Red("Error setting env vars: %s\n", err)
				os.Exit(255)
			}

			fmt.Println("set env variable", key, "to", value)
		}

		analyser := code.NewAnalyser().SetThreshold(coverageThreshold)

		var wg sync.WaitGroup

		wg.Add(1)
		go analyseVersion(analyser, &wg)

		wg.Add(1)
		go analyseSDependencyGraph(analyser, &wg)

		wg.Add(1)
		go analyseTestCoverage(analyser, &wg)

		wg.Add(1)
		go analyseVulnCheck(analyser, &wg)

		wg.Add(1)
		go analyseSASTCheck(analyser, &wg)

		wg.Wait()

		if !useSlack {
			if len(analyser.GetErrors()) > 0 {
				color.Red("pipeline-hero failed\n")

				for _, err := range analyser.GetErrors() {
					color.Red("%s\n", err)
				}

				os.Exit(255)
			}
			return
		}

		client, err := slack.NewClient()
		if err != nil {
			color.Red("error at slack handling: %s\n", err)
			os.Exit(255)
		}

		err = client.StartConversation(analyser, "go")
		if err != nil {
			color.Red("error starting slack conversation: %s\n", err)
			os.Exit(255)
		}

		if len(analyser.GetErrors()) > 0 {
			color.Red("pipeline-hero failed\n")
			os.Exit(255)
		}
	},
}

func init() {
	pipeCmd.AddCommand(pipeAnalyseCmd)
	pipeAnalyseCmd.Flags().StringVarP(&testSetup, "test-setup", "t", "./...", "Test setup to use")
	pipeAnalyseCmd.Flags().BoolVarP(&useSlack, "slack", "s", false, "Send results to slack")
	pipeAnalyseCmd.Flags().Float64VarP(&coverageThreshold, "coverage-threshold", "c", 75.0, "Coverage threshold to use")
	pipeAnalyseCmd.Flags().StringToStringVarP(&envVariables, "env", "e", map[string]string{}, "Environment variables to set")
}

func analyseVersion(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	version, err := cmds.GetVersion()
	if err != nil {
		color.Red("error getting go version: %s\n", err)
		analyser.PushError(fmt.Sprintf("cannot find the go version: %s\n", err))
		return
	}

	analyser.SetGoVersion(version)

	color.Green("%s\n", version)
}

func analyseSDependencyGraph(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	analyser.SetModule()
	analyser.SetUpdatableRequirements()

	dependencyGraph, err := cmds.GetDependencyGraph()
	if err != nil {
		analyser.PushWarning(fmt.Sprintf("internal pipeline-hero error: cannot find the dependency graph: %s\n", err))
		color.Red("error getting dependency graph: %s\n", err)
		return
	}

	color.Green("setting dependency graph\n")

	analyser.SetDependencyGraph(dependencyGraph)
}

func analyseTestCoverage(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("\nrunning tests\n")
	out, err := exec.Command("go", "test", testSetup, "-v", fmt.Sprintf("-coverpkg=%s", testSetup), "-coverprofile=cover.cov").Output()
	if err != nil {
		color.Red("error running tests: %s\n", err)
		color.Red("Tests failed:\n%s\n", string(out))

		//regex find`FAIL(.*)`gm
		fails := ""

		reg := regexp.MustCompile(`FAIL(.*)`)
		matches := reg.FindAllString(string(out), -1)
		for _, match := range matches {
			fails += match + "\n"
		}

		analyser.TestResult = fails

		analyser.PushError(fmt.Sprintf("Tests failed:%s \nThe parsed results can be found in the thread or have a look at the pipeline for the complete output\n", err))
		return
	}

	coverProfile := string(out)
	analyser.SetCoverProfile(coverProfile)

	fmt.Println(coverProfile)
	color.Green("running coverage\n")

	//go tool cover -func=cover.cov | tee coverage.txt
	out, err = exec.Command("go", "tool", "cover", "-func=cover.cov").Output()
	if err != nil {
		fmt.Printf("%s\n", string(out))
		color.Red("error getting coverage: %s\n", err)
		analyser.PushError(fmt.Sprintf("Coverage failed:\n%s\n", string(out)))
		return
	}

	//we have something like total:  (statements)    0.0%
	//grep the total amount with a regex
	totalText := string(out)

	analyser.SetCoverageByTotal(totalText)

	fmt.Println(analyser.GetCoverageInterpretation())

	err = analyser.CheckThreshold()
	if err != nil {
		color.Red("%s\n", err)
		analyser.PushError(fmt.Sprintf("coverage threshold not met: have  %.2f  percent", analyser.Coverage))
		return
	}
}

func analyseVulnCheck(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("starting vuln check\n")

	var vulCheck string
	vulCheck, err := cmds.VulnCheck(testSetup)
	if err != nil {
		color.Red("vuln check error:\n%s\n", err)
		analyser.HasVuln = true
		analyser.PushError(err.Error())
		resp, err := sendVulnToPlatform(err.Error())
		if err != nil {
			color.Red("error sending vuln to platform: %s\n", err)
			return
		}

		color.Green(fmt.Sprintf("vuln sent to platform: %+v\n", resp))

		return
	}

	fmt.Println(vulCheck)
	analyser.SetVulnCheck(vulCheck)
}

func analyseSASTCheck(analyser *code.Analyser, wg *sync.WaitGroup) {
	defer wg.Done()

	color.Green("starting SAST check\n")

	sastStruct := sast.SAST{}
	sastCheckJson, err := cmds.GoSAST()
	if err != nil {
		errString := err.Error()
		analyser.HasSASTCheckFail = true
		color.Red("SAST check failed\n")

		err = json.Unmarshal([]byte(errString), &sastStruct)
		if err != nil {
			color.Red("error unmarshalling SAST: %s\n", err)

			analyser.PushError(fmt.Sprintf("cannot parse the gosec.json: %s\n", err))
		}

		sastCheckString := fmt.Sprintf("Found %d SAST issues\n", sastStruct.Stats.Found)
		for _, issue := range sastStruct.Issues {
			sastCheckString += fmt.Sprintf("- %s (CWE %s) at %s line %s\n Confidence: %s; Severity: %s; \n\n", issue.Details, issue.Cwe.ID, issue.File, issue.Line, issue.Confidence, issue.Severity)
		}

		analyser.SetSASTCheck(sastCheckJson)
		analyser.PushError(sastCheckString)

		color.Red(sastCheckString)

		resp, err := sendSastToPlatform(errString)
		if err != nil {
			color.Red("error sending SAST to platform: %s\n", err)
			return
		}

		color.Green(fmt.Sprintf("SAST sent to platform: %+v\n", resp))

		return
	}

	color.Green("SAST check passed\n")
}
