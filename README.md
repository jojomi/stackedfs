# stackedfs

`stackedfs` allows for overlaying different sources to read a file.
Think of `OverlayFS` for `Go`.
It works great with the embed directive and requires at least Go 1.16 because of the `fs.FS` interface it makes heavy use of.

## Install

``` bash
go get -u github.com/jojomi/stackedfs
```

## Example usage

``` go
package main

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/rs/zerolog/log"
	"gitlab.com/jojomi/zeiterfassung/stackedfs"
)

//go:embed templates/*
var embeddedTemplates embed.FS

func getTemplateReader() (*stackedfs.StackedFS, error) {
	result := stackedfs.NewStackedFS()

	// add the paths for overwriting the embedded resource files
	paths := getLocalOverwritePaths()
	for _, p := range paths {
		result.AddFS(os.DirFS(filepath.Join(p, "templates")))
	}

	// use the embedded variant as default if nothing else has matched
	subFS, err := fs.Sub(embeddedTemplates, "templates")
	if err != nil {
		return nil, err
	}
	result.AddFS(subFS)

	return result, nil
}

func getLocalOverwritePaths() []string {
	// override paths
	// 1. WORKING_DIR
	// 2. ~/BINARY_NAME/
	// 3. BINARY_DIR

	var (
		binaryName = "stackfs-test"
		folders []string
	)

	workingPath, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("working dir not found")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("home dir not found")
	}
	homePath := filepath.Join(home, "."+binaryName)

	binaryPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal().Err(err).Msg("binary path not found")
	}

	return append(folders, workingPath, homePath, binaryPath)
}
```
