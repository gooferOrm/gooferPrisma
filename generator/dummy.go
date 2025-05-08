//go:build tools
// +build tools

// See https://github.com/golang/go/issues/26366
package generator

import (
	_ "github.com/tacherasasi/goofer/generator/templates"
	_ "github.com/tacherasasi/goofer/generator/templates/actions"
)
