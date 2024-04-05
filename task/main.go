package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

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

func downloadFile(srcUrl string, destPath string) (rerr error) {
	defer catch(&rerr, nil)
	try(os.MkdirAll(path.Dir(destPath), 0755))
	file := try1(os.Create(destPath))
	defer file.Close()
	res := try1(http.Get(srcUrl))
	defer res.Body.Close()
	try1(io.Copy(file, res.Body))
	return nil
}

func fetchNodejsDist() (rerr error) {
	defer catch(&rerr, nil)
	const nodejsVersion = "20.12.1"
	archiveUrls := map[string]string{
		"darwin-arm64":  fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-arm64.tar.gz", nodejsVersion),
		"darwin-x64":    fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-x64.tar.gz", nodejsVersion),
		"linux-arm64":   fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-arm64.tar.gz", nodejsVersion),
		"linux-armv7l":  fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-armv7l.tar.gz", nodejsVersion),
		"linux-ppc64le": fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-ppc64le.tar.gz", nodejsVersion),
		"linux-s390x":   fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-s390x.tar.gz", nodejsVersion),
		"linux-x64":     fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-x64.tar.gz", nodejsVersion),
		"windows-x64":   fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x64.zip", nodejsVersion),
		"windows-x86":   fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x86.zip", nodejsVersion),
	}
	for goOsAndGoArch, archiveUrl := range archiveUrls {
		archivePath := strings.Replace(archiveUrl, "https://nodejs.org/dist/", "internal/nodejs_dist/", 1)
		try(downloadFile(archiveUrl, archivePath))
		
	}
	return nil
}

func cleanNodejsDist() (rerr error) {
	defer catch(&rerr, nil)
	try(os.RemoveAll("internal/nodejs_dist"))
	return nil
}

func main() {
	switch os.Args[1] {
	case "fetch-nodejs-dist":
		try(fetchNodejsDist())
	case "clean-nodejs-dist":
		try(cleanNodejsDist())
	}
}
