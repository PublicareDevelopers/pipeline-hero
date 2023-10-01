package code

import "testing"

func TestAnalyser_CheckThreshold(t *testing.T) {
	type fields struct {
		Threshold float64
		Coverage  float64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test CheckThreshold",
			fields: fields{
				Threshold: 75.0,
				Coverage:  75.0,
			},
			wantErr: false,
		},
		{
			name: "Test CheckThreshold",
			fields: fields{
				Threshold: 75.0,
				Coverage:  74.9,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAnalyser()
			a.SetThreshold(tt.fields.Threshold)
			a.Coverage = tt.fields.Coverage

			if err := a.CheckThreshold(); (err != nil) != tt.wantErr {
				t.Errorf("CheckThreshold() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
