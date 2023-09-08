package main

import (
	"github.com/gookit/color"
	"github.com/jqk/futool4go/common"
)

var (
	blue   color.Style = color.New(color.LightBlue)
	green  color.Style = color.New(color.LightGreen)
	white  color.Style = color.New(color.White)
	yellow color.Style = color.New(color.LightYellow)
	red    color.Style = color.New(color.LightRed)
)

func showVersion() {
	white.Println()
	white.Println("Copyright (c) 1999-2023 Not a dream Co., Ltd.")
	white.Print("Visual Studio setup file cleaner (")
	yellow.Print("vssc")
	white.Println(") 1.0.1, 2023-09-08\n")
}

func showHelp() {
	yellow.Println("Usage:")
	yellow.Println("  vssc [option] <visual studio setup path> [package list file]")
	white.Println("       clear obsolete Visual Studio Setup packages")
	yellow.Println("\nOption:")
	white.Println("      The first character of the option determines whether to only show the result or execute the real action.")
	white.Println("      The second one defines whether to statistic the directory of obsolete packages or not.")
	white.Println("      't' is true, 'f' is false.\n")

	yellow.Print("  -tt: default option, can be omitted. ")
	white.Println("Display obsolete packages and their summary.")
	yellow.Print("  -tf: ")
	white.Println("Display obsolete packages only.")
	yellow.Print("  -ft: ")
	white.Println("Clean obsolete packages and display their summary.")
	yellow.Print("  -ff: ")
	white.Println("Clean obsolete packages only.")

	yellow.Println()
	yellow.Println("  otherwise: show this help.")
	yellow.Println()
}

func showError(err error) {
	red.Println("---------- Error! ----------")
	red.Println(err)
	red.Println("----------------------------\n")
}

func showStarting(c *cmdArg) {
	green.Println("Cleaning '" + c.setupPath + "'\n")
}

func showCleanResult(c *cmdArg, cleanInfo *CleanInfo) {
	foundPackages := len(cleanInfo.Packages) > 0

	if foundPackages {
		green.Println()
	}

	green.Println("-------------------")
	green.Println("Show only         :", c.showOnly)

	if !c.showOnly && foundPackages {
		green.Println("Backup directory  :", cleanInfo.BackupPath)
	}
	if c.needDirStat {
		dirCount, fileCount, size := cleanInfo.Stat()
		green.Println("Total directories :", dirCount)
		green.Println("Total files       :", fileCount)
		green.Println("Total size        :", common.ToSizeString(size))
	}

	green.Println("Obsolete packages :", len(cleanInfo.Packages))
	green.Println("Elapsed time      :", cleanInfo.ElapsedTime.String())
	green.Println("-------------------\n")
}

func showObsoletePackage(packageInfo *PackageInfo) error {
	blue.Println(packageInfo.Dir)
	return nil
}
