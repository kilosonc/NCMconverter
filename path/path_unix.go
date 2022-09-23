package path

import (
	"path"
)

func Clean(p string) string {
	return path.Clean(p)
}

func Join(elem ...string) string {
	return path.Join(elem...)
}

func Base(p string) string {
	return path.Base(p)
}

func Ext(p string) string {
	return path.Ext(p)
}

func Dir(p string) string {
	return path.Dir(p)
}
