package config

import (
	"path"
	"runtime"
)

var RootDir string

func InitConfig() {
	_, filename, _, _ := runtime.Caller(0)
	RootDir = path.Dir(path.Dir(filename))
}
