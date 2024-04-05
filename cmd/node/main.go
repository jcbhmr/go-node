package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

const version = "20.12.1"

func try(err error) {
	if err != nil {
		panic(err)
	}
}

func try1[A any](a A, err error) A {
	if err != nil {
		panic(err)
	}
	return a
}

func catch[T any](rerr *T, f func(err T)) {
	if err := recover(); err != nil {
		*rerr = err.(T)
		if f != nil {
			f(*rerr)
		}
	}
}

func main() {
	// 1. Determine which Node.js distribution archive to download.
	var archiveUrl string
	switch fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH) {
	case "darwin-arm64": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-arm64.tar.gz", version)
	case "darwin-x64": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-x64.tar.gz", version)
	case "linux-arm64": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-arm64.tar.gz", version)
	case "linux-armv7l": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-armv7l.tar.gz", version)
	case "linux-ppc64le": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-ppc64le.tar.gz", version)
	case "linux-s390x": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-s390x.tar.gz", version)
	case "linux-x64": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-x64.tar.gz", version)
	case "windows-x64": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x64.zip", version)
	case "windows-x86": archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x86.zip", version)
	default: panic(fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH))
	}

	// 2. Download the Node.js distribution archive.
	archiveFilename := path.Base(try1(url.Parse(archiveUrl)).Path)
	tempDirPath := try1(os.MkdirTemp("", ""))
	defer os.RemoveAll(tempDirPath)
	archivePath := filepath.Join(tempDirPath, archiveFilename)
	try(downloadFile(archiveUrl, archivePath))

	// 3. Extract the Node.js distribution archive right next to the executable.
	exePath := try1(os.Executable())
	goNodeDirPath := filepath.Join(filepath.Dir(exePath), ".go-node")
	switch path.Ext(try1(url.Parse(archiveUrl)).Path) {
	case ".tar.gz": extractTarGz(archivePath, goNodeDirPath)
	case ".zip": extractZip(archivePath, goNodeDirPath)
	}

	// 4. Replace this executable with a symlink to the Node.js executable.
	nodePath := filepath.Join(goNodeDirPath, "bin", "node")
}

func downloadFile(srcUrl string, destPath string) (rerr error) {
	defer catch(&rerr, nil)
	file := try1(os.Create(destPath))
	defer file.Close()
	res := try1(http.Get(srcUrl))
	defer res.Body.Close()
	try1(io.Copy(file, res.Body))
	return nil
}

func extractTarGz(archivePath string, destDirPath string) (rerr error) {
	defer catch(&rerr, nil)
	file := try1(os.Open(archivePath))
	defer file.Close()
	reader := try1(tar.NewReader(file))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		try(err)
		destPath := filepath.Join(destDirPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir: try(os.MkdirAll(destPath, 0755))
		case tar.TypeReg: try(os.WriteFile(destPath, try1(io.ReadAll(reader))))
		}
	}
	return nil
}

func extractZip(archivePath string, destDirPath string) (rerr error) {
	defer catch(&rerr, nil)
	reader := try1(zip.NewReader(try1(os.Open(archivePath)), try1(os.Stat(archivePath)).Size()))
	for _, file := range reader.File {
		destPath := filepath.Join(destDirPath, file.Name)
		if file.FileInfo().IsDir() {
			try(os.MkdirAll(destPath, 0755))
		} else {
			try(os.MkdirAll(filepath.Dir(destPath), 0755))
			srcFile := try1(file.Open())
			defer srcFile.Close()
			destFile := try1(os.Create(destPath))
			defer destFile.Close()
			try1(io.Copy(destFile, srcFile))
		}
	}
	return nil
}
