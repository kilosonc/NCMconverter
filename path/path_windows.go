package path

import (
	"path/filepath"
)

func Clean(path string) string {
	return filepath.Clean(path)
}

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

func Base(p string) string {
	return filepath.Base(p)
}

func Ext(p string) string {
	return filepath.Ext(p)
}

func Dir(p string) string {
	return filepath.Dir(p)
}
