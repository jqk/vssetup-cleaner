package main

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
)

// versionInfo holds version information.
type versionInfo struct {
	fileName string
	version  string
}

const version_is_newer = 1
const version_is_older = -1
const version_is_error = 0

// CleanVSSetupDir run clean function to move old version dirs in vsSetupPath to
// a backup path.
func CleanVSSetupDir(vsSetupPath string, showActionOnly bool) *CleanResult {
	result, list := PrepareEnvironment(vsSetupPath, showActionOnly)
	if result.err != nil {
		return result
	}

	backupOldVersionDirs(result, list, showActionOnly)
	return result
}

// backupOldVersionDirs keeps the newest version dir for each component, move elder dirs to backup path.
func backupOldVersionDirs(result *CleanResult, list *[]fs.DirEntry, showActionOnly bool) {
	dirs := make(map[string]versionInfo)

	for _, fi := range *list {
		if fi.IsDir() {
			// we only process dirs. try to get the version of current dir.
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
				} else {
					// compareResult == 0. it should never run to here, just in case for its happening.
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
func backupDir(result *CleanResult, version string, fileName string, showActionOnly bool) {
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
	// versions are 3 or 4 numbers separated by dot.
	ss1 := strings.Split(version1, ".")
	ss2 := strings.Split(version2, ".")

	len1 := len(ss1)
	len2 := len(ss2)

	result := version_is_error // default return value.
	count := len1              // count of numbers should compare between two versions.

	if len1 > len2 {
		// version1 has more numbers than version2. it may contains version2.
		// version1 is newer if it contains version2.
		result = version_is_newer

		// only check for same length.
		count = len2
	} else if len1 < len2 {
		// version2 is newer if it contains version1.
		result = version_is_older
	}

	for i := 0; i < count; i++ {
		// we must convert string to int because '16' is less than '2',
		// but 16 is greater than 2.
		v1, err1 := strconv.Atoi(ss1[i])
		v2, err2 := strconv.Atoi(ss2[i])

		if err1 != nil || err2 != nil {
			return version_is_error
		}

		if v1 > v2 {
			return version_is_newer
		} else if v1 < v2 {
			return version_is_older
		}
	}

	return result
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
