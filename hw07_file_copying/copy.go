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

	if fileInfo.Mode().IsDir() {
		return errors.New("cannot copy directories")
	}

	if fileInfo.Mode().Perm()&0o400 == 0 {
		return errors.New("source file not readable")
	}

	if fileInfo.Mode()&os.ModeDevice != 0 { // check for os device /dev/urandom, /dev/null, etc.
		return fmt.Errorf("cannot copy device files: %s", from)
	}

	// check files identical
	if from == to {
		return fmt.Errorf("destination file is same as source file: %s", to)
	}

	// validate output dest file against the from file
	if fi, err := os.Stat(to); err == nil {
		if fi.Mode()&os.ModeDevice != 0 {
			return fmt.Errorf("output path is a device file: %s", to)
		}
		if fi.Mode()&os.ModeNamedPipe != 0 {
			return fmt.Errorf("output path is a named pipe: %s", to)
		}
		if fi.IsDir() {
			return fmt.Errorf("output path is a directory: %s", to)
		}

		if os.SameFile(fileInfo, fi) { // where: fileinfo is the source, fi - dest files
			return fmt.Errorf("destination file is symlinked to source file: %s", to)
		}
	}

	// Resolve symlinks and get absolute paths - redundant as SameFile introduced covers this earlier!
	// abs1, err := filepath.Abs(from)
	// if err != nil {
	// 	return fmt.Errorf("cannot check symlink of src file: %w", err)
	// }
	// abs2, err := filepath.Abs(to)
	// if err != nil {
	// 	return fmt.Errorf("cannot check symlink of dest file: %w", err)
	// }

	// if abs1 == abs2 {
	// 	return fmt.Errorf("destination file is same as source file: %s", to)
	// }

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
