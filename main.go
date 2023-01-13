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

	println("Visual Studio setup file cleaner 1.0.0\n")
	println("  Arguments from command line or default value:")
	println("    pkg  = ", *packageListFile)
	println("    show = ", *showActionOnly)
	println("    vs   = ", *setupPath)
	println("\n---------------------------------------")

	var result *CleanResult

	if *packageListFile == no_package_list_file {
		result = CleanVSSetupDir(*setupPath, *showActionOnly)
	} else {
		result = CleanVSSetupDirWithCheckFile(*setupPath, *showActionOnly, *packageListFile)
	}

	if result.err != nil {
		println("---------------------------------------")
		println("Clean Visual Studio setup dir fail:\n", result.err.Error())
		println("---------------------------------------")
	} else {
		if result.backupCount > 0 {
			println("\n---------------------------------------\n")
		}
		println("Clean Visual Studio setup dir finished.")
		println("Visual Studio setup path is [", result.vsPath, "].")
		println("Old version packages will be moved to [", result.backupPath, "].")

		if result.backupCount == 0 {
			println("Old version package is not found. Nothing to do.\n")
		} else {
			packageWord := "packages"
			if result.backupCount == 1 {
				packageWord = "package"
			}

			if *showActionOnly {
				println(result.backupCount, "old version", packageWord, "will be moved.\n")
			} else {
				println(result.backupCount, "old version", packageWord, "have been moved.\n")
			}
		}
	}
}
