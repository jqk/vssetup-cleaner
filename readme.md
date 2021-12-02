# Visual Studio setup cleaner

[中文版](readme_cn.md)

## 1. What to solve

We have 2 ways to install Visual Studio:

- online installation.
- download the selected components and then install it.

We are discussing the second scenario.

Every major version of Visual Studio will have some minor version updates. Each update changes only a small number of installation files. It's a waste of time to download all the components each time when the update is ready.

So, I download the newest installation file first, and then run:

```bash
vs_Community.exe --layout "f:\Software\vs2019-community"
```

`f:\Software\vs2019-community` is the directory where the installation files of previous version are placed. You should set it with your value.

Running above command will only download the changed components. It's faster than downloading all the components.

The problem is it will not remove old version components. So the directory such as `f:\Software\vs2019-community` will go larger and larger.

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

Many components are just different versions.

This project is build for cleaning the old component directories.

- It will first create a directory named `yyyy-mm-dd_ HH-MM-ss` according to current time for backup.
- Move old version component dirs to the backup directory.
- After checking, you can **manually delete** the bauckup directory.

## 2. How to use

- Clone the source code.
- Build a executable and run it.

Or, just run it with GO:

```bash
go run main.go --path YOUR_PATH_TO_VS_SETUP_DIR
```

If you add the `-- showonly` parameter, only information will be displayed, and the backup operation will not be performed.

The project is very simple. There're less than `200 lines` in main.go.

Run `vs_Community.exe --layout "f:\Software\vs2019-community"` again to verify everything is ok. It should not download any component. Then we can remove the buckup directory to release the space.

## 3. What next

There's one problem left. There are a lot of components like below:

```text
Microsoft.AspNetCore.SharedFramework.2.1.2.1.29,version=16.11.31603.221,chip=x64
Microsoft.AspNetCore.SharedFramework.2.1.2.1.30,version=16.11.31701.289,chip=x64
Microsoft.AspNetCore.SharedFramework.2.2.2.2.8,version=16.10.31205.180,chip=x64
```

In addition to the component version, there are versions of other framework.

Backing up these directories is as easy as backing up component version directories. It only takes time to repeat the test to see if the directories we move are still used by the latest version.

It's realy a waste of time. Maybe, maybe, maybe, I'll do it in the future.

BTW, the program released 8.7GB for Visual Studio 2019 community 16.7 setup files.

**Enjoy!**
