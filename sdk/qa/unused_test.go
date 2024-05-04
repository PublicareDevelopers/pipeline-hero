package qa

import "testing"

func TestExtractZipPath(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				line: "zip -j bin/workflow/cluster_export.zip bin/workflow/cluster/export/bootstrap && rm -rf bin/workflow/cluster/export",
			},
			want:    "bin/workflow/cluster_export.zip",
			wantErr: false,
		},
		{
			name: "no zip",
			args: args{
				line: "rm -rf bin/workflow/cluster/export",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractZipPath(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractZipPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractZipPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
