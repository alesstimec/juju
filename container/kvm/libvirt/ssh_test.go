package libvirt

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

// libvirtSSHSuite is gocheck boilerplate
type libvirtSSHSuite struct{}

// gocheck boilerplate
var _ = gc.Suite(libvirtSSHSuite{})

func (libvirtSSHSuite) TestKeepTheImports(c *gc.C) {
	c.Assert(true, jc.IsTrue)
}
