# Visual Studio setup cleaner

[English](readme_cn_v1.md)，[V2 版说明](readme_cn.md)

## 一、 解决什么问题

我们有两种安装 Visual Studio 的方法：

- 在线安装。
- 将安装文件下载到本地，从本地安装。

本文说明本地安装场景。

每一个大版本的 Visual Studio，都会有多个小版本的升级。每次升级都只是升级部分组件而非全部安装文件。这样一来，每次升级都下载整个安装包就太浪费时间了。

所以，我们应该每次只下载更新的部分。首先下载最新的安装文件，此处以社区片为例，是`vs_Community.exe`，约`1.4M`。然后运行：

```bash
vs_Community.exe --layout "f:\Software\vs2019-community"
```

该命令与初次下载安装包时的命令相同，就是这个`vs_Community.exe`要最新的。`f:\Software\vs2019-community`安装包所在目录，应该根据实际情况填写。

以上命令即可只下载更新了的组件，比下载整个版本快很多。

但这样操作不会删除旧版的，不再使用的组件安装目录，导致`f:\Software\vs2019-community`越来越大。下面是组件包目录名示例。

```text
2021/12/01  11:52    <DIR>          AndroidImage_ARM_API25,version=21.0.0.3
2021/12/01  11:41    <DIR>          AndroidImage_ARM_API25,version=21.0.0.5
2021/12/01  14:25    <DIR>          Microsoft.DiagnosticsHub.DatabaseTool.Targeted,version=16.11.31402.2,chip=x64
2021/12/01  14:25    <DIR>          Microsoft.DiagnosticsHub.DatabaseTool.Targeted,version=16.11.31402.2,chip=x86
2021/12/01  14:25    <DIR>          Microsoft.DiagnosticsHub.DatabaseTool.Targeted,version=16.11.31822.4,chip=x64
2021/12/01  14:25    <DIR>          Microsoft.DiagnosticsHub.DatabaseTool.Targeted,version=16.11.31822.4,chip=x86
2021/11/30  15:27    <DIR>          Microsoft.VisualStudio.Net.Eula.Resources,version=16.0.28315.86,language=en-US
2021/12/01  15:09    <DIR>          Microsoft.VisualStudio.Setup.Configuration,version=2.11.47.9733
2021/12/01  15:09    <DIR>          Microsoft.VisualStudio.Setup.Configuration,version=2.7.3111.17308
```

这些组件目录名，只有版本，即 version 部份不同。

本项目即用来清理这些目录：

- 首先以当前时间，建立包为`yyyy-mm-dd_HH-MM-ss`的备份目录。
- 然后将旧版组件目录移动到该备份目录中。
- 确保无误后，人工删除备份目录以释放磁盘空间。

## 二、 如何使用

- 下载代码。
- 编译并执行。

或者，下载代码后直接执行以下命令：

```bash
go run main.go --path 你自己的VS安装包目录
```

加上`--showonly`参数，将只显示信息，而不会执行备份操作。

这个项目很简单，只有一个代码文件 main.go，包含注释在内也才`200行`，有不明白稍看一下即可。

执行完以上命令后再次运行`vs_Community.exe --layout "f:\Software\vs2019-community"`可以检查是否未删除必要的安装包。如果没下载新的组件，就说明是安全的，可以删除备份目录了。

## 三、 3. What next

有一个遗留问题，看下面的目录名：

```text
Microsoft.AspNetCore.SharedFramework.2.1.2.1.29,version=16.11.31603.221,chip=x64
Microsoft.AspNetCore.SharedFramework.2.1.2.1.30,version=16.11.31701.289,chip=x64
Microsoft.AspNetCore.SharedFramework.2.2.2.2.8,version=16.10.31205.180,chip=x64
Microsoft.AspNetCore.SharedFramework.3.0.3.0.3.x64,version=16.10.31205.180,chip=x64
.....
.....
Microsoft.AspNetCore.SharedFramework.3.1.3.1.18-servicing.21365.4.x86,version=16.11.31603.221
Microsoft.AspNetCore.SharedFramework.3.1.3.1.19-servicing.21417.13.x86,version=16.11.31701.289
Microsoft.AspNetCore.SharedFramework.3.1.3.1.20-servicing.21472.42.x86,version=16.11.31729.444
Microsoft.AspNetCore.SharedFramework.3.1.3.1.21-servicing.21523.9.x86,version=16.11.31828.110
Microsoft.AspNetCore.SharedFramework.3.1.3.1.22-servicing.21579.4.x86,version=16.11.32002.110
.....
Microsoft.AspNetCore.SharedFramework.5.0.5.0.9-servicing.21365.3.x86,version=16.11.31603.221
Microsoft.AspNetCore.SharedFramework.5.0.5.0.10-servicing.21410.22.x86,version=16.11.31701.289
Microsoft.AspNetCore.SharedFramework.5.0.5.0.11-servicing.21476.5.x86,version=16.11.31729.444
Microsoft.AspNetCore.SharedFramework.5.0.5.0.12-servicing.21524.1.x86,version=16.11.31828.110
Microsoft.AspNetCore.SharedFramework.5.0.5.0.13-servicing.21572.2.x86,version=16.11.32002.110
```

除我们处理的由`version`指定的版本外，还有针对其它框架的版本。

以同样方式处理这些版本很简单，但需要反复验证是否破坏了安装包。这很耗费时间，以后可能会做。

对 Visual Studio 2019 community 16.7 安装包运行此程序后，释放了 8.7GB 的空间。

**Enjoy!**
