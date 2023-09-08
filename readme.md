# Visual Studio setup cleaner

[中文版](readme_cn.md)

## 1. What is it

Visual Studio setup cleaner (vssc) is a command line tool for cleaning up outdated Visual Studio installation packages.

It has only been tested on Windows 10/11 for Visual Studio 2019 Community. In theory it should be able to handle all versions from 2017 to 2022.

Since it is written in Go and does not use any OS-specific features, it should also be able to run on Mac.

## 2. What to solve

There are two ways to install Visual Studio:

- Online installation.
- Download the installation files locally and install from local.

This article describes the **local installation** scenario.

Each major version of Visual Studio will have multiple minor version upgrades. Each upgrade only upgrades some components rather than all installation files. Downloading the entire installation package each time is wasteful.

Therefore, we should only download the updated parts each time. **Take the community edition as an example**, the command to run is the same as when downloading the installation package for the first time:

```bash
vs_community.exe --layout "f:\Software\vs2019-community"
```

- `vs_community.exe` does not need to be the latest, it will automatically download the latest `vs_layout.exe` and run it.
- `f:\Software\vs2019-community` is the directory where the installation package is located, should be filled in according to the actual situation.

The above command will only download updated components, which is much faster than downloading the entire version.
However, this operation **does not delete the directories of old, no longer used components**, causing f:\Software\vs2019-community to grow larger and larger.

`vssc` finds these no longer used components, and moves them to a backup directory named `_yyyy-mm-dd_HH-MM-ss` generated according to the current time. After verification, they can be manually deleted.

## 3. Installation

There are three ways to install:

1. Build the source code:
   - `git clone https://github.com/jqk/vssetup-cleaner.git`
   - `cd vssetup-cleaner`
   - `go build`. It will generate `vssc.exe` in Windows.
2. You can download the package from <https://github.com/jqk/vssetup-cleaner/releases>, unzip and run directly.
3. Use [scoop](https://github.com/ScoopInstaller/Scoop) in Windows. After installing scoop, execute:
   - `scoop bucket add ajqk https://github/jqk/scoopbucket`
   - `scoop install vssetup-cleaner`

## 4. Usage

### 4.1 vssc Command

```text {.line-numbers}
$ vssc

Copyright (c) 1999-2023 Not a dream Co., Ltd.
Visual Studio setup file cleaner (vssc) 1.0.1, 2023-09-08

Usage:
  vssc [option] <visual studio setup path> [package list file]
        clear obsolete Visual Studio Setup packages

Option:
      The first character of the option determines whether to only show the result or execute the real action.
      The second one defines whether to statistic the directory of obsolete packages or not.
      't' is true, 'f' is false.

  -tt: Default option, can be omitted. Display obsolete packages and their summary.
  -tf: Display obsolete packages only.
  -ft: Clean obsolete packages and display their summary.
  -ff: Clean obsolete packages only.

  otherwise: show this help.
```

The `visual studio setup path` must be specified for `vssc` to work. The `package list file` is optional. Specifying this parameter causes `vssc` to clean up with different logic:

- If provided: `vssc` will compare the directories under `visual studio setup path` against the `package list file`. Directories not listed in the file will be considered outdated and removed.
- If omitted: `vssc` will analyze the directory names under `visual studio setup path` to extract package name and version. Older versions will be removed.

The two methods can validate each other and produce the same results in testing. Using the second method without the file parameter is simpler.

The default is to use the `-tt` option, which only prints the outdated packages found. To actually clean up, use `-ft` or `-ff`.

> See below for how to get the package list file.

### 4.2 Update and Cleanup Process

The update and cleanup process is:

1. Running VS2019 normally will auto detect new minor versions. **_Do not update yet!_**
2. Run `vs_community.exe --layout "f:\Software\vs2019-community"` to download new component packages. Note the path in quotes is where VS is installed, update it accordingly.
3. If interrupted or errors occur, repeat the above command until all new packages are downloaded. Do not close the window.
4. Save the `package list file`. Ignore this step if not using this file. Copy and paste (`Ctrl + A`, `Ctrl + C` to copy from window, paste into a text file e.g. e:\temp\list.txt) to save the list. This is the `package list file` for `vssc`.
5. Run cleanup using `vssc` as described in the commands. Use `-ft` or `-ff` to actually execute cleanup.
6. Run `vs_community.exe --layout "f:\Software\vs2019-community` again to verify component packages. Delete the backup folder manually if no issues. Optional.
7. Run `vs_setup.exe` or start VS and install the update. Downloading will finish instantly. Offline update is supported.

**Enjoy!**
