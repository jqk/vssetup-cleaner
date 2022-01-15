package main

import (
	"flag"
)

func main() {
	const no_package_list_file = ""

	setupPath := flag.String("vs", ".", "Visual Studio's setup path")
	showActionOnly := flag.Bool("show", true, "Show action only")
	packageListFile := flag.String("pkg", no_package_list_file, "package list file")
	flag.Parse()

	var result *CleanResult

	if *packageListFile == no_package_list_file {
		result = CleanVSSetupDir(*setupPath, *showActionOnly)
	} else {
		result = CleanVSSetupDirWithCheckFile(*setupPath, *showActionOnly, *packageListFile)
	}

	if result.err != nil {
		println("-----------------------------------")
		println("Clean Visual Studio setup dir fail:\n", result.err.Error())
		println("-----------------------------------")
	} else {
		println("---------------------------------------")
		println("Clean Visual Studio setup dir finished.")
		println("Visual Studio setup path is [", result.vsPath, "].")
		println("Old version packages will be moved to [", result.backupPath, "].")

		if result.backupCount == 0 {
			println("Old version package is not found. Nothing to do.")
		} else if *showActionOnly {
			println(result.backupCount, "old version packages will be moved.")
		} else {
			println(result.backupCount, "old version packages have been moved.")
		}

		println("---------------------------------------")
	}
}
