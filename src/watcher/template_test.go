package watcher

import (
	"testing"

	"github.com/radovskyb/watcher"
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

// Test with a more complex scenario that includes a real watcher event
func TestWatcher_iterTemplateWithRealEvent(t *testing.T) {
	w := Watcher{}

	// Create a mock event
	mockEvent := watcher.Event{
		Path: "/home/user/documents/test.txt",
	}

	// Create variable map using makeVarMap function
	varMap := w.makeVarMap(mockEvent)

	// Test templates that would be used in real scenarios
	testTemplates := []string{
		"{{.file}}",
		"{{.folder}}",
		"Processing file: {{.file}}",
		"File: {{.file}} in folder: {{.folder}}",
	}

	result := w.iterTemplate(testTemplates, varMap)

	// Expected results
	expected := []string{
		"/home/user/documents/test.txt",
		"/home/user/documents/",
		"Processing file: /home/user/documents/test.txt",
		"File: /home/user/documents/test.txt in folder: /home/user/documents/",
	}

	if len(result) != len(expected) {
		t.Errorf("iterTemplate() length = %d, want %d", len(result), len(expected))
		return
	}

	for i, v := range result {
		if v != expected[i] {
			t.Errorf("iterTemplate() = %q, want %q", v, expected[i])
		}
	}
}
