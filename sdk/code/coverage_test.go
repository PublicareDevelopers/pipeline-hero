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

func TestAnalyser_GetCoverageInterpretation(t *testing.T) {
	type fields struct {
		Coverage float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 100.0,
			},
			want: "coverage is good, have 100.00  percent\n",
		},
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 93.84,
			},
			want: "coverage is good, have 93.84  percent\n",
		},
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 75.0,
			},
			want: "coverage is good, have 75.00  percent\n",
		},
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 74.99,
			},
			want: "coverage is ok, have 74.99  percent\n",
		},
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 50.0,
			},
			want: "coverage is ok, have 50.00  percent\n",
		},
		{
			name: "Test GetCoverageInterpretation",
			fields: fields{
				Coverage: 49.99,
			},
			want: "coverage is BAD, have 49.99  percent\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyser()
			a.Coverage = tt.fields.Coverage

			if got := a.GetCoverageInterpretation(); got != tt.want {
				t.Errorf("GetCoverageInterpretation() = %v, want %v", got, tt.want)
			}
		})
	}
}
