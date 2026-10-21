package watcher

import (
	"testing"
)

func TestWatcher_iterTemplate(t *testing.T) {
	w := Watcher{}

	// Create a test variable map
	varMap := make(map[string]tVarMapEntry)
	varMap["file"] = tVarMapEntry{
		Val:  "/path/to/test.txt",
		Desc: "file that triggered the event",
	}
	varMap["folder"] = tVarMapEntry{
		Val:  "/path/to/",
		Desc: "folder of the file that triggered the event",
	}

	tests := []struct {
		name   string
		arr    []string
		varMap map[string]tVarMapEntry
		want   []string
	}{
		{
			name:   "simple template with file variable",
			arr:    []string{"{{.file}}"},
			varMap: varMap,
			want:   []string{"/path/to/test.txt"},
		},
		{
			name:   "simple template with folder variable",
			arr:    []string{"{{.folder}}"},
			varMap: varMap,
			want:   []string{"/path/to/"},
		},
		{
			name:   "multiple templates with different variables",
			arr:    []string{"{{.file}}", "{{.folder}}"},
			varMap: varMap,
			want:   []string{"/path/to/test.txt", "/path/to/"},
		},
		{
			name:   "template with mixed content",
			arr:    []string{"File: {{.file}} in folder {{.folder}}"},
			varMap: varMap,
			want:   []string{"File: /path/to/test.txt in folder /path/to/"},
		},
		{
			name:   "empty array",
			arr:    []string{},
			varMap: varMap,
			want:   []string{},
		},
		{
			name:   "empty template string",
			arr:    []string{""},
			varMap: varMap,
			want:   []string{""},
		},
		{
			name:   "template with non-existent variable",
			arr:    []string{"{{.nonexistent}}"},
			varMap: varMap,
			want:   []string{"<no value>"},
		},
		{
			name:   "complex template with multiple variables",
			arr:    []string{"{{.file}} - {{.folder}} - {{.file}}"},
			varMap: varMap,
			want:   []string{"/path/to/test.txt - /path/to/ - /path/to/test.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := w.iterTemplate(tt.arr, tt.varMap)
			if len(result) != len(tt.want) {
				t.Errorf("iterTemplate() length = %d, want %d", len(result), len(tt.want))
				return
			}
			for i, v := range result {
				if v != tt.want[i] {
					t.Errorf("iterTemplate() = %q, want %q", v, tt.want[i])
				}
			}
		})
	}
}
