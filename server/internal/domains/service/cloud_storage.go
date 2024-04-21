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

// GCSアップロード
func UploadToGCS(ctx *gin.Context, config util.Config, file *multipart.FileHeader, filename string, fileType string) (string, error) {
	// setting for upload request
	client, err := createGCSClient(ctx, config)
	if err != nil {
		return "", fmt.Errorf("cannot create client : %w", err)
	}

	bucket := client.Bucket(config.BucketName)
	if filename == "" {
		filename = time.Now().Format("20060102150405")
	}
	gcsFileName := fmt.Sprintf("%s/%s/%s.png", fileType, config.Environment, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file : %w", err)
	}
	defer src.Close()
	obj := bucket.Object(gcsFileName)

	// write file to gcs
	wc := obj.NewWriter(ctx)
	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("error writing file : %w", err)
	}
	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("error closing file : %w", err)
	}

	err = updateMetadata(ctx, client, bucket, obj)
	if err != nil {
		return "", fmt.Errorf("error update storage metadata : %w", err)
	}

	resImagePath := fmt.Sprintf("https://storage.googleapis.com/%s/%s", config.BucketName, gcsFileName)
	return resImagePath, nil
}

// GCS上の画像を削除する
func DeleteFileFromGCS(ctx *gin.Context, deleteSrcPath string, config util.Config) error {
	client, err := createGCSClient(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create : %w", err)
	}

	bucket := client.Bucket(config.BucketName)
	objectPath := strings.TrimPrefix(deleteSrcPath, fmt.Sprintf("https://storage.googleapis.com/%s/", config.BucketName))
	obj := bucket.Object(objectPath)

	err = obj.Delete(ctx)
	if err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("failed to delete object : %w", err)
	}

	return nil
}

// GCSクライアントとの接続
func createGCSClient(ctx *gin.Context, config util.Config) (*storage.Client, error) {
	credentialFilePath := config.JsonPath
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create : %w", err)
	}

	defer client.Close()

	return client, err
}

func updateMetadata(ctx *gin.Context, client *storage.Client, bucket *storage.BucketHandle, object *storage.ObjectHandle) error {
	// メタデータの更新
	attrsToUpdate := storage.ObjectAttrsToUpdate{
		CacheControl: "no-cache",
	}
	if _, err := object.Update(ctx, attrsToUpdate); err != nil {
		return err
	}
	return nil
}
