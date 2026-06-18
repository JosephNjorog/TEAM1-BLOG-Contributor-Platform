package cloudinary

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func encodePNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 232, G: 65, B: 66, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test PNG: %v", err)
	}
	return buf.Bytes()
}

func encodeJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to encode test JPEG: %v", err)
	}
	return buf.Bytes()
}

func TestValidateBanner_Valid(t *testing.T) {
	if err := ValidateBanner(encodePNG(t, BannerMinSide, BannerMinSide)); err != nil {
		t.Errorf("expected a %dx%d PNG to be valid, got error: %v", BannerMinSide, BannerMinSide, err)
	}
	if err := ValidateBanner(encodeJPEG(t, 1400, 1400)); err != nil {
		t.Errorf("expected a 1400x1400 JPEG to be valid, got error: %v", err)
	}
}

func TestValidateBanner_TooSmall(t *testing.T) {
	err := ValidateBanner(encodePNG(t, 400, 400))
	if !errors.Is(err, ErrImageTooSmall) {
		t.Errorf("expected ErrImageTooSmall for a 400x400 image, got %v", err)
	}
}

func TestValidateBanner_TooLarge(t *testing.T) {
	oversized := make([]byte, BannerMaxBytes+1)
	err := ValidateBanner(oversized)
	if !errors.Is(err, ErrImageTooLarge) {
		t.Errorf("expected ErrImageTooLarge for a file over %d bytes, got %v", BannerMaxBytes, err)
	}
}

func TestValidateBanner_UnsupportedFormat(t *testing.T) {
	err := ValidateBanner([]byte("this is not an image at all"))
	if !errors.Is(err, ErrUnsupportedFormat) {
		t.Errorf("expected ErrUnsupportedFormat for non-image bytes, got %v", err)
	}
}
