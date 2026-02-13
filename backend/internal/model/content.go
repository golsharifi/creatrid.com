package model

import "time"

type ContentItem struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	ContentType  string    `json:"contentType"`
	MimeType     string    `json:"mimeType"`
	FileSize     int64     `json:"fileSize"`
	FileURL      string    `json:"-"`
	ThumbnailURL *string   `json:"thumbnailUrl"`
	HashSHA256   string    `json:"hashSha256"`
	IsPublic     bool      `json:"isPublic"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PublicContentItem struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	ContentType  string    `json:"contentType"`
	MimeType     string    `json:"mimeType"`
	FileSize     int64     `json:"fileSize"`
	ThumbnailURL *string   `json:"thumbnailUrl"`
	HashSHA256   string    `json:"hashSha256"`
	Tags         []string  `json:"tags"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (c *ContentItem) ToPublic() *PublicContentItem {
	return &PublicContentItem{
		ID:           c.ID,
		UserID:       c.UserID,
		Title:        c.Title,
		Description:  c.Description,
		ContentType:  c.ContentType,
		MimeType:     c.MimeType,
		FileSize:     c.FileSize,
		ThumbnailURL: c.ThumbnailURL,
		HashSHA256:   c.HashSHA256,
		Tags:         c.Tags,
		CreatedAt:    c.CreatedAt,
	}
}

type LicenseOffering struct {
	ID          string    `json:"id"`
	ContentID   string    `json:"contentId"`
	LicenseType string    `json:"licenseType"`
	PriceCents  int       `json:"priceCents"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"isActive"`
	TermsText   *string   `json:"termsText"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LicensePurchase struct {
	ID                string    `json:"id"`
	OfferingID        string    `json:"offeringId"`
	ContentID         string    `json:"contentId"`
	BuyerUserID       *string   `json:"buyerUserId"`
	BuyerEmail        string    `json:"buyerEmail"`
	BuyerCompany      *string   `json:"buyerCompany"`
	StripeSessionID   *string   `json:"stripeSessionId"`
	AmountCents       int       `json:"amountCents"`
	PlatformFeeCents  int       `json:"platformFeeCents"`
	CreatorPayoutCents int      `json:"creatorPayoutCents"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
}

type TakedownRequest struct {
	ID              string     `json:"id"`
	ReporterEmail   string     `json:"reporterEmail"`
	ReporterName    string     `json:"reporterName"`
	ContentID       string     `json:"contentId"`
	Reason          string     `json:"reason"`
	EvidenceURL     *string    `json:"evidenceUrl"`
	Status          string     `json:"status"`
	ResolvedBy      *string    `json:"resolvedBy"`
	ResolutionNotes *string    `json:"resolutionNotes"`
	CreatedAt       time.Time  `json:"createdAt"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
}
