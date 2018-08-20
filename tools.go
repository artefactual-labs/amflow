// +build tools

// This file is used to track the version of the build tools in our module's `go.mod` file.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module for more details.

package tools

import (
	_ "github.com/goadesign/goa/goagen"
	_ "github.com/gobuffalo/packr/v2/packr2"
)
