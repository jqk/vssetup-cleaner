package main

import (
	"strings"
	"time"
)

/*
PackageInfo describes information of a package.

PackageInfo 是每一个包的信息。
*/
type PackageInfo struct {
	Dir       string // the directory of the package, not full path.
	Name      string // the name of the package.
	Version   string // the version of the package.
	DirCount  int    // the number of directories.
	FileCount int    // the number of files.
	Size      int64  // the size of the package.
}

/*
ObsoletePackageHandler is a callback function that handles obsolete packages.

Parameters:
  - packageInfo: the information of an obsolete package, not nil.

Returns:
  - any error will terminate the clean process.
    SkipAll or SkipDir is considered a valid signal to terminate the process.

ObsoletePackageHandler 是处理过期包的回调函数。

参数:
  - packageInfo: 一个过期包的信息，不为 nil。

返回:
  - 任何错误都会终止清理过程。SkipAll 和 SkipDir 被视为有效的信号来终止过程。
*/
type ObsoletePackageHandler func(packageInfo *PackageInfo) error

/*
CleanInfo describes information of a clean process.

CleanInfo 说明清理过程的信息。
*/
type CleanInfo struct {
	SetupPath   string         // the path of Visual Studio Setup.
	BackupPath  string         // the backup path of obsolete packages.
	Packages    []*PackageInfo // the list of obsolete packages.
	ElapsedTime time.Duration  // the elapsed time of the clean process.
}

/*
getPackageInfo returns package name and version according to given dir name.

See test for more detail.
*/
func getPackageInfo(dir string) *PackageInfo {
	// 可以通过正则 ^(.+),version=((\d+\.?)+)(,chip=(\w+))?(,language=(\w+\-?\w+))?
	// 解析出所需字段，但相对麻烦一些，所以这里采用直接解析的方式。
	// 包的目录名是这样的：
	// Microsoft.VisualStudio.Web.Mvc4,version=16.11.115.10959
	// Win10SDK_10.0.16299,version=10.0.16299.5
	// Microsoft.VisualStudio.Vsto.Runtime.Resources,version=16.0.28315.86,chip=x64,language=zh-CN
	// Microsoft.VisualStudio.Vsto.Runtime,version=16.0.28315.86,chip=x64
	// Microsoft.Build,version=16.11.2.2150704,chip=neutral,language=neutral
	// Microsoft.DiagnosticsHub.DatabaseTool.Resources,version=16.11.31901.3,language=zh-CN

	ss := strings.Split(dir, ",")
	ver := ""
	name := ""

	for _, s := range ss {
		if strings.Contains(s, "version=") {
			vv := strings.Split(s, "=")

			if len(vv) >= 2 {
				ver = vv[1]
			}
		} else {
			name += s + ","
		}
	}

	if ver == "" {
		// skip dirs without version, such as dir 'certificates'.
		return nil
	}

	return &PackageInfo{
		Dir:     dir,
		Name:    name,
		Version: ver,
	}
}
