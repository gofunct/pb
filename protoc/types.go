package protoc

import (
	"github.com/gofunct/pb/protoc/build"
)

// RootDir represents a project root directory.
type RootDir struct {
	build.Path
}

// BinDir returns the directory path contains executable binaries.
func (d *RootDir) BinDir() build.Path {
	return d.Join("bin")
}
