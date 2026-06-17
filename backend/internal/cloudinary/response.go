package cloudinary

import "encoding/json"

type uploadResponse struct {
	SecureURL string `json:"secure_url"`
}

func extractSecureURL(body []byte) (string, error) {
	var resp uploadResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
