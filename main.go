package main

import (
	"errors"
	"os"
	"strings"
)

func main() {
	showVersion()

	argCount := len(os.Args)
	if argCount == 1 {
		// argCount 为 1，说明没有参数。args[0] 是程序名。
		showHelp()
		os.Exit(0)
	}

	// 此时 argCount 已经肯定大于 1 了。这是调用 parseArguments() 的要求。
	c, err := parseArguments(os.Args)
	if err != nil {
		showError(err)
		showHelp()
		os.Exit(1)
	}

	showStarting(c)

	cleanInfo, err := Clean(c.setupPath, c.listFilename, showObsoletePackage, c.showOnly, c.needDirStat)
	if err != nil {
		showError(err)
		os.Exit(2)
	}

	showCleanResult(c, cleanInfo)
	os.Exit(0)
}

type cmdArg struct {
	setupPath    string // Visual Studio Setup 的路径。
	listFilename string // 包列表文件的路径。
	showOnly     bool   // 是否只显示，而不实际处理过期的包。
	needDirStat  bool   // 是否需要统计目录数量、文件大小。
}

func parseOption(option string) (showOnly bool, needDirStat bool, hasOption bool, err error) {
	// 默认值。
	showOnly = true
	needDirStat = true
	hasOption = false
	err = nil

	// 看下面的 case 分支，判断的字符串长度都是 3。
	if len(option) != 3 {
		return
	}

	hasOption = true
	option = strings.ToLower(option)

	// '-' 后是两个 bool 值缩写。前面是 showOnly，后面是 needDirStat。
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
		err = errors.New("invalid option '" + option + "'")
	}

	return
}

// parseArguments 解析命令行参数。
// 不处理无参数情况，此时 args 的长度必需大于 1。
func parseArguments(args []string) (cmd *cmdArg, err error) {
	var hasOption bool
	cmd = &cmdArg{
		listFilename: "", // 显示地初始化。
	}
	if cmd.showOnly, cmd.needDirStat, hasOption, err = parseOption(args[1]); err != nil {
		return
	}

	n := len(args)
	if hasOption {
		if n > 4 {
			// 提供 option 时，参数最多为：option setupPath listFilename
			err = errors.New("too many arguments")
		} else if n == 2 {
			// 提供 option 时，参数最少为：option setupPath。
			// n == 2 说明提供了 option，但是没有提供 setupPath。
			err = errors.New("visual studio setup path must be specified")
		} else { // n == 3 || n == 4
			cmd.setupPath = os.Args[2]
			if n == 4 {
				cmd.listFilename = os.Args[3]
			}
		}
	} else { // option 使用了默认值。
		if n > 3 {
			// 未提供 option 时，参数最多为：setupPath listFilename
			err = errors.New("too many arguments")
		} else { // n == 2 || n == 3
			cmd.setupPath = os.Args[1]
			if n == 3 {
				cmd.listFilename = os.Args[2]
			}
		}
	}

	return
}
