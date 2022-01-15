# Visual Studio setup cleaner

[English](readme.md)

## 一、 解决什么问题

我们有两种安装 Visual Studio 的方法：

- 在线安装。
- 将安装文件下载到本地，从本地安装。

本文说明**本地安装场景**。

每一个大版本的 Visual Studio，都会有多个小版本的升级。每次升级都只是升级部分组件而非全部安装文件。这样一来，每次升级都下载整个安装包就太浪费时间了。

所以，我们应该每次只下载更新的部分，此处**以社区版为例**，运行命令与初次下载安装包时的命令相同：

```bash
vs_community.exe --layout "f:\Software\vs2019-community"
```

- `vs_community.exe`不需要最新，因为它会自动下载最新的`vs_layout.exe`并运行。
- `f:\Software\vs2019-community`安装包所在目录，应该根据实际情况填写。

以上命令即可只下载更新了的组件，比下载整个版本快很多。

但这样操作**不会删除旧版的、不再使用的组件目录**，导致`f:\Software\vs2019-community`越来越大。

本程序即删除这些不再使用的组件。这些组件被移动到根据当前时间生成的格式为`_yyyy-mm-dd_HH-MM-ss`的备份目录中。检查无误后可以手动删除。

## 二、 如何使用

### 2.1 Visual Studio setup cleaner

- 下载代码。
- 编译并执行：`vssetup-cleaner.exe -h`显示程序可使用的参数。

```cmd
vssetup-cleaner.exe -h
Usage of vssetup-cleaner.exe:
  -pkg string
        package list file
  -show
        Show action only (default true)
  -vs string
        Visual Studio's setup path (default ".")
```

当忽略pkg或将其设置为空字符串时，将使用[版本1的操作逻辑](readme_cn_v1.md)。

### 2.2 执行过程

执行过程如下：

1. 执行`vs_community.exe --layout "f:\Software\vs2019-community"`下载本地安装包。
2. 执行`vs_setup.exe`进行安装。
3. 通常运行`VS2019`会自动发现新的小版本。**此时不更新**。
4. 执行`vs_community.exe --layout "f:\Software\vs2019-community"`下载新版本发生变动的组件包。
5. 如果中断或有其它错误发生，反复执行以上命令直到所有新的组件包全部下载完毕。
6. 再次执行`vs_community.exe --layout "f:\Software\vs2019-community"`会检查一遍组件包是否完整。所有组件都被列在命令行窗口中。通过复制、粘贴的方式，将这总分内容保存到文本文件中，例如`e:\temp\list.txt`。该文件即为`vssetup-cleaner.exe`的`pkg`参数所需的`package list file`。
7. 执行`vssetup-cleaner.exe -pkg="e:\temp\list.txt" -vs="f:\Software\vs2019-community" -show=false`会将`f:\Software\vs2019-community`中的目录与`e:\temp\list.txt`中的组件名作对比，将过期不用的组件目录移动到`f:\Software\vs2019-community`下的`_yyyy-mm-dd_HH-MM-ss`备份目录中。
8. 再次执行`vs_community.exe --layout "f:\Software\vs2019-community"`会检查一遍组件包是否完整。如果没有问题，可以手工删除备份目录。
9. 执行`vs_setup.exe`安装更新即可。

对 Visual Studio 2019 community 16.9 安装包运行此程序后，释放了 6.4 GB 的空间。

**Enjoy!**
