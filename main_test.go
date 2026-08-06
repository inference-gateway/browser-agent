package main

import "testing"

func TestTaskWorkers(t *testing.T) {
	tests := []struct {
		name string
		env  string
		want int
	}{
		{"default", "", 4},
		{"explicit", "2", 2},
		{"invalid", "abc", 4},
		{"zero falls back", "0", 4},
		{"negative falls back", "-3", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TASK_WORKERS", tt.env)
			if got := taskWorkers(); got != tt.want {
				t.Fatalf("taskWorkers() with TASK_WORKERS=%q = %d, want %d", tt.env, got, tt.want)
			}
		})
	}
}
