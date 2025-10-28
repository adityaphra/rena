package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

type file struct {
	baseDir string
	name    string
	ext     string
}

func (f *file) getFullPath() string {
	return path.Join(f.baseDir, f.name+f.ext)
}

func (f *file) getFullName() string {
	return f.name + f.ext
}

func (f *file) setFullName(value string) {
	f.ext = path.Ext(value)
	f.name = strings.TrimSuffix(value, f.ext)
}

func (f *file) getExt() string {
	return strings.TrimPrefix(f.ext, ".")
}

func CreateFile(filePath string) file {
	base, fileName := path.Split(filePath)
	ext := path.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	return file{baseDir: base, name: name, ext: ext}
}

func RenameFile(old file, new file) error {
	newBaseDir := new.baseDir
	if old.baseDir != new.baseDir {
		info, err := os.Stat(newBaseDir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(newBaseDir, 0755); err != nil {
				return fmt.Errorf("Can't create directory: %v", newBaseDir)
			}
		} else if !info.IsDir() {
			return fmt.Errorf("%v is not a directory", newBaseDir)
		}
	}

	oldPath := expandPath(old.getFullPath())
	newPath := expandPath(new.getFullPath())
	err := os.Rename(oldPath, newPath)
	if err == nil {
		return nil
	}

	// fallback to copy then remove
	in, err := os.Open(oldPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	if err = out.Sync(); err != nil {
		return err
	}

	return os.Remove(oldPath)
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			return path.Join(homeDir, p[2:])
		}
	}
	return p
}

// references:
// - https://stackoverflow.com/questions/1976007/what-characters-are-forbidden-in-windows-and-linux-directory-names
func IsSafeName(f file) bool {
	fullName := f.getFullName()

	if runtime.GOOS == "windows" {
		forbiddenWords := []string{
			"CON", "PRN", "AUX", "NUL",
			"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
			"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
		}
		for _, word := range forbiddenWords {
			if strings.EqualFold(word, f.name) {
				return false
			}
		}

		forbiddenSymbols := []string{
			"<", ">", ":", "\"", "/",
			"\\", "|", "?", "*",
		}
		for _, symbol := range forbiddenSymbols {
			if strings.Contains(fullName, symbol) {
				return false
			}
			if symbol != "/" && symbol != "\\" {
				if strings.Contains(f.baseDir, symbol) {
					return false
				}
			}
		}
	}

	if runtime.GOOS == "linux" {
		if strings.Contains(fullName, "/") || fullName == "." || fullName == ".." {
			return false
		}
	}

	if runtime.GOOS == "darwin" {
		if strings.Contains(fullName, "/") || strings.Contains(fullName, ":") {
			return false
		}
	}

	return true
}
