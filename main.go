package main

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// cleanResult holds task result information.
type cleanResult struct {
	vsPath      string
	backupPath  string
	backupCount int
	err         error
}

// versionInfo holds version information.
type versionInfo struct {
	fileName string
	version  string
}

// cleanVSSetupDir run clean function to move old version dirs in vsSetupPath to
// a backup path.
func cleanVSSetupDir(vsSetupPath string, showActionOnly bool) *cleanResult {
	result := &cleanResult{
		vsPath:      "",
		backupPath:  "",
		backupCount: 0,
		err:         nil,
	}

	// make sure the given path is ok.
	if result.vsPath, result.err = filepath.Abs(vsSetupPath); result.err != nil {
		return result
	}

	// read file list in the give path. the result is sorted by file name, asc.
	var list []fs.DirEntry
	if list, result.err = os.ReadDir(result.vsPath); result.err != nil {
		return result
	}

	// create backup dir just inside give visual studio setup path.
	if result.backupPath, result.err = createBackupDir(vsSetupPath, showActionOnly); result.err != nil {
		return result
	}

	backupOldVersionDirs(result, &list, showActionOnly)
	return result
}

// backupOldVersionDirs keeps the newest version dir for each component, move elder dirs to backup path.
func backupOldVersionDirs(result *cleanResult, list *[]fs.DirEntry, showActionOnly bool) {
	dirs := make(map[string]versionInfo)

	for _, fi := range *list {
		if fi.IsDir() { // we only process dirs.
			key, ver := getDirInfo(fi.Name())

			if ver == "" {
				continue // skip dirs without version.
			}

			// different version of same component dirs have same key extracted by getDirInfo().
			if info, ok := dirs[key]; ok {
				compareResult := compareVersion(info.version, ver)

				if compareResult < 0 {
					// exist version is older than current one, backup exist one and save current one.
					backupDir(result, ver, info.fileName, showActionOnly)
				} else if compareResult > 0 {
					// exist version is newer than current one, backup current one and no need to change the map.
					backupDir(result, info.version, fi.Name(), showActionOnly)
					continue
				} else { // compareResult == 0. it should never run to here, just in case for its happening.
					result.err = errors.New("IMPOSSIBLE: key[" + key + "], old version=" + info.version + ", new version=" + ver)
					break
				}

				result.backupCount++
			}

			// replace exist one with newer version info or just add it when it isn't in the map.
			dirs[key] = versionInfo{
				fileName: fi.Name(),
				version:  ver,
			}
		}
	}
}

// backupDir moves given dir specified by fileName to backup path.
func backupDir(result *cleanResult, version string, fileName string, showActionOnly bool) {
	if showActionOnly {
		println("New version [", version, "] is found, [", fileName, "] will be moved.")
	} else {
		oldPath := path.Join(result.vsPath, fileName)
		newPath := path.Join(result.backupPath, fileName)

		println("Moving ...... [", fileName, "]")
		result.err = os.Rename(oldPath, newPath)
	}
}

// compareVersion compares two version.
// return 1 for version1 is newer, -1 for version2 is newer, 0 for equal or something wrong.
// because given value should never be equal, so return 0 means error.
func compareVersion(version1 string, version2 string) int {
	ss1 := strings.Split(version1, ".")
	ss2 := strings.Split(version2, ".")

	count := len(ss1)
	temp := len(ss2)
	if temp < count {
		count = temp
	}

	for i := 0; i < count; i++ {
		v1, err1 := strconv.Atoi(ss1[i])
		v2, err2 := strconv.Atoi(ss2[i])

		if err1 != nil || err2 != nil {
			return 0
		}

		if v1 > v2 {
			return 1
		} else if v1 < v2 {
			return -1
		}
	}

	return 0
}

// createBackupDir creates specified root dir when showActionOnly is false.
// It returns the backup dir name.
func createBackupDir(root string, showActionOnly bool) (string, error) {
	dirName := path.Join(root, time.Now().Format("_2006-01-02_15-04-05"))

	if showActionOnly {
		return dirName, nil
	}

	err := os.Mkdir(dirName, os.ModeDir)
	return dirName, err
}

// getDirInfo returns dir name without version segment and
// the version string extracted from version segment.
func getDirInfo(dirName string) (string, string) {
	// The componet dir name is like:
	// 1) Microsoft.Net.Core.SDK.MSBuildExtensions,version=16.11.31603.221
	// 	  Microsoft.Net.Core.SDK.MSBuildExtensions,version=16.11.31701.289
	// 2) AndroidNDK_R16B,version=16.0,chip=x64
	//    AndroidNDK_R16B,version=16.0,chip=x86
	// The version number was preceded by 'version='.
	// A few directory names do not have a version number, such as `certificates`.
	// Return values for group 1 are:
	//    "Microsoft.Net.Core.SDK.MSBuildExtensions,", "16.11.31603.221"
	//    "Microsoft.Net.Core.SDK.MSBuildExtensions,", "16.11.31701.289"
	// Return values for group 2 are:
	//	  "AndroidNDK_R16B,chip=x64,", "16.0"
	//	  "AndroidNDK_R16B,chip=x86,", "16.0"
	ss := strings.Split(dirName, ",")
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

	return name, ver
}

func main() {
	setupPath := flag.String("path", ".", "Visualt Studio's setup path")
	showActionOnly := flag.Bool("showonly", false, "Show process only, no real action")
	flag.Parse()

	result := cleanVSSetupDir(*setupPath, *showActionOnly)

	if result.err != nil {
		println("-----------------------------------")
		println("Clean Visual Studio setup dir fail:\n", result.err.Error())
		println("-----------------------------------")
	} else {
		println("---------------------------------------")
		println("Clean Visual Studio setup dir finished.")
		println("Visual Studio setup path is [", result.vsPath, "].")
		println("Old version dirs will be moved to [", result.backupPath, "].")

		if result.backupCount == 0 {
			println("Old version dir is not found. Nothing to do.")
		} else if *showActionOnly {
			println(result.backupCount, "old version dirs will be moved.")
		} else {
			println(result.backupCount, "old version dirs have been moved.")
		}

		println("---------------------------------------")
	}
}
