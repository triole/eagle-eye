package watcher

import (
	"testing"
	"time"
)

func TestRunCmd(t *testing.T) {
	w := Watcher{}

	tests := []struct {
		name    string
		cmdArr  []string
		pause   time.Duration
		verbose bool
		wantErr bool
	}{
		{
			name:    "successful command",
			cmdArr:  []string{"echo", "hello world"},
			pause:   0,
			verbose: false,
			wantErr: false,
		},
		{
			name:    "command with pause",
			cmdArr:  []string{"echo", "test"},
			pause:   1 * time.Millisecond,
			verbose: false,
			wantErr: false,
		},
		{
			name:    "non-existent command",
			cmdArr:  []string{"this_command_does_not_exist"},
			pause:   0,
			verbose: false,
			wantErr: true,
		},
		{
			name:    "command with error",
			cmdArr:  []string{"sh", "-c", "exit 1"},
			pause:   0,
			verbose: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := w.runCmd(tt.cmdArr, tt.pause, tt.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCmdOutputs(t *testing.T) {
	w := Watcher{}

	output, exitcode, err := w.runCmd([]string{"echo", "test output"}, 0, false)

	if err != nil {
		t.Errorf("runCmd() failed unexpectedly: %v", err)
	}

	if exitcode != 0 {
		t.Errorf("runCmd() exit code = %d, want 0", exitcode)
	}

	expected := "test output"
	if string(output) != expected+"\n" {
		t.Errorf("runCmd() output = %q, want %q", string(output), expected+"\n")
	}
}

func TestRunCmdNoOutput(t *testing.T) {
	w := Watcher{}

	output, exitcode, err := w.runCmd([]string{"true"}, 0, false)

	if err != nil {
		t.Errorf("runCmd() failed unexpectedly: %v", err)
	}

	if exitcode != 0 {
		t.Errorf("runCmd() exit code = %d, want 0", exitcode)
	}

	if string(output) != "" {
		t.Errorf("runCmd() output = %q, want empty string", string(output))
	}
}
