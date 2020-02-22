package server

import (
	"errors"
	"io"

	"github.com/minio/minio-go/v6"
)

// MinioInteractor _
type MinioInteractor interface {
	PutObject(bucketName string, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error)
	StatObject(bucketName string, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
}

// FileURLBuilder _
type FileURLBuilder = func(checksum string) string

// Storage _
type Storage struct {
	S3Bucket   string
	S3Path     string
	URLBuilder FileURLBuilder

	s3Client MinioInteractor
}

// StorageFile _
type StorageFile struct {
	storage  *Storage
	s3Object *minio.Object
	URL      string `json:"url"`
	Size     int    `json:"size"`
	Checksum string `json:"checksum"`
	MimeType string `json:"mime_type"`
}

// Read _
func (sf *StorageFile) Read(p []byte) (n int, err error) {
	if sf.s3Object == nil {
		sf.s3Object, err = sf.getS3Object()

		if err != nil {
			return 0, err
		}
	}

	return sf.s3Object.Read(p)
}

// getS3Object _
func (sf *StorageFile) getS3Object() (*minio.Object, error) {
	if sf.storage == nil {
		panic(errors.New("StorageFile is missing .storage"))
	}

	if sf.storage.s3Client == nil {
		panic(errors.New("StorageFile.storage is missing .s3Client"))
	}

	object, err := sf.storage.s3Client.GetObject(
		sf.storage.S3Bucket,
		sf.storage.s3FilePath(sf.Checksum),
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	return object, nil
}

// ErrFileNotFoundCode _
const ErrFileNotFoundCode = "FILE_NOT_FOUND"

// ErrFileNotFound _
var ErrFileNotFound = errors.New(ErrFileNotFoundCode)

// StorageConfig _
type StorageConfig struct {
	S3Host         string
	S3AccessKey    string
	S3AccessSecret string
	S3Bucket       string
	S3Path         string
	URLBuilder     FileURLBuilder
}

// NewStorage _
func NewStorage(cfg StorageConfig) (*Storage, error) {
	s3Client, err := minio.New(
		cfg.S3Host,
		cfg.S3AccessKey,
		cfg.S3AccessSecret,
		true,
	)

	if err != nil {
		return nil, err
	}

	return &Storage{
		S3Bucket:   cfg.S3Bucket,
		S3Path:     cfg.S3Path,
		URLBuilder: cfg.URLBuilder,
		s3Client:   s3Client,
	}, nil
}

// Filer _
type Filer interface {
	Checksum() string
	ContentType() string
	Size() int
	Open() io.Reader
}

// UploadFileIfNotExists _
func (u *Storage) UploadFileIfNotExists(file Filer) error {
	isFileExists, err := u.isFileExists(file.Checksum())

	if err != nil {
		return err
	}

	if !isFileExists {
		return u.UploadFile(file, file.Open())
	}

	return nil
}

// UploadFile _
func (u *Storage) UploadFile(file Filer, reader io.Reader) error {
	_, err := u.s3Client.PutObject(
		u.S3Bucket,
		u.s3FilePath(file.Checksum()),
		reader,
		int64(file.Size()),
		minio.PutObjectOptions{
			ContentType: file.ContentType(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// GetFile _
func (u *Storage) GetFile(file Filer) (*StorageFile, error) {
	return u.GetFileByChecksum(file.Checksum())
}

// GetFileByChecksum _
func (u *Storage) GetFileByChecksum(checksum string) (*StorageFile, error) {
	fileObj, err := u.s3Client.StatObject(
		u.S3Bucket,
		u.s3FilePath(checksum),
		minio.StatObjectOptions{},
	)

	if err != nil {
		switch e := err.(type) {
		case minio.ErrorResponse:
			if e.Code == "NoSuchKey" {
				return nil, ErrFileNotFound
			}
		}

		return nil, err
	}

	return &StorageFile{
		storage:  u,
		Checksum: checksum,
		URL:      u.URLBuilder(checksum),
		Size:     int(fileObj.Size),
		MimeType: fileObj.ContentType,
	}, nil
}

func (u *Storage) isFileExists(checksum string) (bool, error) {
	uploadedFile, err := u.GetFileByChecksum(checksum)

	if uploadedFile != nil {
		return true, nil
	}

	if err == ErrFileNotFound {
		return false, nil
	}

	return false, err
}

func (u *Storage) s3FilePath(checksum string) string {
	return u.S3Path + checksum
}
