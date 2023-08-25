package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var packages []*PackageInfo = make([]*PackageInfo, 0, 2)

func TestClean(t *testing.T) {
	packages = packages[0:0]
	info, err := Clean("test-data", "", packageHandler, true, true)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, len(info.Packages), len(packages))

	assert.Equal(t, 1, info.Packages[0].FileCount)
	assert.Equal(t, int64(52), info.Packages[0].Size)
	assert.Equal(t, 0, info.Packages[1].FileCount)
	assert.Equal(t, int64(0), info.Packages[1].Size)

	packages = packages[0:0]
	info, err = Clean("test-data", "test-data/certificates/packageList.txt", packageHandler, true, false)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, len(info.Packages), len(packages))
}

func packageHandler(packageInfo *PackageInfo) error {
	packages = append(packages, packageInfo)
	return nil
}
