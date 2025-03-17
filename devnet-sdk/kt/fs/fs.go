package fs

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// EnclaveContextIface abstracts the EnclaveContext for testing
type EnclaveContextIface interface {
	GetAllFilesArtifactNamesAndUuids(ctx context.Context) ([]*kurtosis_core_rpc_api_bindings.FilesArtifactNameAndUuid, error)
	DownloadFilesArtifact(ctx context.Context, name string) ([]byte, error)
	UploadFiles(pathToUpload string, artifactName string) (services.FilesArtifactUUID, services.FileArtifactName, error)
}

type EnclaveFS struct {
	enclaveCtx EnclaveContextIface
}

func NewEnclaveFS(ctx context.Context, enclave string) (*EnclaveFS, error) {
	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		return nil, err
	}

	enclaveCtx, err := kurtosisCtx.GetEnclaveContext(ctx, enclave)
	if err != nil {
		return nil, err
	}

	return &EnclaveFS{enclaveCtx: enclaveCtx}, nil
}

// NewEnclaveFSWithContext creates an EnclaveFS with a provided context (useful for testing)
func NewEnclaveFSWithContext(ctx EnclaveContextIface) *EnclaveFS {
	return &EnclaveFS{enclaveCtx: ctx}
}

type Artifact struct {
	rawData []byte
	reader  *tar.Reader
}

func (fs *EnclaveFS) GetAllArtifactNames(ctx context.Context) ([]string, error) {
	artifacts, err := fs.enclaveCtx.GetAllFilesArtifactNamesAndUuids(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(artifacts))
	for i, artifact := range artifacts {
		names[i] = artifact.GetFileName()
	}

	return names, nil
}

func (fs *EnclaveFS) GetArtifact(ctx context.Context, name string) (*Artifact, error) {
	artifact, err := fs.enclaveCtx.DownloadFilesArtifact(ctx, name)
	if err != nil {
		return nil, err
	}

	// Store the raw data
	buffer := bytes.NewBuffer(artifact)
	zipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(zipReader)
	return &Artifact{
		rawData: artifact,
		reader:  tarReader,
	}, nil
}

func (a *Artifact) newReader() (*tar.Reader, error) {
	buffer := bytes.NewBuffer(a.rawData)
	zipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	return tar.NewReader(zipReader), nil
}

func (a *Artifact) Download(path string) error {
	// Create a new reader for this operation
	reader, err := a.newReader()
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		fpath := filepath.Join(path, filepath.Clean(header.Name))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
		case tar.TypeReg:
			// Create parent directories if they don't exist
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", fpath, err)
			}

			// Create the file
			f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", fpath, err)
			}

			// Copy contents from tar reader to file
			if _, err := io.Copy(f, reader); err != nil {
				f.Close()
				return fmt.Errorf("failed to write contents to %s: %w", fpath, err)
			}
			f.Close()
		default:
			return fmt.Errorf("unsupported file type %d for %s", header.Typeflag, header.Name)
		}
	}
}

func (a *Artifact) ExtractFiles(writers ...*ArtifactFileWriter) error {
	// Create a new reader for this operation
	reader, err := a.newReader()
	if err != nil {
		return fmt.Errorf("failed to create reader: %w", err)
	}

	paths := make(map[string]io.Writer)
	for _, writer := range writers {
		canonicalPath := filepath.Clean(writer.path)
		paths[canonicalPath] = writer.writer
	}

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		headerPath := filepath.Clean(header.Name)
		if _, ok := paths[headerPath]; !ok {
			continue
		}

		writer := paths[headerPath]
		_, err = io.Copy(writer, reader)
		if err != nil {
			return fmt.Errorf("failed to copy content: %w", err)
		}
	}

	return nil
}

func (fs *EnclaveFS) PutArtifact(ctx context.Context, name string, readers ...*ArtifactFileReader) error {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "artifact-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up temp dir when we're done

	// Process each reader
	for _, reader := range readers {
		// Create the full path in the temp directory
		fullPath := filepath.Join(tempDir, reader.path)

		// Ensure the parent directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}

		// Create the file
		file, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		// Copy the content
		_, err = io.Copy(file, reader.reader)
		file.Close() // Close file after writing
		if err != nil {
			return err
		}
	}

	// Upload the directory to Kurtosis
	_, _, err = fs.enclaveCtx.UploadFiles(tempDir, name)
	return err
}

type ArtifactFileReader struct {
	path   string
	reader io.Reader
}

func NewArtifactFileReader(path string, reader io.Reader) *ArtifactFileReader {
	return &ArtifactFileReader{
		path:   path,
		reader: reader,
	}
}

type ArtifactFileWriter struct {
	path   string
	writer io.Writer
}

func NewArtifactFileWriter(path string, writer io.Writer) *ArtifactFileWriter {
	return &ArtifactFileWriter{
		path:   path,
		writer: writer,
	}
}
