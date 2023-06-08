package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/devusSs/crosshairs/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	userPPBucketName     = "profiles"
	allowedFileExtension = ".png"
)

type Service struct {
	client *minio.Client
}

func NewMinioConnection(cfg *config.Config) (*Service, error) {
	var endpoint string
	var useSSL bool

	endpoint = fmt.Sprintf("%s:%d", cfg.MinioHost, cfg.MinioPort)

	if cfg.MinioDomain != "" {
		endpoint = cfg.MinioDomain

		if strings.Contains(cfg.MinioDomain, "https") {
			useSSL = true
		}
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioUser, cfg.MinioPassword, ""),
		Secure: useSSL,
	})

	return &Service{client}, err
}

func (s *Service) CheckMinioConnection() bool {
	return s.client.IsOnline()
}

func (s *Service) CreateUserProfilePicturesBucket() error {
	if !s.CheckMinioConnection() {
		return errors.New("minio client not online")
	}

	userPPBucketExists, err := s.client.BucketExists(context.Background(), userPPBucketName)
	if err != nil {
		return err
	}

	if !userPPBucketExists {
		if err := s.client.MakeBucket(context.Background(), userPPBucketName, minio.MakeBucketOptions{
			Region:        "eu-west-1",
			ObjectLocking: false,
		}); err != nil {
			return err
		}
	}

	readOnlyPolicy := `{"Version":"2012-10-17",
	"Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation"],
	"Resource":["arn:aws:s3:::` + userPPBucketName + `"]},
	{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + userPPBucketName + `/*"]}]}`

	return s.client.SetBucketPolicy(context.Background(), userPPBucketName, readOnlyPolicy)
}

func (s *Service) UpdateUserProfilePicture(fileName, filePath string) error {
	if !s.CheckMinioConnection() {
		return errors.New("minio client not online")
	}

	_, err := s.client.FPutObject(context.Background(), userPPBucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType:     "image/png",
		ContentLanguage: "en-US",
	})
	return err
}

func (s *Service) GetUserProfilePictureLink(userID string) (string, error) {
	if !s.CheckMinioConnection() {
		return "", errors.New("minio client not online")
	}

	for object := range s.client.ListObjects(context.Background(), userPPBucketName, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			return "", object.Err
		}

		if object.Key == fmt.Sprintf("%s.png", userID) {
			return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL().String(), userPPBucketName, object.Key), nil
		}
	}

	return "", errors.New("no matching object found")
}

// TODO: also replace fileName in route with .Base
func CheckFileValid(file *multipart.FileHeader) error {
	fileName := filepath.Base(file.Filename)

	fileNameSplit := strings.Split(fileName, ".")
	ext := fileNameSplit[len(fileNameSplit)-1]

	if ext != allowedFileExtension {
		return fmt.Errorf("invalid extension, only %s allowed", allowedFileExtension)
	}

	readFile, err := file.Open()
	if err != nil {
		return err
	}
	defer readFile.Close()

	bytes, err := io.ReadAll(readFile)
	if err != nil {
		return err
	}

	mimeTypeIncipit := mimeFromIncipit(bytes)

	if mimeTypeIncipit == "" {
		return errors.New("could not determine mime type")
	}

	if ext != strings.Split(mimeTypeIncipit, "/")[1] {
		return fmt.Errorf("filename and mime type mismatch: %s <-> %s", ext, strings.Split(mimeTypeIncipit, "/")[1])
	}

	mimeType := http.DetectContentType(bytes)

	if ext != strings.Split(mimeType, "/")[1] {
		return fmt.Errorf("filename and mime type mismatch: %s <-> %s", ext, strings.Split(mimeType, "/")[1])
	}

	if mimeTypeIncipit != mimeType {
		return fmt.Errorf("mime type incipit and mime type http mismatch: %s <-> %s", mimeTypeIncipit, mimeType)
	}

	return nil
}

func mimeFromIncipit(incipit []byte) string {
	var magicTable = map[string]string{
		"\xff\xd8\xff":      "image/jpeg",
		"\x89PNG\r\n\x1a\n": "image/png",
		"GIF87a":            "image/gif",
		"GIF89a":            "image/gif",
	}

	incipitStr := string([]byte(incipit))
	for magic, mime := range magicTable {
		if strings.HasPrefix(incipitStr, magic) {
			return mime
		}
	}

	return ""
}