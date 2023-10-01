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
