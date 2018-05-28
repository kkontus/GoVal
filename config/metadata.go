package config

import (
	"strings"
	"os"
	"path/filepath"
	. "GoVal/util"
)

var skipFiles = []string{"metadata.go", "util.go"}

func Filter(info os.FileInfo) bool {
	name := info.Name()

	if info.IsDir() {
		return false
	}

	if Contains(skipFiles, name) {
		return false
	}

	if filepath.Ext(name) != ".go" {
		return false
	}

	if strings.HasSuffix(name, "_test.go") {
		return false
	}

	return true
}
