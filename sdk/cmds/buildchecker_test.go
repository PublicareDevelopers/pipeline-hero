package cmds

import "testing"

func TestCheckBuild(t *testing.T) {
	_, err := CheckBuild("./")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestCheckLocalBuild(t *testing.T) {
	_, err := CheckLocalBuild("./")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestCheckYml(t *testing.T) {
	_, err := CheckYml("./")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
