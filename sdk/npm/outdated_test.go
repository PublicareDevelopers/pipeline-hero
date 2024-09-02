package npm

import (
	"testing"
)

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

func TestOutDate_NewMajorVersionAvailable(t *testing.T) {
	type fields struct {
		Name      string
		Current   string
		Wanted    string
		Latest    string
		Dependent string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "major version available 1",
			/*
				"portal-vue": {
				    "current": "2.1.7",
				    "wanted": "2.1.7",
				    "latest": "3.0.0",
				    "dependent": "foobar"
				  },
			*/
			fields: fields{
				Name:    "portal-vue",
				Current: "2.1.7",
				Wanted:  "2.1.7",
				Latest:  "3.0.0",
			},
			want: true,
		},
		{
			name: "major version available 2",
			/*
				"vue": {
				    "current": "2.7.16",
				    "wanted": "2.7.16",
				    "latest": "3.4.38",
				    "dependent": "foobar"
				  },
			*/
			fields: fields{
				Name:    "vue",
				Current: "2.7.16",
				Wanted:  "2.7.16",
				Latest:  "3.4.38",
			},
			want: true,
		},
		{
			name: "major version not available 1",
			/*
				"vue-json-viewer": {
				    "current": "2.2.19",
				    "wanted": "2.2.22",
				    "latest": "2.2.22",
				    "dependent": "foobar"
				  },
			*/
			fields: fields{
				Name:    "vue-json-viewer",
				Current: "2.2.19",
				Wanted:  "2.2.22",
				Latest:  "2.2.22",
			},
			want: false,
		},
		{
			name: "major version not available 2",
			/*
				"vue-workflow-chart": {
				    "current": "0.4.5",
				    "wanted": "0.4.5",
				    "latest": "0.5.0",
				    "dependent": "foobar"
				  },
			*/
			fields: fields{
				Name:    "vue-workflow-chart",
				Current: "0.4.5",
				Wanted:  "0.4.5",
				Latest:  "0.5.0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOutDate(tt.fields.Name, tt.fields.Current, tt.fields.Wanted, tt.fields.Latest, tt.fields.Dependent)
			if err != nil {
				t.Errorf("NewOutDate() error = %v, want nil", err)
			}
			if got := o.NewMajorVersionAvailable(); got != tt.want {
				t.Errorf("NewMajorVersionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutDate_NewMinorVersionAvailable(t *testing.T) {
	type fields struct {
		Name      string
		Current   string
		Wanted    string
		Latest    string
		Dependent string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "minor version available 1",
			/*
				"vuex": {
					"current": "3.6.2",
					"wanted": "3.6.2",
					"latest": "4.1.0",
					"dependent": "foobar"
				},
			*/
			fields: fields{
				Name:    "vuex",
				Current: "3.6.2",
				Wanted:  "3.6.2",
				Latest:  "4.1.0",
			},
			want: false,
		},
		{
			name: "minor version available 2",
			/*
				"vuejs-storage": {
					"current": "3.1.0",
					"wanted": "3.1.1",
					"latest": "3.2.1",
					"dependent": "foobar"
				},
			*/
			fields: fields{
				Name:    "vuejs-storage",
				Current: "3.1.0",
				Wanted:  "3.1.1",
				Latest:  "3.2.1",
			},
			want: true,
		},
		{
			name: "minor version not available 1",
			/*
				"vue-workflow-chart": {
					"current": "0.4.5",
					"wanted": "0.4.5",
					"latest": "0.5.0",
					"dependent": "foobar"
				},
			*/
			fields: fields{
				Name:    "vue-workflow-chart",
				Current: "0.4.5",
				Wanted:  "0.4.5",
				Latest:  "0.4.8",
			},
			want: false,
		},
		{
			name: "minor version not available 2",
			/*
				"vue-json-viewer": {
					"current": "2.2.19",
					"wanted": "2.2.22",
					"latest": "2.2.22",
					"dependent": "foobar"
				},
			*/
			fields: fields{
				Name:    "vue-json-viewer",
				Current: "2.2.19",
				Wanted:  "2.2.22",
				Latest:  "2.2.22",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOutDate(tt.fields.Name, tt.fields.Current, tt.fields.Wanted, tt.fields.Latest, tt.fields.Dependent)
			if err != nil {
				t.Errorf("NewOutDate() error = %v, want nil", err)
			}
			if got := o.NewMinorVersionAvailable(); got != tt.want {
				t.Errorf("NewMinorVersionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutDate_NewPatchVersionAvailable(t *testing.T) {
	type fields struct {
		Name      string
		Current   string
		Wanted    string
		Latest    string
		Dependent string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "patch version available 1",
			/*
			   "vuex": {
			    "current": "3.6.2",
			    "wanted": "3.6.2",
			    "latest": "3.6.3",
			    "dependent": "foobar"
			   },
			*/
			fields: fields{
				Name:    "vuex",
				Current: "3.6.2",
				Wanted:  "3.6.2",
				Latest:  "3.6.3",
			},
			want: true,
		},
		{
			name: "patch version available 2",
			/*
			   "vuejs-storage": {
			    "current": "3.1.0",
			    "wanted": "3.1.0",
			    "latest": "3.1.1",
			    "dependent": "foobar"
			   },
			*/
			fields: fields{
				Name:    "vuejs-storage",
				Current: "3.1.0",
				Wanted:  "3.1.0",
				Latest:  "3.1.1",
			},
			want: true,
		},
		{
			name: "patch version not available 1",
			/*
			   "vue-workflow-chart": {
			    "current": "0.4.5",
			    "wanted": "0.4.5",
			    "latest": "0.4.5",
			    "dependent": "foobar"
			   },
			*/
			fields: fields{
				Name:    "vue-workflow-chart",
				Current: "0.4.5",
				Wanted:  "0.4.5",
				Latest:  "0.4.5",
			},
			want: false,
		},
		{
			name: "patch version not available 2",
			/*
			   "vue-json-viewer": {
			    "current": "2.2.19",
			    "wanted": "2.2.19",
			    "latest": "2.2.19",
			    "dependent": "foobar"
			   },
			*/
			fields: fields{
				Name:    "vue-json-viewer",
				Current: "2.2.19",
				Wanted:  "2.2.19",
				Latest:  "3.0.0",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := NewOutDate(tt.fields.Name, tt.fields.Current, tt.fields.Wanted, tt.fields.Latest, tt.fields.Dependent)
			if err != nil {
				t.Errorf("NewOutDate() error = %v, want nil", err)
			}
			if got := o.NewPatchVersionAvailable(); got != tt.want {
				t.Errorf("NewPatchVersionAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
