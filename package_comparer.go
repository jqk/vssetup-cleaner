package main

import (
	"bufio"
	"os"
	"path"
	"strings"
)

// packageList is used as a list.
type packageList map[string]bool

// not_a_package is the symbol that it is not a package name.
const not_a_package = ""

// CleanVSSetupDirWithCheckFile run clean function to move old version dirs in vsSetupPath to
// a backup path according to vs_installer checking result.
func CleanVSSetupDirWithCheckFile(vsSetupPath string, showActionOnly bool, packageListFile string) *CleanResult {
	result, list := PrepareEnvironment(vsSetupPath, showActionOnly)
	if result.err != nil {
		return result
	}

	// read package list file.
	var packages packageList
	if packages, result.err = readPackageListFile(packageListFile); result.err != nil {
		return result
	}

	for _, fn := range *list {
		if fn.IsDir() { // we only process dirs.
			dirName := fn.Name()

			// ok == false means the dir in setup path is not in the package list.
			// it should be moved out in most cases.
			if _, ok := packages[dirName]; !ok {
				// 'Archive' is the install history. It is removable, but please do it manually.
				// I think we should leave 'certificates' there although I haven't tried it yet.
				// Dir starts with '_' is the backup directory. Remove it manually.
				if dirName == "Archive" || dirName == "certificates" || strings.Index(dirName, "_") == 0 {
					continue
				}

				// dirs that are not in the package list will be moved out.
				if backupAbandonedPackage(result, dirName, showActionOnly); result.err != nil {
					break
				}

				result.backupCount++
			}
		}
	}

	return result
}

// backupAbandonedPackage moves abandoned packages to backupDir.
func backupAbandonedPackage(result *CleanResult, packageName string, showActionOnly bool) {
	if showActionOnly {
		println("Abandoned package [", packageName, "] will be moved.")
	} else {
		oldPath := path.Join(result.vsPath, packageName)
		newPath := path.Join(result.backupPath, packageName)

		println("Moving package ...... [", packageName, "]")
		result.err = os.Rename(oldPath, newPath)
	}
}

// readPackageListFile reads package list file create by vs_layout.exe.
func readPackageListFile(filename string) (packageList, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// by visual studio version 16.8 (2019), the package list file has 3633 lines.
	// give it a little bit more space.
	list := make(packageList, 4000)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()       // read a text line.
		path := getPackagePath(line) // analyse it, get package name.

		if path != not_a_package {
			list[path] = true // save package name to the list.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// getPackagePath get package name as the result, otherwise returns empty string.
func getPackagePath(line string) string {
	if strings.Index(line, "Verified existing package '") == 0 {
		// each line is like below:
		// Verified existing package 'Unity3d.x86,version=3.1,chip=x86'
		// the package name is the dir name, quoted with "'".
		start := 27 // the length of "Verified existing package '". calculated manually.
		end := strings.LastIndex(line, "'")

		if end > start {
			s := line[start:end]
			return s
		}
	}

	return not_a_package
}
