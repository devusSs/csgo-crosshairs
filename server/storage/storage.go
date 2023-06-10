package storage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/devusSs/crosshairs/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	userPPBucketName     = "profiles"
	allowedFileExtension = "png"
)

type Service struct {
	client *minio.Client
}

func NewMinioConnection(cfg *config.Config) (*Service, error) {
	var endpoint string
	var useSSL bool

	endpoint = fmt.Sprintf("%s:%d", cfg.MinioHost, cfg.MinioPort)

	if cfg.MinioDomain != "" {
		if strings.Contains(cfg.MinioDomain, "https://") {
			useSSL = true
			endpoint = strings.Replace(cfg.MinioDomain, "https://", "", 1)
		} else if strings.Contains(cfg.MinioDomain, "http://") {
			useSSL = false
			endpoint = strings.Replace(cfg.MinioDomain, "http://", "", 1)
		} else {
			return nil, errors.New("missing http schema in minio domain")
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

func (s *Service) CheckMinioVersion() (string, error) {
	if !s.CheckMinioConnection() {
		return "", errors.New("minio client not online")
	}

	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe = filepath.Dir(exe)

	httpFilePath := filepath.Join(exe, "minio.txt")

	httpFile, err := os.Create(httpFilePath)
	if err != nil {
		return "", err
	}

	s.client.TraceOn(httpFile)

	_, err = s.client.GetBucketPolicy(context.Background(), userPPBucketName)
	if err != nil {
		return "", err
	}

	s.client.TraceOff()

	httpFile.Close()

	httpFile, err = os.Open(httpFilePath)
	if err != nil {
		return "", err
	}

	var version string

	scanner := bufio.NewScanner(httpFile)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "User-Agent: ") {
			version = strings.Replace(line, "User-Agent: ", "", 1)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	httpFile.Close()

	if err := os.Remove(filepath.Join(exe, "minio.txt")); err != nil {
		return "", err
	}

	if version == "" {
		return "", errors.New("missing minio version info in header")
	}

	return version, nil
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

func (s *Service) DeleteUserProfilePicture(userID string) error {
	if !s.CheckMinioConnection() {
		return errors.New("minio client not online")
	}

	for object := range s.client.ListObjects(context.Background(), userPPBucketName, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			return object.Err
		}

		if object.Key == fmt.Sprintf("%s.png", userID) {
			if err := s.client.RemoveObject(context.Background(), userPPBucketName, object.Key, minio.RemoveObjectOptions{}); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("no matching object found")
}

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
