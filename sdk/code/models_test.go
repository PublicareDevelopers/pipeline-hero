package code

import (
	"reflect"
	"testing"
)

func TestAnalyser_SetCoverageByTotal(t *testing.T) {
	type args struct {
		totalText string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Test SetCoverageByTotal",
			args: args{
				totalText: "total: (statements) 100.0%",
			},
			want: 100.0,
		},
		{
			name: "Test SetCoverageByTotal",
			args: args{
				totalText: "total: (statements) 93.84%",
			},
			want: 93.84,
		},
		{
			name: "Test SetCoverageByTotal",
			args: args{
				totalText: "total: (statements) 0.0%",
			},
			want: 0.0,
		},
		{
			name: "Test SetCoverageByTotal",
			args: args{
				totalText: "total: nonsens",
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyser()
			a.SetCoverageByTotal(tt.args.totalText)
			if got := a.Coverage; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCoverageByTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}
