package backup

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dest.
func copyFile(src, dest string, mode os.FileMode) error {
	// Open the source file
	srcFile, err := os.OpenFile(src, os.O_RDONLY, mode)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		err = srcFile.Close()
		if err != nil {
			panic(err)
		}
	}(srcFile)

	// Create or open the destination file
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		// Close the source file if opening destination file fails
		err = srcFile.Close()
		if err != nil {
			return err
		}
		return err
	}
	defer func(destFile *os.File) {
		err = destFile.Close()
		if err != nil {
			panic(err)
		}
	}(destFile)

	// Copy the contents of the source file to the destination file
	if _, err = io.Copy(destFile, srcFile); err != nil {
		return err
	}

	// Set permissions on destination files
	return os.Chmod(dest, 0644)
}

// CopyDir copies a directory from src to dest.
func CopyDir(src, dest string) error {
	// Get the source directory information
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err = os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	// Get the contents of the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy the contents of the source directory to the destination directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// If it is a subdirectory, call CopyDir recursively
			if err = CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// If it is a file, call copyFile to copy the file content
			if err = copyFile(srcPath, destPath, entry.Type()); err != nil {
				return err
			}
		}
	}

	return nil
}
