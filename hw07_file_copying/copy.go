package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	//nolint:depguard // this import is allowed here for CLI U
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(from, to string, offset, limit int64) error {
	// Open source file
	src, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("cannot open source: %w", err) // Error wrapping
	}
	defer src.Close()

	// Get file size if limit is 0 (copy entire file)
	fileInfo, err := src.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat file: %w", err)
	}
	fileSize := fileInfo.Size()

	// Validate offset BEFORE any operations
	if offset > fileSize {
		return ErrOffsetExceedsFileSize // Fail fast
	}

	// Seek to offset (if specified)
	if offset > 0 {
		_, err = src.Seek(offset, io.SeekStart)
		if err != nil {
			return ErrOffsetExceedsFileSize
		}
	}

	// Adjust limit:
	remainingBytes := fileSize - offset
	if limit == 0 {
		limit = remainingBytes // Copy all remaining bytes
	} else if limit > remainingBytes {
		log.Printf("Warning: --limit (%d) exceeds remaining bytes (%d), adjusting", limit, remainingBytes)
		limit = remainingBytes
	}

	// Create progress bar
	bar := pb.Start64(limit)
	bar.Set(pb.Bytes, true) // Show in bytes
	bar.SetWidth(80)        // Set width of progress bar

	// Create proxy reader that updates progress bar
	proxyReader := bar.NewProxyReader(src)

	// Create destination file (overwrite if exists)
	dst, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer dst.Close()

	// Copy using the proxy reader
	_, err = io.CopyN(dst, proxyReader, limit)
	if err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	// finish bar
	bar.Finish()
	return nil
}
