package main

import (
	process "util/internal"
)

func main() {
	var mode string
	var sourcePath string
	var resultPath string

	process.ParseFlags(&mode, &sourcePath, &resultPath)
	process.ArchiveProcesses(mode, sourcePath, resultPath)

}
