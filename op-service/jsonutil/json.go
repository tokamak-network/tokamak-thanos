package jsonutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tokamak-network/tokamak-thanos/op-service/ioutil"
)

func LoadJSON[X any](inputPath string) (*X, error) {
	if inputPath == "" {
		return nil, errors.New("no path specified")
	}
	var f io.ReadCloser
	f, err := ioutil.OpenDecompressed(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", inputPath, err)
	}
	defer f.Close()
	var state X
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&state); err != nil {
		return nil, fmt.Errorf("failed to decode file %q: %w", inputPath, err)
	}
	// We are only expecting 1 JSON object - confirm there is no trailing data
	if _, err := decoder.Token(); err != io.EOF {
		return nil, fmt.Errorf("unexpected trailing data in file %q", inputPath)
	}
	return &state, nil
}

func WriteJSON[X any](outputPath string, value X, perm os.FileMode) error {
	if outputPath == "" {
		return nil
	}
	var out io.Writer
	finish := func() error { return nil }
	if outputPath != "-" {
		f, err := ioutil.NewAtomicWriterCompressed(outputPath, perm)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}
		// Ensure we close the stream even if failures occur.
		defer f.Close()
		out = f
		// Closing the file causes it to be renamed to the final destination
		// so make sure we handle any errors it returns
		finish = f.Close
	} else {
		out = os.Stdout
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("failed to encode to JSON: %w", err)
	}
	_, err := out.Write([]byte{'\n'})
	if err != nil {
		return fmt.Errorf("failed to append new-line: %w", err)
	}
	if err := finish(); err != nil {
		return fmt.Errorf("failed to finish write: %w", err)
	}
	return nil
}

// LoadJSONFieldStrict loads a JSON file and extracts a specific field, returning an error if
// any unknown fields are present.
func LoadJSONFieldStrict[X any](inputPath string, field string) (*X, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", inputPath, err)
	}
	defer file.Close()

	var raw map[string]json.RawMessage
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	fieldData, ok := raw[field]
	if !ok {
		return nil, fmt.Errorf("field %q not found in %q", field, inputPath)
	}

	var x X
	dec := json.NewDecoder(bytes.NewReader(fieldData))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&x); err != nil {
		return nil, fmt.Errorf("failed to decode field %q: %w", field, err)
	}
	return &x, nil
}

// WriteJSONToTarget writes JSON to an OutputTarget.
func WriteJSONToTarget[X any](value X, target ioutil.OutputTarget) error {
	out, closer, abort, err := target()
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	defer abort()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	return closer.Close()
}
