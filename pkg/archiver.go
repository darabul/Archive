package archiver

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	TempFile = "temp.txt"
	BufSize  = 4096
)

type Archiver struct {
	archive string
}

func New(archive string) *Archiver {
	return &Archiver{archive: archive}
}

func (a *Archiver) Archive(inputDirName string, outputDirName string) error {
	dir, err := os.OpenFile(inputDirName, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer dir.Close()

	temp, err := os.OpenFile(TempFile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(TempFile)
	}()

	var headers strings.Builder
	err = a.archiveRecursive(dir, temp, &headers, path.Base(inputDirName))
	if err != nil {
		return err
	}

	resultPath := path.Join(outputDirName, a.archive)
	out, err := os.OpenFile(resultPath, os.O_WRONLY|os.O_CREATE, 0622)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(headers.String() + "\n")
	if err != nil {
		return err
	}

	buf := make([]byte, BufSize)
	temp.Seek(0, 0)
	for {
		n, err := temp.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("%w", err)
		}
		_, err = out.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Archiver) archiveRecursive(dir, out *os.File, headers *strings.Builder, curDir string) error {
	files, err := dir.ReadDir(0)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		headers.WriteString(fmt.Sprintf("%s/,%d;", curDir, 0))
	}

	for _, file := range files {
		if file.Name() == out.Name() {
			continue
		}
		filePath := path.Join(dir.Name(), file.Name())
		fileName := path.Base(filePath)
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if file.IsDir() {
			err = a.archiveRecursive(f, out, headers, path.Join(curDir, fileName))
			if err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			headers.WriteString(fmt.Sprintf("%s,%d;", path.Join(curDir, fileName), len(data)))
			_, err = out.Write(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Archiver) Unarchive(inputDirName string, outputDirName string) error {
	sourcePath := path.Join(inputDirName, a.archive)
	arch, err := os.OpenFile(sourcePath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer arch.Close()

	scanner := bufio.NewScanner(arch)
	scanner.Scan()
	header := scanner.Text()
	arch.Seek(int64(len(header)+1), 0)

	filesInfo := strings.Split(header[:len(header)-1], ";")
	for _, fi := range filesInfo {
		fileInfo := strings.Split(fi, ",")
		filePath := fileInfo[0]
		fileSize, err := strconv.Atoi(fileInfo[1])
		if err != nil {
			return err
		}
		fileData := make([]byte, fileSize)
		absoluteFilePath := path.Join(outputDirName, filePath)

		var isDir bool
		if len(strings.Split(filePath, ".")) == 1 && fileSize == 0 {
			isDir = true
			err := os.MkdirAll(absoluteFilePath, 0777)
			if err != nil {
				return err
			}
		}

		_, err = arch.Read(fileData)
		if err != nil {
			return err
		}

		if _, err = os.Stat(absoluteFilePath); os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(absoluteFilePath), 0777)
			if err != nil {
				return err
			}
		}

		if !isDir {
			err = os.WriteFile(absoluteFilePath, fileData, 0777)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
