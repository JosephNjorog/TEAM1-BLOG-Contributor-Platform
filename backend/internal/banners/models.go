package banners

import (
	"time"

	"github.com/google/uuid"
)

type Banner struct {
	ID            uuid.UUID  `json:"id"`
	ArticleID     uuid.UUID  `json:"articleId"`
	DesignerID    uuid.UUID  `json:"designerId"`
	DesignerName  string     `json:"designerName"`
	CloudinaryURL string     `json:"cloudinaryUrl"`
	UploadedAt    time.Time  `json:"uploadedAt"`
	MarkedReadyAt *time.Time `json:"markedReadyAt"`
}
