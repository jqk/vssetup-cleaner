package main

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"
)

// CleanResult holds task result information.
type CleanResult struct {
	vsPath      string
	backupPath  string
	backupCount int
	err         error
}

// CreateBackupDir creates specified root dir when showActionOnly is false.
// It returns the backup dir name.
func CreateBackupDir(root string, showActionOnly bool) (dirName string, err error) {
	dirName = path.Join(root, time.Now().Format("_2006-01-02_15-04-05"))
	if dirName, err = filepath.Abs(dirName); err != nil {
		return dirName, err
	} else if showActionOnly {
		return dirName, nil
	}

	err = os.Mkdir(dirName, os.ModeDir)
	return dirName, err
}

// PrepareEnvironment prepare working environment, load vs dir list.
func PrepareEnvironment(vsSetupPath string, showActionOnly bool) (*CleanResult, *[]fs.DirEntry) {
	result := &CleanResult{
		vsPath:      "",
		backupPath:  "",
		backupCount: 0,
		err:         nil,
	}

	// make sure the given path is ok.
	if result.vsPath, result.err = filepath.Abs(vsSetupPath); result.err != nil {
		return result, nil
	}

	// read file list in the give path. the result is sorted by file name, asc.
	var list []fs.DirEntry
	if list, result.err = os.ReadDir(result.vsPath); result.err != nil {
		return result, nil
	}

	// create backup dir just inside given visual studio setup path.
	if result.backupPath, result.err = CreateBackupDir(vsSetupPath, showActionOnly); result.err != nil {
		return result, nil
	}

	return result, &list
}
