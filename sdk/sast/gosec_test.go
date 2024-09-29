package sast

import (
	"encoding/json"
	"os"
	"testing"
)

func TestAnalyseSAST(t *testing.T) {
	sastRes, err := os.ReadFile("gosec.json")
	if err != nil {
		t.Fatal(err)
	}

	sastStruct := SAST{}
	err = json.Unmarshal(sastRes, &sastStruct)
	if err != nil {
		t.Fatal(err)
	}

	res, err := HasFailedSAST(sastStruct)
	if err != nil {
		t.Fatal(err)
	}

	if !res {
		t.Fatal("expected to have failed sast")
	}
}
