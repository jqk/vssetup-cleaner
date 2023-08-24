package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOldPackagesByDirVersion(t *testing.T) {
	pkgs, err := FindOldPackagesByDirVersion("test-data")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(pkgs))
	assert.Equal(t, "Microsoft.VisualStudio.Web.Mvc,version=16.10.115.10959", pkgs[0].Dir)
	assert.Equal(t, "Microsoft.VisualStudio.Web.Mvc,version=16.11.115.10951", pkgs[1].Dir)
	assert.Equal(t, "Microsoft.VisualStudio.Web.Mvc,", pkgs[0].Name)
	assert.Equal(t, "Microsoft.VisualStudio.Web.Mvc,", pkgs[1].Name)
	assert.Equal(t, "16.10.115.10959", pkgs[0].Version)
	assert.Equal(t, "16.11.115.10951", pkgs[1].Version)
}
