package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Uploader stores an image and returns its publicly servable URL. Folder
// mirrors the PRD's /team1/articles/{article-id}/{banner|inline} layout.
type Uploader interface {
	Upload(ctx context.Context, file []byte, filename, folder string) (secureURL string, err error)
}

func NewUploader(cloudName, apiKey, apiSecret, publicAPIURL string, mock bool, uploadsDir string) Uploader {
	if mock {
		return &localUploader{baseURL: publicAPIURL, dir: uploadsDir}
	}
	return &cloudinaryUploader{cloudName: cloudName, apiKey: apiKey, apiSecret: apiSecret}
}

// localUploader stands in for Cloudinary when no API credentials are
// configured: it writes the file to disk and serves it from this same API
// process under /uploads, so every other module (banners, inline images)
// works identically whether or not Cloudinary credentials exist yet.
type localUploader struct {
	baseURL string
	dir     string
}

func (u *localUploader) Upload(_ context.Context, file []byte, filename, folder string) (string, error) {
	dir := filepath.Join(u.dir, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	name := fmt.Sprintf("%s-%s", uuid.NewString(), filepath.Base(filename))
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, file, 0o644); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/uploads/%s/%s", u.baseURL, folder, name), nil
}

type cloudinaryUploader struct {
	cloudName string
	apiKey    string
	apiSecret string
}

func (u *cloudinaryUploader) Upload(ctx context.Context, file []byte, filename, folder string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := u.sign(map[string]string{"folder": folder, "timestamp": timestamp})

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("api_key", u.apiKey)
	_ = writer.WriteField("timestamp", timestamp)
	_ = writer.WriteField("folder", folder)
	_ = writer.WriteField("signature", signature)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/auto/upload", u.cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("cloudinary upload failed (%d): %s", resp.StatusCode, string(respBody))
	}

	secureURL, err := extractSecureURL(respBody)
	if err != nil {
		return "", err
	}
	return secureURL, nil
}

// sign implements Cloudinary's signed-upload scheme: sort params
// alphabetically, join as key=value pairs with &, append the API secret,
// then SHA-1 hash. https://cloudinary.com/documentation/signatures
func (u *cloudinaryUploader) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	// folder, timestamp - already alphabetical for our fixed param set.
	joined := ""
	for i, k := range sortedKeys(keys) {
		if i > 0 {
			joined += "&"
		}
		joined += k + "=" + params[k]
	}
	h := sha1.New()
	h.Write([]byte(joined + u.apiSecret))
	return hex.EncodeToString(h.Sum(nil))
}

func sortedKeys(keys []string) []string {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j-1] > keys[j]; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	return keys
}
