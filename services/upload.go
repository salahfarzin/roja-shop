package services

import (
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
)

type Upload interface {
	SaveFile(ctx *fiber.Ctx, formField string, cfg *configs.Configs) (*types.File, string, error)
}

type upload struct {
}

func NewUpload() Upload {
	return &upload{}
}

func (u *upload) SaveFile(ctx *fiber.Ctx, formField string, cfg *configs.Configs) (*types.File, string, error) {
	fileHeader, err := ctx.FormFile(formField)
	if err != nil || fileHeader == nil {
		return nil, "", err
	}

	os.MkdirAll(cfg.UploadPath, os.ModePerm)
	ext := filepath.Ext(fileHeader.Filename)

	fileID := uuid.New().String()
	filename := fileID + ext
	destPath := filepath.Join(cfg.UploadPath, filename)
	if err := ctx.SaveFile(fileHeader, destPath); err != nil {
		return nil, "", err
	}

	publicPath := "/uploads/" + filename
	file := &types.File{
		ID:        fileID,
		Name:      filename,
		Path:      publicPath,
		Type:      fileHeader.Header.Get("Content-Type"),
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	return file, publicPath, nil
}
