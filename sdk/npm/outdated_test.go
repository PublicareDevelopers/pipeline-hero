package npm

import (
	"testing"
)

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

	_, err = ConvertVersion("foobar")
	if err == nil {
		t.Errorf("ConvertVersion() error = %v, want not nil", err)
	}

	_, err = ConvertVersion("1.2.3.4")
	if err == nil {
		t.Errorf("ConvertVersion() error = %v, want not nil", err)
	}
}

func TestNewOutDate(t *testing.T) {
	outdate, err := NewOutDate("@codemirror/autocomplete", "0.18.8", "0.18.8", "6.18.0", "foobar")
	if err != nil {
		t.Errorf("NewOutDate() error = %v, want nil", err)
	}

	if outdate.Name != "@codemirror/autocomplete" {
		t.Errorf("NewOutDate() Name = %v, want @codemirror/autocomplete", outdate.Name)
	}

	if outdate.CurrentVersion.Major != 0 {
		t.Errorf("NewOutDate() CurrentVersion.Major = %v, want 0", outdate.CurrentVersion.Major)
	}

	if outdate.CurrentVersion.Minor != 18 {
		t.Errorf("NewOutDate() CurrentVersion.Minor = %v, want 18", outdate.CurrentVersion.Minor)
	}

	if outdate.CurrentVersion.Patch != 8 {
		t.Errorf("NewOutDate() CurrentVersion.Patch = %v, want 8", outdate.CurrentVersion.Patch)
	}

	_, err = NewOutDate("@codemirror/autocomplete", "not-parseable", "0.18.8", "6.18.0", "foobar")
	if err == nil {
		t.Errorf("NewOutDate() error = %v, want not nil", err)
	}
}
