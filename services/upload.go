package services

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/salahfarzin/roja-shop/configs"
	"github.com/salahfarzin/roja-shop/types"
)

var (
	UploadService Upload = &upload{}
	ErrNotImage          = errors.New("only image files are allowed")
)

type Upload interface {
	SaveFile(ctx *fiber.Ctx, fileHeader *multipart.FileHeader, cfg *configs.Configs) (*types.File, string, error)
}

type upload struct {
}

func NewUpload() Upload {
	return &upload{}
}

func (u *upload) SaveFile(ctx *fiber.Ctx, fileHeader *multipart.FileHeader, cfg *configs.Configs) (*types.File, string, error) {
	// Only allow image uploads
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" || contentType[:6] != "image/" {
		return nil, "", ErrNotImage
	}

	os.MkdirAll(cfg.UploadPath, os.ModePerm)
	// Prevent path traversal by ignoring the original filename except for extension
	ext := filepath.Ext(filepath.Base(fileHeader.Filename))
	fileID := uuid.New().String()
	filename := fileID + ext
	destPath := filepath.Join(cfg.UploadPath, filename)

	if err := ctx.SaveFile(fileHeader, destPath); err != nil {
		return nil, "", err
	}
	// Set file permissions to 0644 (readable, not executable)
	if err := os.Chmod(destPath, 0644); err != nil {
		return nil, "", err
	}

	publicPath := "/uploads/" + filename
	file := &types.File{
		ID:        fileID,
		Name:      filename,
		Path:      publicPath,
		Type:      contentType,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	return file, publicPath, nil
}
