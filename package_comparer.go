package main

import (
	"bufio"
	"os"
	"strings"
)

/*
FindOldPackagesByFileList find old packages according to package list file.

FindOldPackagesByFileList 以目录列表中给出的包为基础，查找指定目录下过期的包。
*/
func FindOldPackagesByFileList(vsSetupPath string, listFilename string) ([]*PackageInfo, error) {
	// 读取 vsSetupPath 下的目录列表。
	packageDirs, err := os.ReadDir(vsSetupPath)
	if err != nil {
		return nil, err
	}

	packages, err := readPackageListFile(listFilename)
	if err != nil {
		return nil, err
	}

	// VS 安装目录下一般有 3000 多个包。每个小版本更新其中的数十个。给数组预先分配一些空间。
	result := make([]*PackageInfo, 0, 100)

	for _, fn := range packageDirs {
		if fn.IsDir() { // we only process dirs.
			dirName := fn.Name()

			// ok == false means the dir in setup path is not in the package list.
			// it should be moved out in most cases.
			if _, ok := packages[dirName]; !ok {
				pkg := getPackageInfo(dirName)
				if pkg == nil {
					continue
				}

				// 目录具有版本信息，但又不在包列表中。这个目录应该被移除。
				result = append(result, pkg)
			}
		}
	}

	return result, nil
}

// readPackageListFile reads package list file create by vs_layout.exe.
func readPackageListFile(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// by visual studio version 16.8 (2019), the package list file has 3633 lines.
	// give it a little bit more space.
	list := make(map[string]bool, 4000)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()       // read a text line.
		path := getPackagePath(line) // analyse it, get package name.

		if path != notPackage {
			list[path] = true // save package name to the list.
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

const packageHeader = "Verified existing package '"
const packageHeaderLenght = len(packageHeader)
const notPackage = ""

// getPackagePath get package name as the result, otherwise returns empty string.
func getPackagePath(line string) string {
	// each line is like below:
	// Verified existing package 'Unity3d.x86,version=3.1,chip=x86'
	// the package name is the dir name, quoted with "'".
	if strings.Index(line, packageHeader) == 0 {
		end := strings.LastIndex(line, "'")

		if end > packageHeaderLenght {
			s := line[packageHeaderLenght:end]
			return s
		}
	}

	return notPackage
}
