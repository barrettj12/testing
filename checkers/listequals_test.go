// Copyright 2024 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package checkers_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type listEqualsSuite struct{}

var _ = gc.Suite(&listEqualsSuite{})

func (s *listEqualsSuite) TestFoo(c *gc.C) {
	c.Check(
		[]rune("ABDEZGHILJK"),
		jc.ListEquals,
		[]rune("ABCDEFGHIJK"),
	)
}
