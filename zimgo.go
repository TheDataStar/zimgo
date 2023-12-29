package zimgo

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	zimMagicHeader = "ZIM\x00\x00\x00\x00\x00\x00\x00\x00"
)

// Header struct stores information about the .zim file header. includes metadata like version, creation time, etc.
type ZimHeader struct {
	Version      uint32
	CreationTime time.Time
}

// ZimFile struct to represent a .zim file and initialize it with necessary fields
// os.File instance to manage file operations, allows you to open, read, and seek within the file.
// Compression, if the .zim file supports different compression algorithms, will include a field to store the selected algorithm.
// Checksum, included a field for storing checksum information if the file uses checksums for integrity verification.
type ZimFile struct {
	FileHandle           *os.File
	File                 *os.File
	FilePath             string
	Header               ZimHeader
	Index                map[string]ZimIndexEntry
	CompressionAlgorithm string
	Checksum             string
	// Add other fields as needed.
}

// Index map or slice in ZimFile Struct. If the .zim file has an index structure it will be stored in the ZimFile. This could help in efficiently locating and accessing data within the file.
type ZimIndexEntry struct {
	Offset int64
	Size   int
}

// Function to open a .zim file and return a 'ZimFile' instance.
// Add your implementation to open and validate the .zim file
// Return a ZimFile instance or and error if the file cannot open
func Open(filePath string) (*ZimFile, error) {

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the magic header
	magicHeader := make([]byte, len(zimMagicHeader))
	_, err = file.Read(magicHeader)
	if err != nil {
		return nil, err
	}

	// Validate the magic header
	if string(magicHeader) != zimMagicHeader {
		return nil, errors.New("invalid zim file format")
	}

	//Create and return a ZimFile instance
	return &ZimFile{FilePath: filePath}, nil
}

// ZimFile struct methods for reading and extractin data from the .zim files.
func (zf *ZimFile) Read(offset int64, size int) ([]byte, error) {

	// Seek to the specified offset
	_, err := zf.File.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// Read the specified number of bytes
	data := make([]byte, size)
	n, err := zf.File.Read(data)
	if err != nil {
		return nil, err
	}

	//Check if the expected number of bytes were read
	if n != size {
		return nil, errors.New("unexpected end of the file")
	}

	return data, nil
}

func (zf *ZimFile) Extract(destPath string) error {
	// Extract the content of the .zim file to the specified destination.
	zipFile, err := zip.OpenReader(zf.FilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create the destination directory if it doesn't exist.
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Extract each file from the .zim archive.
	for _, file := range zipFile.File {
		//construct the destination path
		destFilePath := filepath.Join(destPath, file.Name)

		// Create the directory structure if it doesn't exist.
		err := os.MkdirAll(filepath.Dir(destFilePath), os.ModePerm)
		if err != nil {
			return err
		}

		// Open the source file in the .zim archive
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create the desitination file
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		//Copy the contents of the source file to the desitination file
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func (zf *ZimFile) Close() error {

	// Close the file.
	if zf.File != nil {
		err := zf.File.Close()
		zf.File = nil // Set nil to indicate that the file is closed.
		return err
	}
	return nil
}
