package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"shin-monta-no-mori/server/pkg/util"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type StorageService interface {
	UploadFile(ctx *gin.Context, file multipart.File, filename string, fileType string) (string, error)
	DeleteFile(ctx *gin.Context, filePath string) error
}

type GCSStorageService struct {
	Config util.Config
}

// GCSアップロード
func (g *GCSStorageService) UploadFile(ctx *gin.Context, file multipart.File, filename string, fileType string) (string, error) {
	client, err := createClient(ctx, g.Config)
	if err != nil {
		return "", fmt.Errorf("cannot create client : %w", err)
	}

	bucket := client.Bucket(g.Config.BucketName)
	if filename == "" {
		filename = time.Now().Format("20060102150405")
	}
	gcsFileName := fmt.Sprintf("%s/%s/%s.png", fileType, g.Config.Environment, filename)

	obj := bucket.Object(gcsFileName)

	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("error writing file : %w", err)
	}
	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("error closing file : %w", err)
	}

	err = updateMetadata(ctx, obj)
	if err != nil {
		return "", fmt.Errorf("error update storage metadata : %w", err)
	}

	resImagePath := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.Config.BucketName, gcsFileName)
	return resImagePath, nil
}

// GCS上の画像を削除する
func (g *GCSStorageService) DeleteFile(ctx *gin.Context, deleteSrcPath string) error {
	client, err := createClient(ctx, g.Config)
	if err != nil {
		return fmt.Errorf("failed to create : %w", err)
	}

	bucket := client.Bucket(g.Config.BucketName)
	objectPath := strings.TrimPrefix(deleteSrcPath, fmt.Sprintf("https://storage.googleapis.com/%s/", g.Config.BucketName))
	obj := bucket.Object(objectPath)

	err = obj.Delete(ctx)
	if err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("failed to delete object : %w", err)
	}

	return nil
}

// GCSクライアントとの接続
func createClient(ctx *gin.Context, config util.Config) (*storage.Client, error) {
	credentialFilePath := config.JsonPath
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create : %w", err)
	}

	defer client.Close()

	return client, err
}

func updateMetadata(ctx *gin.Context, object *storage.ObjectHandle) error {
	// メタデータの更新
	attrsToUpdate := storage.ObjectAttrsToUpdate{
		CacheControl: "no-cache",
	}
	if _, err := object.Update(ctx, attrsToUpdate); err != nil {
		return err
	}
	return nil
}
