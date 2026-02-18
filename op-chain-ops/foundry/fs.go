package foundry

import (
	"io/fs"
	"os"

	"github.com/tokamak-network/tokamak-thanos/op-chain-ops/srcmap"
)

// StatDirFs is an fs.FS that also supports Stat and ReadDir.
type StatDirFs interface {
	fs.FS
	fs.StatFS
	fs.ReadDirFS
}

// OpenArtifactsDir opens a directory as a StatDirFs for reading foundry artifacts.
func OpenArtifactsDir(dir string) StatDirFs {
	return os.DirFS(dir).(StatDirFs)
}

// NewSourceMapFS returns an FS for reading source maps from foundry output.
func NewSourceMapFS(fsys fs.FS) *SourceMapFS {
	return &SourceMapFS{fsys: fsys}
}

// SourceMapFS provides source-map access from a foundry artifacts FS.
type SourceMapFS struct {
	fsys fs.FS
}

// SourceMap reads the source map for a given artifact.
func (s *SourceMapFS) SourceMap(artifact *Artifact, contractName string) (*srcmap.SourceMap, error) {
	// Stub: actual source map parsing not implemented for old geth compat
	return nil, nil
}
