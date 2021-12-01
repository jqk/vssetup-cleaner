package main

import (
	"errors"
	"flag"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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
	lastFileName := ""
	lastName := ""
	lastVer := ""

	for _, fi := range *list {
		if fi.IsDir() { // we only process dirs.
			fileName := fi.Name()
			name, ver := getDirInfo(fileName)

			if ver != "" { // skip dirs without version.
				if lastName == name { // deferent version of same component dirs have same name extracted by getDirInfo().
					// because the file name list is sorted, elder version is preceed the new one.
					// we should never go into the other two branches.
					if lastVer < ver {
						if showActionOnly {
							println("New version [", ver, "] is found, [", lastFileName, "] will be moved.")
						} else {
							oldPath := path.Join(result.vsPath, lastFileName)
							newPath := path.Join(result.backupPath, lastFileName)

							println("Moving ...... [", lastFileName, "]")
							if result.err = os.Rename(oldPath, newPath); result.err != nil {
								break
							}
						}

						result.backupCount++
					} else if lastVer == ver {
						result.err = errors.New("IMPOSSIBLE: " + name + " has duplicated version " + ver + "")
						break
					} else {
						result.err = errors.New("IMPOSSIBLE: " + name + " version '" + ver + "' is after " + lastName)
						break
					}
				}
			}

			lastFileName = fileName
			lastName = name
			lastVer = ver
		}
	}
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

// getDirInfo returns dir name without version and
// the version info according to given dirName.
func getDirInfo(dirName string) (string, string) {
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
