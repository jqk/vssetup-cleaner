package main

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jqk/futool4go/fileutils"
	"github.com/jqk/futool4go/timeutils"
)

/*
Clean cleans up obsolete packages in Visual Studio Setup directory.

Parameters:
  - setupPath: The path of the Visual Studio Setup.
  - listFilename: The path of the package list file.
  - handler: The callback function for handling expired packages. Cannot be nil.
  - showOnly: Whether show obsolete package or do the real clean action.
  - needDirStat: Whether to statistic directories.

Returns:
  - Information about the cleaning process.
  - Error. Returns nil if the handler returns SkipAll or SkipDir.

Clean 用于清理 Visual Studio Setup 目录中的过期包。

参数:
  - setupPath: Visual Studio Setup 的路径。
  - listFilename: 包列表文件的路径。
  - handler: 处理过期包的回调函数。不能为 nil。
  - showOnly: 是否仅显示处理过程。
  - needDirStat: 是否需要统计目录数量、文件大小。

返回:
  - 清理过程的信息。
  - 错误。如果 handler 返回 SkipAll 或 SkipDir，则返回 nil。
*/
func Clean(
	setupPath string,
	listFilename string,
	handler ObsoletePackageHandler,
	showOnly bool,
	needDirStat bool) (*CleanInfo, error) {

	useListFile := strings.TrimSpace(listFilename) != ""
	if err := varifyParameters(setupPath, listFilename, handler, useListFile); err != nil {
		return nil, err
	}

	info := &CleanInfo{}
	// 将可能的相对路径转换为绝对路径，并创建备份目录名称。
	info.SetupPath, _ = filepath.Abs(setupPath)
	info.BackupPath, _ = filepath.Abs(path.Join(setupPath, time.Now().Format("_2006-01-02_15-04-05")))

	var errOuter error
	sw := timeutils.Stopwatch{}

	info.ElapsedTime, errOuter = sw.Elapsing(func() error {
		var err error

		// 查找旧包。
		if useListFile {
			info.Packages, err = FindOldPackagesByFileList(info.SetupPath, listFilename)
		} else {
			info.Packages, err = FindOldPackagesByDirVersion(info.SetupPath)
		}
		if err != nil {
			return err
		}

		if !showOnly && len(info.Packages) > 0 {
			// 只有必要时才创建备份目录，确保其可用。
			if err := os.Mkdir(info.BackupPath, os.ModeDir); err != nil {
				return err
			}
		}

		// 逐一处理旧包。
		for _, pkg := range info.Packages {
			if needDirStat {
				pkgPath := path.Join(info.SetupPath, pkg.Dir)
				if pkg.DirCount, pkg.FileCount, pkg.Size, err = fileutils.GetDirStatistics(pkgPath); err != nil {
					return err
				}
			}
			// 调用回调函数处理旧包。返回任何错误均结束处理。
			if err = handler(pkg); err != nil {
				return err
			}
			if !showOnly {
				if err = backupObsoletePackage(info, pkg); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if errOuter == filepath.SkipAll || errOuter == filepath.SkipDir {
		// SkipAll 或 SkipDir 被认为是正常的信号。
		errOuter = nil
	}
	return info, errOuter
}

func backupObsoletePackage(info *CleanInfo, oldPackage *PackageInfo) error {
	oldPath := path.Join(info.SetupPath, oldPackage.Dir)
	newPath := path.Join(info.BackupPath, oldPackage.Dir)

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

func varifyParameters(setupPath string, listFilename string, handler ObsoletePackageHandler, useListFile bool) error {
	if handler == nil {
		return errors.New("handler is nil")
	}

	// 判断 setupPath 是否有效。
	if exists, isDir, err := fileutils.FileExists(setupPath); err != nil {
		return err
	} else if !exists || !isDir {
		return errors.New("invalid Visual Studio setup path: " + setupPath)
	}

	if useListFile {
		// 判断 listFilename 是否有效。
		if exists, isDir, err := fileutils.FileExists(listFilename); err != nil {
			return err
		} else if !exists || isDir {
			return errors.New("invalid Visual Studio package list file: " + listFilename)
		}
	}

	return nil
}
