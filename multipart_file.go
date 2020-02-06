package idemploader

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
)

// MultipartFile _
type MultipartFile struct {
	file     *multipart.FileHeader
	checksum string
	reader   io.Reader
	size     int
}

// NewMultipartFile _
func NewMultipartFile(file *multipart.FileHeader) (*MultipartFile, error) {
	reader, err := file.Open()

	if err != nil {
		return nil, errors.New("Opening file error: " + err.Error())
	}

	checksum, err := calculateChecksum(reader)

	if err != nil {
		return nil, err
	}

	return &MultipartFile{
		file:     file,
		reader:   reader,
		checksum: checksum,
	}, nil
}

// Size _
func (mf *MultipartFile) Size() int {
	return mf.size
}

// Checksum _
func (mf *MultipartFile) Checksum() string {
	return mf.checksum
}

// ContentType _
func (mf *MultipartFile) ContentType() string {
	return mf.file.Header.Get("Content-Type")
}

// Read _
func (mf *MultipartFile) Read(p []byte) (n int, err error) {
	if mf.reader == nil {
		mf.reader, err = mf.file.Open()

		if err != nil {
			return 0, err
		}
	}

	return mf.reader.Read(p)

}

func calculateChecksum(reader io.Reader) (string, error) {
	hash := sha256.New()

	_, err := io.Copy(hash, reader)

	if err != nil {
		return "", errors.New("IO file copying error: " + err.Error())
	}

	sum := hash.Sum(nil)
	return hex.EncodeToString(sum), nil
}
