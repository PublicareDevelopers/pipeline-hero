package code

import (
	"reflect"
	"testing"
)

func TestAnalyser_SetThreshold(t *testing.T) {
	type args struct {
		threshold float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test SetThreshold",
			args: args{
				threshold: 75.0,
			},
			want: 75.0,
		},
		{
			name: "Test SetThreshold",
			args: args{
				threshold: 50.0,
			},
			want: 50.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyser().SetThreshold(tt.args.threshold)

			if got := a.Threshold; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetThreshold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalyser_SetVulnCheck(t *testing.T) {
	a := NewAnalyser().SetVulnCheck("trivy")

	if a.VulnCheck != "trivy" {
		t.Errorf("SetVulnCheck() = %v, want %v", a.VulnCheck, "trivy")
	}
}

func TestAnalyser_SetGoVersion(t *testing.T) {
	a := NewAnalyser().SetGoVersion("1.14.1")

	if a.GoVersion != "1.14.1" {
		t.Errorf("SetGoVersion() = %v, want %v", a.GoVersion, "1.14.1")
	}
}

func TestAnalyser_PushError(t *testing.T) {
	a := NewAnalyser()

	if len(a.errors) != 0 {
		t.Errorf("PushError() = %v, want %v", len(a.errors), 0)
	}

	a.PushError("error 1")

	if len(a.errors) != 1 {
		t.Errorf("PushError() = %v, want %v", len(a.errors), 1)
	}

	errors := a.GetErrors()

	if len(errors) != 1 {
		t.Errorf("GetErrors() = %v, want %v", len(errors), 1)
	}

	if errors[0] != "error 1" {
		t.Errorf("GetErrors() = %v, want %v", errors[0], "error 1")
	}
}

func TestAnalyser_PushWarning(t *testing.T) {
	a := NewAnalyser()

	if len(a.warnings) != 0 {
		t.Errorf("PushWarning() = %v, want %v", len(a.warnings), 0)
	}

	a.PushWarning("warning 1")

	if len(a.warnings) != 1 {
		t.Errorf("PushWarning() = %v, want %v", len(a.warnings), 1)
	}

	warnings := a.GetWarnings()

	if len(warnings) != 1 {
		t.Errorf("GetWarnings() = %v, want %v", len(warnings), 1)
	}

	if warnings[0] != "warning 1" {
		t.Errorf("GetWarnings() = %v, want %v", warnings[0], "warning 1")
	}
}
