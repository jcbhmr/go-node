package osutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Touch(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func WaitUntilRemoved(path string) error {
	if exists, err := Exists(path); err != nil {
		return err
	} else if exists {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return err
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher closed")
			}
			if exists, _ := Exists(path); !exists {
				return nil
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher closed")
			}
			return err
		}
	}
}

func WaitUntilCreated(path string) error {
	if exists, err := Exists(path); err != nil {
		return err
	} else if !exists {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	if err := watcher.Add(filepath.Dir(path)); err != nil {
		return err
	}

	for {
		select {
		case _, ok := <-watcher.Events:
			if !ok {
				return errors.New("watcher closed")
			}
			if exists, _ := Exists(path); exists {
				return nil
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("watcher closed")
			}
			return err
		}
	}
}

func Exists(path string) (bool, error) {
	if _, err := os.Lstat(path); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func ReplaceExecutableWithSymlink(dest string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	err = os.Rename(exe, exe + ".old")
	if err != nil {
		return err
	}
	err = os.Symlink(dest, exe)
	if err != nil {
		return err
	}
	if err := os.Remove(exe + ".old"); err != nil {
		_ = os.Rename(exe + ".old", filepath.Join(os.TempDir(), filepath.Base(exe) + ".old"))
	}
	return nil
}
