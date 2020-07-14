package app

import (
	"github.com/kardianos/osext"
	"path/filepath"
)

var BinPath string
var ResourcesPath string
var BuiltinCloudSQLProxy string

func init() {
	binPath, err := osext.Executable()
	if err != nil {
		panic(err)
	}

	BinPath = binPath
	ResourcesPath = filepath.Join(filepath.Dir(filepath.Dir(BinPath)), "Resources")
	BuiltinCloudSQLProxy = filepath.Join(ResourcesPath, "bin", "cloud_sql_proxy")
}
