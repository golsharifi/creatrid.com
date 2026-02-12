package store

import (
	"context"
	"time"
)

type AnalyticsSummary struct {
	TotalViews    int            `json:"totalViews"`
	ViewsToday    int            `json:"viewsToday"`
	ViewsThisWeek int           `json:"viewsThisWeek"`
	TotalClicks   int            `json:"totalClicks"`
	ClicksByType  map[string]int `json:"clicksByType"`
	ViewsByDay    []DayStat      `json:"viewsByDay"`
	ClicksByDay   []DayStat      `json:"clicksByDay"`
	TopReferrers  []ReferrerStat `json:"topReferrers"`
	ViewsByHour   []HourStat     `json:"viewsByHour"`
}

type DayStat struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

type HourStat struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

func (s *Store) RecordProfileView(ctx context.Context, userID, viewerIP, referrer string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO profile_views (user_id, viewer_ip, referrer) VALUES ($1, $2, $3)`,
		userID, viewerIP, referrer,
	)
	return err
}

func (s *Store) RecordLinkClick(ctx context.Context, userID, linkType, linkValue, clickerIP string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO link_clicks (user_id, link_type, link_value, clicker_ip) VALUES ($1, $2, $3, $4)`,
		userID, linkType, linkValue, clickerIP,
	)
	return err
}

func (s *Store) GetAnalyticsSummary(ctx context.Context, userID string) (*AnalyticsSummary, error) {
	summary := &AnalyticsSummary{
		ClicksByType: make(map[string]int),
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -7)

	// Total views
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM profile_views WHERE user_id = $1`, userID,
	).Scan(&summary.TotalViews)
	if err != nil {
		return nil, err
	}

	// Views today
	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM profile_views WHERE user_id = $1 AND created_at >= $2`,
		userID, todayStart,
	).Scan(&summary.ViewsToday)
	if err != nil {
		return nil, err
	}

	// Views this week
	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM profile_views WHERE user_id = $1 AND created_at >= $2`,
		userID, weekStart,
	).Scan(&summary.ViewsThisWeek)
	if err != nil {
		return nil, err
	}

	// Total clicks
	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM link_clicks WHERE user_id = $1`, userID,
	).Scan(&summary.TotalClicks)
	if err != nil {
		return nil, err
	}

	// Clicks by type
	rows, err := s.pool.Query(ctx,
		`SELECT link_type, COUNT(*) FROM link_clicks WHERE user_id = $1 GROUP BY link_type`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var linkType string
		var count int
		if err := rows.Scan(&linkType, &count); err != nil {
			return nil, err
		}
		summary.ClicksByType[linkType] = count
	}

	// Views by day (last 30 days)
	dayRows, err := s.pool.Query(ctx,
		`SELECT DATE(created_at) AS day, COUNT(*)
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		 GROUP BY day ORDER BY day`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer dayRows.Close()

	for dayRows.Next() {
		var ds DayStat
		var d time.Time
		if err := dayRows.Scan(&d, &ds.Count); err != nil {
			return nil, err
		}
		ds.Date = d.Format("2006-01-02")
		summary.ViewsByDay = append(summary.ViewsByDay, ds)
	}

	if summary.ViewsByDay == nil {
		summary.ViewsByDay = []DayStat{}
	}

	// Clicks by day (last 30 days)
	clickDayRows, err := s.pool.Query(ctx,
		`SELECT DATE(created_at) AS day, COUNT(*)
		 FROM link_clicks
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		 GROUP BY day ORDER BY day`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer clickDayRows.Close()

	for clickDayRows.Next() {
		var ds DayStat
		var d time.Time
		if err := clickDayRows.Scan(&d, &ds.Count); err != nil {
			return nil, err
		}
		ds.Date = d.Format("2006-01-02")
		summary.ClicksByDay = append(summary.ClicksByDay, ds)
	}
	if summary.ClicksByDay == nil {
		summary.ClicksByDay = []DayStat{}
	}

	// Top referrers (last 30 days)
	refRows, err := s.pool.Query(ctx,
		`SELECT COALESCE(referrer, 'Direct') AS referrer, COUNT(*) AS count
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		 GROUP BY referrer ORDER BY count DESC LIMIT 5`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer refRows.Close()

	for refRows.Next() {
		var rs ReferrerStat
		if err := refRows.Scan(&rs.Referrer, &rs.Count); err != nil {
			return nil, err
		}
		summary.TopReferrers = append(summary.TopReferrers, rs)
	}
	if summary.TopReferrers == nil {
		summary.TopReferrers = []ReferrerStat{}
	}

	// Views by hour of day (last 30 days)
	hourRows, err := s.pool.Query(ctx,
		`SELECT EXTRACT(HOUR FROM created_at)::int AS hour, COUNT(*) AS count
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
		 GROUP BY hour ORDER BY hour`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer hourRows.Close()

	for hourRows.Next() {
		var hs HourStat
		if err := hourRows.Scan(&hs.Hour, &hs.Count); err != nil {
			return nil, err
		}
		summary.ViewsByHour = append(summary.ViewsByHour, hs)
	}
	if summary.ViewsByHour == nil {
		summary.ViewsByHour = []HourStat{}
	}

	return summary, nil
}
