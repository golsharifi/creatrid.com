package store

import (
	"context"
	"time"
)

type ContentView struct {
	ID        int64     `json:"id"`
	ContentID string    `json:"contentId"`
	ViewerIP  *string   `json:"viewerIp"`
	Referrer  *string   `json:"referrer"`
	CreatedAt time.Time `json:"createdAt"`
}

type ContentDownload struct {
	ID               int64     `json:"id"`
	ContentID        string    `json:"contentId"`
	DownloaderUserID *string   `json:"downloaderUserId"`
	CreatedAt        time.Time `json:"createdAt"`
}

type ContentAnalyticsSummary struct {
	ContentID    string `json:"contentId"`
	Title        string `json:"title"`
	TotalViews   int    `json:"totalViews"`
	TotalDownloads int  `json:"totalDownloads"`
	RevenueCents int    `json:"revenueCents"`
}

func (s *Store) RecordContentView(ctx context.Context, contentID, viewerIP, referrer string) error {
	var ipPtr, refPtr *string
	if viewerIP != "" {
		ipPtr = &viewerIP
	}
	if referrer != "" {
		refPtr = &referrer
	}
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_views (content_id, viewer_ip, referrer) VALUES ($1, $2, $3)`,
		contentID, ipPtr, refPtr,
	)
	return err
}

func (s *Store) RecordContentDownload(ctx context.Context, contentID string, downloaderUserID *string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO content_downloads (content_id, downloader_user_id) VALUES ($1, $2)`,
		contentID, downloaderUserID,
	)
	return err
}

func (s *Store) GetContentAnalytics(ctx context.Context, contentID string) (*ContentAnalyticsSummary, error) {
	summary := &ContentAnalyticsSummary{ContentID: contentID}

	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(ci.title, '') FROM content_items ci WHERE ci.id = $1`, contentID,
	).Scan(&summary.Title)
	if err != nil {
		return nil, err
	}

	_ = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_views WHERE content_id = $1`, contentID,
	).Scan(&summary.TotalViews)

	_ = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM content_downloads WHERE content_id = $1`, contentID,
	).Scan(&summary.TotalDownloads)

	_ = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(amount_cents), 0) FROM license_purchases WHERE content_id = $1 AND status = 'completed'`, contentID,
	).Scan(&summary.RevenueCents)

	return summary, nil
}

func (s *Store) GetCreatorContentAnalytics(ctx context.Context, userID string) ([]*ContentAnalyticsSummary, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT ci.id, ci.title,
		        COALESCE((SELECT COUNT(*) FROM content_views cv WHERE cv.content_id = ci.id), 0) AS views,
		        COALESCE((SELECT COUNT(*) FROM content_downloads cd WHERE cd.content_id = ci.id), 0) AS downloads,
		        COALESCE((SELECT SUM(lp.amount_cents) FROM license_purchases lp WHERE lp.content_id = ci.id AND lp.status = 'completed'), 0) AS revenue
		 FROM content_items ci
		 WHERE ci.user_id = $1
		 ORDER BY views DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ContentAnalyticsSummary
	for rows.Next() {
		var s ContentAnalyticsSummary
		if err := rows.Scan(&s.ContentID, &s.Title, &s.TotalViews, &s.TotalDownloads, &s.RevenueCents); err != nil {
			return nil, err
		}
		items = append(items, &s)
	}
	if items == nil {
		items = []*ContentAnalyticsSummary{}
	}
	return items, nil
}
