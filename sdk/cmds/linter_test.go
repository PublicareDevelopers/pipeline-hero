package cmds

import "testing"

func TestLinter(t *testing.T) {
	res, err := Linter("./...")
	if err != nil {
		t.Errorf("Error: %s\n", err)
	}

	t.Log(res)
}
