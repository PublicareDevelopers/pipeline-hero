package helper

import "testing"

func TestConvertVersion(t *testing.T) {
	version, err := ConvertVersion("1.2.3")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 1 {
		t.Errorf("ConvertVersion() Major = %v, want 1", version.Major)
	}

	if version.Minor != 2 {
		t.Errorf("ConvertVersion() Minor = %v, want 2", version.Minor)
	}

	if version.Patch != 3 {
		t.Errorf("ConvertVersion() Patch = %v, want 3", version.Patch)
	}

	version, err = ConvertVersion("1.2")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 1 {
		t.Errorf("ConvertVersion() Major = %v, want 1", version.Major)
	}

	if version.Minor != 2 {
		t.Errorf("ConvertVersion() Minor = %v, want 2", version.Minor)
	}

	if version.Patch != 0 {
		t.Errorf("ConvertVersion() Patch = %v, want 0", version.Patch)
	}

	version, err = ConvertVersion("1")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 1 {
		t.Errorf("ConvertVersion() Major = %v, want 1", version.Major)
	}

	if version.Minor != 0 {
		t.Errorf("ConvertVersion() Minor = %v, want 0", version.Minor)
	}

	if version.Patch != 0 {
		t.Errorf("ConvertVersion() Patch = %v, want 0", version.Patch)
	}

	version, err = ConvertVersion("8.4.0-rc2")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 8 {
		t.Errorf("ConvertVersion() Major = %v, want 8", version.Major)
	}

	if version.Minor != 4 {
		t.Errorf("ConvertVersion() Minor = %v, want 4", version.Minor)
	}

	if version.Patch != 0 {
		t.Errorf("ConvertVersion() Patch = %v, want 0", version.Patch)
	}

	_, err = ConvertVersion("foobar")
	if err == nil {
		t.Errorf("ConvertVersion() error = %v, want not nil", err)
	}

	_, err = ConvertVersion("1.2.3.4")
	if err == nil {
		t.Errorf("ConvertVersion() error = %v, want not nil", err)
	}

	version, err = ConvertVersion("")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 0 {
		t.Errorf("ConvertVersion() Major = %v, want 0", version.Major)
	}

	if version.Minor != 0 {
		t.Errorf("ConvertVersion() Minor = %v, want 0", version.Minor)
	}

	if version.Patch != 0 {
		t.Errorf("ConvertVersion() Patch = %v, want 0", version.Patch)
	}

	version, err = ConvertVersion("v1.2.3")
	if err != nil {
		t.Errorf("ConvertVersion() error = %v, want nil", err)
	}

	if version.Major != 1 {
		t.Errorf("ConvertVersion() Major = %v, want 1", version.Major)
	}

	if version.Minor != 2 {
		t.Errorf("ConvertVersion() Minor = %v, want 2", version.Minor)
	}

	if version.Patch != 3 {
		t.Errorf("ConvertVersion() Patch = %v, want 3", version.Patch)
	}
}
