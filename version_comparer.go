package main

import (
	"errors"
	"os"

	"github.com/jqk/futool4go/common"
)

/*
FindOldPackagesByDirVersion find old packages in given dir list.

FindOldPackagesByDirVersion 在给定的目录列表中查找过期的包。
*/
func FindOldPackagesByDirVersion(vsSetupPath string) ([]*PackageInfo, error) {
	// 读取 vsSetupPath 下的目录列表。
	packageDirs, err := os.ReadDir(vsSetupPath)
	if err != nil {
		return nil, err
	}

	// VS 安装目录下一般有 3000 多个包。每个小版本更新其中的数十个。给数组预先分配一些空间。
	oldPackages := make([]*PackageInfo, 0, 100)      // 过期包集合。
	currentPackages := make(map[string]*PackageInfo) // 当前包集合。以包名为 key，存储包信息。

	for _, dir := range packageDirs {
		if dir.IsDir() { // we only process dirs.
			nextPack := getPackageInfo(dir.Name())
			if nextPack == nil {
				continue
			}

			if currentPack, ok := currentPackages[nextPack.Name]; !ok {
				// 在当前包集合内，没有找到刚刚从目录包转换得到的包信息。需将其添加到当前包集合。
				currentPackages[nextPack.Name] = nextPack
			} else {
				// 当前包集合内已存在刚刚从目录包转换得到的同名包信息，则需要比较版本。
				compareResult := common.CompareVersions(currentPack.Version, nextPack.Version)

				if compareResult < 0 {
					// current version is older than the next one.
					// move current package info to obsolete package list.
					// and save next package info to current package list.
					oldPackages = append(oldPackages, currentPack)
					currentPackages[nextPack.Name] = nextPack
				} else if compareResult > 0 {
					// current version is newer than the next one.
					// it should not happened becase the packageDirs is sorted by os.ReadDir().
					// but, in case it is happening, keep the code below. no need to change currentPackages.
					oldPackages = append(oldPackages, nextPack)
				} else {
					// the 2 versions are identical. it should never run to here.
					return nil, errors.New("Duplicated packages: " + nextPack.Name + ", old version=" +
						currentPack.Version + ", new version=" + nextPack.Version)
				}
			}
		}
	}

	return oldPackages, nil
}
