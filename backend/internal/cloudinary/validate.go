package cloudinary

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

const (
	BannerMaxBytes = 5 * 1024 * 1024
	BannerMinSide  = 1360
)

var (
	ErrImageTooLarge     = errors.New("file exceeds the 5MB maximum")
	ErrImageTooSmall     = errors.New("image must be at least 1360x1360px")
	ErrUnsupportedFormat = errors.New("only JPG and PNG images are accepted")
)

// ValidateBanner enforces the PRD's banner upload constraints: JPG/PNG only,
// minimum 1360x1360px, maximum 5MB.
func ValidateBanner(file []byte) error {
	if len(file) > BannerMaxBytes {
		return ErrImageTooLarge
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(file))
	if err != nil || (format != "jpeg" && format != "png") {
		return ErrUnsupportedFormat
	}
	if cfg.Width < BannerMinSide || cfg.Height < BannerMinSide {
		return ErrImageTooSmall
	}
	return nil
}
