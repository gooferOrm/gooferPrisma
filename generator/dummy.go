//go:build tools
// +build tools

// See https://github.com/golang/go/issues/26366
package generator

import (
	_ "github.com/gooferOrm/goofer/generator/templates"
	_ "github.com/gooferOrm/goofer/generator/templates/actions"
)
