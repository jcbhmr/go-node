package nodejsinstall

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/codeclysm/extract"
)

var archiveBytes []byte

func getArchive(version string) (io.Reader, error) {
	if archiveBytes == nil {
		var archiveUrl string
		switch runtime.GOOS + "/" + runtime.GOARCH {
		case "darwin/arm64":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-arm64.tar.gz", version)
		case "darwin/amd64":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-darwin-x64.tar.gz", version)
		case "linux/arm64":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-arm64.tar.gz", version)
		case "linux/armv7l":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-armv7l.tar.gz", version)
		case "linux/ppc64le":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-ppc64le.tar.gz", version)
		case "linux/s390x":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-s390x.tar.gz", version)
		case "linux/amd64":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-linux-x64.tar.gz", version)
		case "windows/amd64":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x64.zip", version)
		case "windows/386":
			archiveUrl = fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%[1]s-win-x86.zip", version)
		default:
			return nil, fmt.Errorf("unsupported platform: %s-%s", runtime.GOOS, runtime.GOARCH)
		}
		res, err := http.Get(archiveUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", archiveUrl, err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s not OK: %s", res.Request.URL, res.Status)
		}
		if strings.Contains(res.Header.Get("Content-Type"), "text/html") {
			panic(fmt.Errorf("%s returned %s html", archiveUrl, res.Status))
		}
		archiveBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body of %s: %w", res.Request.URL, err)
		}
	}
	return bytes.NewReader(archiveBytes), nil
}

func getLatestNodejsVersion() (string, error) {
	const url = "https://nodejs.org/dist/index.json"
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s not OK: %s", res.Request.URL, res.Status)
	}
	if res.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("%s not application/json: %s", res.Request.URL, res.Header.Get("Content-Type"))
	}
	var versions []struct{ Version string }
	if err := json.NewDecoder(res.Body).Decode(&versions); err != nil {
		return "", fmt.Errorf("failed to decode %s: %w", res.Request.URL, err)
	}
	return versions[0].Version, nil
}

func Install(vOpt *string, nodejsInstallOpt *string) (string, string, error) {
	var v string
	if vOpt == nil {
		latest, err := getLatestNodejsVersion()
		if err != nil {
			return "", "", err
		}
		v = latest
	} else {
		v = *vOpt
	}
	var nodejsInstall string
	if nodejsInstallOpt == nil {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			userCacheDir = os.TempDir()
		}
		cacheDir := filepath.Join(userCacheDir, "nodejsinstall", v)
		nodejsInstall = cacheDir
	} else {
		nodejsInstall = *nodejsInstallOpt
	}
	archiveReader, err := getArchive(v)
	if err != nil {
		return "", "", err
	}
	err = os.MkdirAll(nodejsInstall, 0755)
	if err != nil {
		return "", "", err
	}
	err = extract.Archive(context.TODO(), archiveReader, nodejsInstall, func(s string) string {
		parts := filepath.SplitList(s)
		parts = parts[1:]
		return filepath.Join(parts...)
	})
	if err != nil {
		_ = os.RemoveAll(nodejsInstall)
		return "", "", err
	}
	var bin string
	if runtime.GOOS == "windows" {
		bin = nodejsInstall
	} else {
		bin = path.Join(nodejsInstall, "bin")
	}
	return nodejsInstall, bin, nil
}
