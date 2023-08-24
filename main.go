package main

import (
	"errors"
	"os"
	"strings"
)

func main() {
	argCount := len(os.Args)
	if argCount == 1 {
		showHelp()
		os.Exit(0)
	}

	c, err := parseArguments(os.Args)
	if err != nil {
		showError(err)
		showHelp()
		os.Exit(0)
	}

	cleanInfo, err := Clean(c.setupPath, c.listFilename, showObsoletePackage, c.showOnly, c.needDirStat)
	if err != nil {
		showError(err)
		os.Exit(1)
	}

	showCleanInfo(cleanInfo)
	os.Exit(0)
}

type cmdArg struct {
	setupPath    string
	listFilename string
	showOnly     bool
	needDirStat  bool
}

func parseOption(option string) (showOnly bool, needDirStat bool, hasOption bool, err error) {
	// 默认值。
	showOnly = true
	needDirStat = true
	hasOption = false
	err = nil

	if len(option) != 3 {
		return
	}

	hasOption = true
	option = strings.ToLower(option)

	switch option {
	case "-tt":
		// 使用默认值。
	case "-tf":
		needDirStat = false
	case "-ft":
		showOnly = false
	case "-ff":
		showOnly = false
		needDirStat = false
	default:
		err = errors.New("invalid option")
	}

	return
}

func parseArguments(args []string) (cmd *cmdArg, err error) {
	n := len(args)
	if n > 4 {
		return nil, errors.New("too many arguments")
	}

	var hasOption bool
	cmd = &cmdArg{
		listFilename: "", // 显示地初始化。
	}
	if cmd.showOnly, cmd.needDirStat, hasOption, err = parseOption(args[1]); err != nil {
		return
	}

	if hasOption {
		if n == 2 {
			return nil, errors.New("visual studio setup path must be specified")
		}

		cmd.setupPath = os.Args[2] // n == 3 || n == 4
		if n == 4 {
			cmd.listFilename = os.Args[3]
		}
	} else { // option 使用了默认值。
		if n == 4 {
			return nil, errors.New("too many arguments")
		}

		cmd.setupPath = os.Args[1] // n == 2 || n == 3
		if n == 3 {
			cmd.listFilename = os.Args[2]
		}
	}

	return
}
