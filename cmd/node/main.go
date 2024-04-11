package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jcbhmr/go-node/internal/nodejsinstall"
	"github.com/jcbhmr/go-node/internal/osutil"
)

const v = "20.12.1"

var cacheDir string

func init() {
	log.SetFlags(0)
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		userCacheDir = os.TempDir()
	}
	cacheDir = filepath.Join(userCacheDir, "go-node")
}

func ptr[T any](v T) *T {
	return &v
}

func main() {
	root, bin, err := nodejsinstall.Install(ptr(v), ptr(filepath.Join(cacheDir, "node", v)))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Installed Node.js %s to %s", v, root)

	path := filepath.Join(bin, "node")
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	path, err = filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}
	err = osutil.Execve(path, os.Args, os.Environ())
	if err != nil {
		log.Fatalln(err)
	}
}
