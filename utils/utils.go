package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"github.com/creamdog/gonfig"
)

func ProcessError(e error) {
    if e != nil {
        panic(e)
    }
}

func GetRootPath() string {
    var (
    	_, b, _, _ = runtime.Caller(0)
    	basepath = filepath.Dir(b)
    )
    return basepath + "/.."
}

func GetConfig() gonfig.Gonfig {
	path := GetRootPath()
	f, err := os.Open(path + "/config.json")
	ProcessError(err)
	defer f.Close()
	config, err := gonfig.FromJson(f)
	ProcessError(err)
	return config
 }

func GetIntConfigVar(config gonfig.Gonfig, path string) int {
	prefix, err := config.GetInt(path, 0)
	ProcessError(err)
	return prefix
}

func GetStringConfigVar(config gonfig.Gonfig, path string) string {
	prefix, err := config.GetString(path, "")
	ProcessError(err)
	return prefix
}