package store

import (
	"context"
	"time"
)

type AnalyticsSummary struct {
	TotalViews       int                `json:"totalViews"`
	ViewsToday       int                `json:"viewsToday"`
	ViewsThisWeek    int                `json:"viewsThisWeek"`
	TotalClicks      int                `json:"totalClicks"`
	ClicksByType     map[string]int     `json:"clicksByType"`
	ViewsByDay       []DayStat          `json:"viewsByDay"`
	ClicksByDay      []DayStat          `json:"clicksByDay"`
	TopReferrers     []ReferrerStat     `json:"topReferrers"`
	ViewsByHour      []HourStat         `json:"viewsByHour"`
	DeviceBreakdown  []DeviceBreakdown  `json:"deviceBreakdown,omitempty"`
	BrowserBreakdown []BrowserBreakdown `json:"browserBreakdown,omitempty"`
	GeoBreakdown     []GeoBreakdown     `json:"geoBreakdown,omitempty"`
	BounceRate       *float64           `json:"bounceRate,omitempty"`
}

type DeviceBreakdown struct {
	DeviceType string `json:"deviceType"`
	Count      int    `json:"count"`
}

type BrowserBreakdown struct {
	Browser string `json:"browser"`
	Count   int    `json:"count"`
}

type GeoBreakdown struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
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

func (s *Store) RecordProfileView(ctx context.Context, userID, viewerIP, referrer, browser, os, deviceType, country, city string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO profile_views (user_id, viewer_ip, referrer, browser, os, device_type, country, city)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, viewerIP, referrer, nilIfEmpty(browser), nilIfEmpty(os), nilIfEmpty(deviceType), nilIfEmpty(country), nilIfEmpty(city),
	)
	return err
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
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

	// Device breakdown (last 30 days)
	devRows, err := s.pool.Query(ctx,
		`SELECT COALESCE(device_type, 'Unknown') AS device_type, COUNT(*) AS count
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days' AND device_type IS NOT NULL
		 GROUP BY device_type ORDER BY count DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer devRows.Close()

	for devRows.Next() {
		var d DeviceBreakdown
		if err := devRows.Scan(&d.DeviceType, &d.Count); err != nil {
			return nil, err
		}
		summary.DeviceBreakdown = append(summary.DeviceBreakdown, d)
	}

	// Browser breakdown (last 30 days, top 10)
	browRows, err := s.pool.Query(ctx,
		`SELECT COALESCE(browser, 'Unknown') AS browser, COUNT(*) AS count
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days' AND browser IS NOT NULL
		 GROUP BY browser ORDER BY count DESC LIMIT 10`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer browRows.Close()

	for browRows.Next() {
		var b BrowserBreakdown
		if err := browRows.Scan(&b.Browser, &b.Count); err != nil {
			return nil, err
		}
		summary.BrowserBreakdown = append(summary.BrowserBreakdown, b)
	}

	// Geo breakdown (last 30 days, top 10)
	geoRows, err := s.pool.Query(ctx,
		`SELECT COALESCE(country, 'Unknown') AS country, COUNT(*) AS count
		 FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days' AND country IS NOT NULL
		 GROUP BY country ORDER BY count DESC LIMIT 10`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer geoRows.Close()

	for geoRows.Next() {
		var g GeoBreakdown
		if err := geoRows.Scan(&g.Country, &g.Count); err != nil {
			return nil, err
		}
		summary.GeoBreakdown = append(summary.GeoBreakdown, g)
	}

	// Bounce rate (views with only 1 visit from same IP in last 30 days)
	var totalIPs, singleVisitIPs int
	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT viewer_ip) FROM profile_views
		 WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'`,
		userID,
	).Scan(&totalIPs)
	if err != nil {
		return nil, err
	}

	if totalIPs > 0 {
		err = s.pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM (
				SELECT viewer_ip FROM profile_views
				WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
				GROUP BY viewer_ip HAVING COUNT(*) = 1
			) AS single_visits`,
			userID,
		).Scan(&singleVisitIPs)
		if err != nil {
			return nil, err
		}
		rate := float64(singleVisitIPs) / float64(totalIPs) * 100
		summary.BounceRate = &rate
	}

	return summary, nil
}
