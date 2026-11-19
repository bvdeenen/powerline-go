package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func symlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

func Test_pythonVersion(t *testing.T) {
	root := t.TempDir()

	// System python with headers: /usr/bin/python3 -> python3.14
	writeFile(t, filepath.Join(root, "usr/bin/python3.14"), "")
	symlink(t, "python3.14", filepath.Join(root, "usr/bin/python3"))
	writeFile(t, filepath.Join(root, "usr/include/python3.14/patchlevel.h"),
		"#define PY_MINOR_VERSION 14\n#define PY_VERSION              \"3.14.8\"\n")

	// System python without headers
	writeFile(t, filepath.Join(root, "opt/bin/python3.12"), "")
	symlink(t, "python3.12", filepath.Join(root, "opt/bin/python3"))

	// venv module: version key
	symlink(t, filepath.Join(root, "usr/bin/python3"), filepath.Join(root, "venv/bin/python3"))
	writeFile(t, filepath.Join(root, "venv/pyvenv.cfg"), "home = /usr/bin\nversion = 3.14.7\n")

	// virtualenv: version_info key with release level
	symlink(t, filepath.Join(root, "usr/bin/python3"), filepath.Join(root, "venv2/bin/python3"))
	writeFile(t, filepath.Join(root, "venv2/pyvenv.cfg"), "home = /usr/bin\nversion_info = 3.13.1.final.0\n")

	// Plain python3 binary, single header set under its prefix
	writeFile(t, filepath.Join(root, "plain/bin/python3"), "")
	writeFile(t, filepath.Join(root, "plain/include/python3.11/patchlevel.h"), "#define PY_VERSION \"3.11.9\"\n")

	// Plain python3 binary, nothing to go on
	writeFile(t, filepath.Join(root, "bare/bin/python3"), "")

	tests := []struct {
		exe  string
		want string
	}{
		{"usr/bin/python3", "3.14.8"},
		{"opt/bin/python3", "3.12"},
		{"venv/bin/python3", "3.14.7"},
		{"venv2/bin/python3", "3.13.1"},
		{"plain/bin/python3", "3.11.9"},
		{"bare/bin/python3", ""},
		{"missing/bin/python3", ""},
	}
	for _, tt := range tests {
		t.Run(tt.exe, func(t *testing.T) {
			if got := pythonVersion(filepath.Join(root, tt.exe)); got != tt.want {
				t.Errorf("pythonVersion(%q) = %q, want %q", tt.exe, got, tt.want)
			}
		})
	}
}
