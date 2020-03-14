package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	filepath.Walk("./testdata", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") {
			cmd := exec.Command("bash", filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			_, err := cmd.Output()
			if err != nil {
				t.Errorf("FAIL: command execution")
			}

			outFile := strings.TrimSuffix(path, filepath.Ext(path)) + ".html"
			output, err := ioutil.ReadFile(outFile)
			if err != nil {
				t.Errorf("FAIL: Reading on output file: %s\n", outFile)
			}
			expFile := strings.TrimSuffix(path, filepath.Ext(path)) + "_exp.html"
			expected, err := ioutil.ReadFile(expFile)
			if err != nil {
				t.Errorf("FAIL: Reading on expected file: %s\n", expFile)
			}

			if strings.HasPrefix(string(output), strings.TrimSuffix(string(expected), "\n")) {
				t.Logf("PASS: %s\n", path)
			} else {
				t.Errorf("FAIL: output differs: %s\n", path)
			}

			if err = os.Remove(outFile); err != nil {
				t.Logf("Unremoved %s file\n", outFile)
			}
		}
		return nil
	})
}
