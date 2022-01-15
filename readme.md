# Visual Studio setup cleaner

[中文版](readme_cn.md)

## 1. What to solve

We have 2 ways to install Visual Studio:

- online installation.
- download the selected components and then install it.

We are discussing the second scenario.

Every major version of Visual Studio will have some minor version updates. Each update changes only a small number of installation files. It's a waste of time to download all the components each time when the update is ready.

So, download the newest installation file first, and then run:

```bash
vs_community.exe --layout "f:\Software\vs2019-community"
```

`f:\Software\vs2019-community` is the directory where the installation files of previous version are placed. You should set it with your value.

Running above command will only download the changed components. It's faster than downloading all the packages.

The problem is it will not remove old version components. So the directory such as `f:\Software\vs2019-community` will go larger and larger.

This project is build for cleaning the old component directories.

- It will create a backup directory named `_yyyy-mm-dd_HH-MM-ss` according to current time for backup.
- Move old version component directories to the new created backup directory.
- Then you can check it and **manually delete** the backup directory.

## 2. How to use

### 2.1 Visual Studio setup cleaner

- Clone the source code.
- Build an executable and run it.

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

[Version 1 process logic](readme_v1.md) is used when PKG is missing or set to an empty string.

### 2.2 Update step by step

Run `vs_Community.exe --layout "f:\Software\vs2019-community"` again to verify everything is ok. It should not download any component. Then we can remove the backup directory to release the space.

1. Download local installation files by running `vs_community.exe --layout "f:\Software\vs2019-community"`.
2. Install `VS2019` by running `vs_setup.exe`.
3. `VS2019` will sent out notification when there's update available. **Don't update now.**.
4. Download new packages by running `vs_community.exe --layout "f:\Software\vs2019-community"`.
5. Repeat the above operation until download is successful.
6. Run `vs_community.exe --layout "f:\Software\vs2019-community"` again. It will check the local packages and log all package names to command windows. Save the output to a file, for example, `e:\temp\list.txt` which will be used as the `package list file` for `vssetup-cleaner.exe`.
7. Run `vssetup-cleaner.exe -pkg="e:\temp\list.txt" -vs="f:\Software\vs2019-community" -show=false`.It'll move old packages to the backup dir。
8. Check the installation by running `vs_community.exe --layout "f:\Software\vs2019-community"` again.
9. Update `VS2019` by running `vs_setup.exe`. Then delete backup dir.

The program released 6.4 GB for Visual Studio 2019 community 16.9 setup files.

**Enjoy!**
