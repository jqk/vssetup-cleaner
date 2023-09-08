# Visual Studio setup cleaner

[English](readme.md)

## 一、 Visual Studio Setup Cleaner

Visual Studio setup cleaner (vssc) 是清理 Visual Studo 过期安装包的命令行工具。

仅在 Windows 10/11 下针对 Visual Studio 2019 Community 测试过。理论上可以处理 2017 至 2022 的所有版本。

因为使用 go 编写，且未使用与操作系统相关的特有功能，所以应也可在 Mac 上运行。

## 二、 解决什么问题

我们有两种安装 Visual Studio 的方法：

- 在线安装。
- 将安装文件下载到本地，从本地安装。

本文说明**本地安装场景**。

每一个大版本的 Visual Studio，都会有多个小版本的升级。每次升级都只是升级部分组件而非全部安装文件。这样一来，每次升级都下载整个安装包就太浪费时间了。

所以，我们应该每次只下载更新的部分，此处**以社区版为例**，运行命令与初次下载安装包时的命令相同：

```bash
vs_community.exe --layout "f:\Software\vs2019-community"
```

- `vs_community.exe` 不需要最新，它会自动下载最新的 `vs_layout.exe` 并运行。
- `f:\Software\vs2019-community` 安装包所在目录，应该根据实际情况填写。

以上命令即可只下载更新了的组件，比下载整个版本快很多。

但这样操作**不会删除旧版的、不再使用的组件目录**，导致 `f:\Software\vs2019-community` 越来越大。

`vssc` 即查找这些不再使用的组件，并将这些组件移动到根据当前时间生成的格式为`_yyyy-mm-dd_HH-MM-ss`的备份目录中。检查无误后可以手动删除。

## 三、 安装

有两种方法进行安装：

1. 编译源码：
   - `git clone https://github.com/jqk/vssetup-cleaner.git`
   - `cd vssetup-cleaner`
   - `go build`。该命令在 Windows 下生成 `vssc.exe`。
2. 可以从 <https://github.com/jqk/vssetup-cleaner/releases> 下载程序包，解压后直接运行。
3. 在 Windows 下使用 [scoop](https://github.com/ScoopInstaller/Scoop) 。在安装完 scoop 后执行：
   - `scoop bucket add ajqk https://github/jqk/scoopbucket`
   - `scoop install vssc`

## 四、 使用

### 4.1 vssc 命令

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

必须指定 `visual studio setup path` 才能工作。而 `package list file` 是可选的，该参数将使 `vssc` 以不同逻辑执行清理操作：

- 提供该参数：对比 `visual studio setup path` 中的目录，存在于 `package list file` 则保留；不在其中的被认为是过期的，将被清除。
- 省略该参数：分析 `visual studio setup path` 中的目录名，得到包名和版本号。版本号旧的将被清除。

两种方式可相互验证，经测试结果相同。使用`省略该参数`的方式更简单易用。

默认使用 `-tt` 选项，仅给出查到的过期包信息。若要真正清理，请使用 `-ft` 或 `-ff` 选项。

> 获得 `package list file` 的方法见下面说明。

### 4.2 更新及清理过程

更新及清理过程如下：

1. 通常运行 `VS2019` 会自动发现新的小版本。**此时不要更新！**。
1. 执行 `vs_community.exe --layout "f:\Software\vs2019-community"` 下载新版本发生变动的组件包。注意引号内的是安装程序路径，应根据实际情况填写。
1. 如果中断或有其它错误发生，反复执行以上命令直到所有新的组件包全部下载完毕。注意不要关闭窗口。
1. 本步保存 `package list file`。如果不想使用该文件则请忽略本步。通过复制、粘贴（`Ctrl + A`，`Ctrl + C`从窗口中选择复制，再粘贴到文本文件中）的方式，将这总分内容保存到文本文件中，例如 `e:\temp\list.txt` 。该文件即为 `vssc` 的`package list file`。
1. 按前面的命令说明执行清理操作。如果需要实际执行，请使用 `-ft` 或 `-ff` 选项。
1. 再次执行 `vs_community.exe --layout "f:\Software\vs2019-community"` 会检查一遍组件包是否完整。如果没有问题，可以手工删除备份目录。本步可选。
1. 执行 `vs_setup.exe` 或启动 Visual Studio 后安装更新即可。此时仍然会提示下载更新包，但点击继续后，下载过程瞬间完成。此操作支持离线更新。

**Enjoy!**
