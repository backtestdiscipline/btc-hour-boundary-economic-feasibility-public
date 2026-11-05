package study

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"

	"research-platform/seal"
)

const maxArchiveBytes = 16 << 20
const maxCSVBytes = 10 << 20

// ReadArchive reads and closes the platform-bound source before interpreting it.
// The caller specifies the timestamp unit; the filename is never used to infer it.
func ReadArchive(in Inputs, ref seal.ArtifactRef, micros bool) ([]Candle, DecodeReport, error) {
	var d DecodeReport
	if ref.Size <= 0 || ref.Size > maxArchiveBytes {
		return nil, d, fmt.Errorf("source archive size outside bounds")
	}
	if !strings.HasSuffix(ref.Path, ".zip") {
		return nil, d, fmt.Errorf("source must be a ZIP archive")
	}
	r, err := in.OpenSource(ref)
	if err != nil {
		return nil, d, err
	}
	return decodeArchive(r, path.Base(ref.Path), micros)
}

func decodeArchive(r io.ReadCloser, archiveName string, micros bool) ([]Candle, DecodeReport, error) {
	var d DecodeReport
	b, readErr := io.ReadAll(io.LimitReader(r, maxArchiveBytes+1))
	closeErr := r.Close()
	if readErr != nil {
		return nil, d, fmt.Errorf("read source: %w", readErr)
	}
	if closeErr != nil {
		return nil, d, fmt.Errorf("close source: %w", closeErr)
	}
	if len(b) > maxArchiveBytes {
		return nil, d, fmt.Errorf("source archive exceeds limit")
	}
	z, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, d, fmt.Errorf("open archive: %w", err)
	}
	expected := strings.TrimSuffix(archiveName, ".zip") + ".csv"
	if !strings.HasSuffix(archiveName, ".zip") || len(z.File) != 1 || z.File[0].Name != expected || z.File[0].UncompressedSize64 > maxCSVBytes {
		return nil, d, fmt.Errorf("unexpected archive structure")
	}
	f, err := z.File[0].Open()
	if err != nil {
		return nil, d, fmt.Errorf("open csv: %w", err)
	}
	c, d, decodeErr := DecodeDetailed(io.LimitReader(f, maxCSVBytes+1), micros)
	closeErr = f.Close()
	if decodeErr != nil {
		return nil, d, decodeErr
	}
	if closeErr != nil {
		return nil, d, fmt.Errorf("close csv: %w", closeErr)
	}
	return c, d, nil
}
