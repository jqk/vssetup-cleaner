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
	white.Print("visual studio setup file cleaner (")
	yellow.Print("vssc")
	white.Println(") 1.0.0, 2023-08-25\n")
}

func showHelp() {
	yellow.Println("Usage:")
	yellow.Println("  vssc [option] <path/to/visual studio setup> [path/to/package list]")
	white.Println("       clear obsolete Visual Studio Setup packages")
	yellow.Println("\nCommand:")
	white.Println("  The first character of the option determines whether to only show the result or execute the real action.")
	white.Println("  The second one defines whether to statistic the directory of obsolete packages or not.")
	white.Println("      't' is true, 'f' is false.\n")

	yellow.Print("  -tt: default option, can be omitted. ")
	white.Println("show only and statistic the result.")
	yellow.Print("  -tf: ")
	white.Println("show only but not statistic the result.")
	yellow.Print("  -ft: ")
	white.Println("clean and statistic the result.")
	yellow.Print("  -ff: ")
	white.Println("clean but not to statistic the result.")

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
