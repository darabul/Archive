package process

import (
	"flag"
	archiver "util/pkg"
)

const archiveFileName = "archive.txt"

func ArchiveProcesses(mode string, sourcePath string, resultPath string) {
	switch mode {
	case "arch":
		a := archiver.New(archiveFileName)
		a.Archive(sourcePath, resultPath)
	case "unarch":
		a := archiver.New(archiveFileName)
		a.Unarchive(sourcePath, resultPath)
	}
}

func ParseFlags(mode *string, sourcePath *string, resultPath *string) {
	flag.StringVar(mode, "m", "", "режим работы программы")
	flag.StringVar(sourcePath, "s", "", "путь к папке, с которой мы будем проводить манипуляции")
	flag.StringVar(resultPath, "r", "", "путь, в котором будем лежать результат программы")
	flag.Parse()
}
