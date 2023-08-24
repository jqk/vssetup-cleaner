package main

import (
	"testing"

	"github.com/jqk/futool4go/common"

	"github.com/stretchr/testify/assert"
)

func TestGetPackageInfo(t *testing.T) {
	// only have version number.
	pkg := getPackageInfo("Microsoft.Net.Core.SDK.MSBuildExtensions,version=16.11.31603.221")
	assert.Equal(t, "Microsoft.Net.Core.SDK.MSBuildExtensions,", pkg.Name)
	assert.Equal(t, "16.11.31603.221", pkg.Version)

	// the package name has version number.
	pkg = getPackageInfo("Win10SDK_10.0.16299,version=10.0.16299.5")
	assert.Equal(t, "Win10SDK_10.0.16299,", pkg.Name)
	assert.Equal(t, "10.0.16299.5", pkg.Version)

	// have version number, chip set and language.
	pkg = getPackageInfo("Microsoft.VisualStudio.Vsto.Runtime.Resources,version=16.0.28315.86,chip=x64,language=zh-CN")
	assert.Equal(t, "Microsoft.VisualStudio.Vsto.Runtime.Resources,chip=x64,language=zh-CN,", pkg.Name)
	assert.Equal(t, "16.0.28315.86", pkg.Version)

	// have version number and chip set.
	pkg = getPackageInfo("Microsoft.VisualStudio.Vsto.Runtime,version=16.0.28315.86,chip=x64")
	assert.Equal(t, "Microsoft.VisualStudio.Vsto.Runtime,chip=x64,", pkg.Name)
	assert.Equal(t, "16.0.28315.86", pkg.Version)

	// have version number, chip set and language in neutral.
	pkg = getPackageInfo("Microsoft.Build,version=16.11.2.2150704,chip=neutral,language=neutral")
	assert.Equal(t, "Microsoft.Build,chip=neutral,language=neutral,", pkg.Name)
	assert.Equal(t, "16.11.2.2150704", pkg.Version)

	// have version number and language.
	pkg = getPackageInfo("Microsoft.DiagnosticsHub.DatabaseTool.Resources,version=16.11.31901.3,language=zh-CN")
	assert.Equal(t, "Microsoft.DiagnosticsHub.DatabaseTool.Resources,language=zh-CN,", pkg.Name)
	assert.Equal(t, "16.11.31901.3", pkg.Version)
}

func TestCompareVersions(t *testing.T) {
	assert.Equal(t, -1, common.CompareVersions("16.10.31603.221", "16.11.31603.221"))
	assert.Equal(t, 0, common.CompareVersions("16.11.31603.221", "16.11.31603.221"))
	assert.Equal(t, 1, common.CompareVersions("16.11.31604.221", "16.11.31603.221"))
}
