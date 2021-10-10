package cache

import "os"

type TileCacheLocal struct {
	path string
	perm os.FileMode
}

// NewTileCache ...
func NewTileCache(rootPath string, perm os.FileMode) *TileCacheLocal {
	return &TileCacheLocal{
		path: rootPath,
		perm: perm,
	}
}

// Path ...
func (c *TileCacheLocal) Path() string {
	return c.path
}

// Perm ...
func (c *TileCacheLocal) Perm() os.FileMode {
	return c.perm
}
