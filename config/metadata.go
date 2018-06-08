package config

import (
	"strings"
	"os"
	"path/filepath"
	. "GoVal/util"
)

const (
	RECURSIVE   string = "recursive"
	DESCRIPTION string = "addDesc"
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

type GoVal struct {
	Packages  int
	Files     int
	Functions int
	Internal  int
	Exported  int
	NoDocs    int
	WithDocs  int
}

type GoData struct {
	Package  string
	File     string
	Function string
	Start    int
	End      int
	Lines    int
}
