package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Agency struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Name        string    `json:"name"`
	Website     *string   `json:"website"`
	LogoURL     *string   `json:"logoUrl"`
	Description *string   `json:"description"`
	IsVerified  bool      `json:"isVerified"`
	MaxCreators int       `json:"maxCreators"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AgencyCreator struct {
	ID        string     `json:"id"`
	AgencyID  string     `json:"agencyId"`
	CreatorID string     `json:"creatorId"`
	Status    string     `json:"status"`
	InvitedAt time.Time  `json:"invitedAt"`
	JoinedAt  *time.Time `json:"joinedAt"`
	// Joined fields for display
	CreatorName     *string `json:"creatorName,omitempty"`
	CreatorUsername *string `json:"creatorUsername,omitempty"`
	CreatorImage    *string `json:"creatorImage,omitempty"`
	CreatorScore    *int    `json:"creatorScore,omitempty"`
	// Agency fields for invite display
	AgencyName        *string `json:"agencyName,omitempty"`
	AgencyDescription *string `json:"agencyDescription,omitempty"`
}

type APIUsageStat struct {
	Date     string `json:"date"`
	Count    int    `json:"count"`
	Endpoint string `json:"endpoint,omitempty"`
}

type AgencyStats struct {
	TotalCreators int     `json:"totalCreators"`
	TotalViews    int     `json:"totalViews"`
	AvgScore      float64 `json:"avgScore"`
	APICalls30d   int     `json:"apiCalls30d"`
}

func (s *Store) CreateAgency(ctx context.Context, agency *Agency) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO agencies (id, user_id, name, website, logo_url, description, is_verified, max_creators, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		agency.ID, agency.UserID, agency.Name, agency.Website, agency.LogoURL, agency.Description,
		agency.IsVerified, agency.MaxCreators, agency.CreatedAt, agency.UpdatedAt,
	)
	return err
}

func (s *Store) FindAgencyByUserID(ctx context.Context, userID string) (*Agency, error) {
	var a Agency
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, website, logo_url, description, is_verified, max_creators, created_at, updated_at
		 FROM agencies WHERE user_id = $1`, userID,
	).Scan(&a.ID, &a.UserID, &a.Name, &a.Website, &a.LogoURL, &a.Description,
		&a.IsVerified, &a.MaxCreators, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) FindAgencyByID(ctx context.Context, id string) (*Agency, error) {
	var a Agency
	err := s.pool.QueryRow(ctx,
		`SELECT id, user_id, name, website, logo_url, description, is_verified, max_creators, created_at, updated_at
		 FROM agencies WHERE id = $1`, id,
	).Scan(&a.ID, &a.UserID, &a.Name, &a.Website, &a.LogoURL, &a.Description,
		&a.IsVerified, &a.MaxCreators, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (s *Store) UpdateAgency(ctx context.Context, id, name string, website, description *string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE agencies SET name = $1, website = $2, description = $3, updated_at = NOW() WHERE id = $4`,
		name, website, description, id,
	)
	return err
}

func (s *Store) InviteCreator(ctx context.Context, id, agencyID, creatorID string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO agency_creators (id, agency_id, creator_id, status, invited_at)
		 VALUES ($1, $2, $3, 'pending', NOW())`,
		id, agencyID, creatorID,
	)
	return err
}

func (s *Store) RespondToAgencyInvite(ctx context.Context, inviteID string, accept bool) error {
	if accept {
		_, err := s.pool.Exec(ctx,
			`UPDATE agency_creators SET status = 'active', joined_at = NOW() WHERE id = $1 AND status = 'pending'`,
			inviteID,
		)
		return err
	}
	_, err := s.pool.Exec(ctx,
		`UPDATE agency_creators SET status = 'removed' WHERE id = $1 AND status = 'pending'`,
		inviteID,
	)
	return err
}

func (s *Store) RemoveAgencyCreator(ctx context.Context, agencyID, creatorID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE agency_creators SET status = 'removed' WHERE agency_id = $1 AND creator_id = $2 AND status = 'active'`,
		agencyID, creatorID,
	)
	return err
}

func (s *Store) ListAgencyCreators(ctx context.Context, agencyID string, limit, offset int) ([]AgencyCreator, int, error) {
	// Count
	var total int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM agency_creators WHERE agency_id = $1 AND status = 'active'`,
		agencyID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.pool.Query(ctx,
		`SELECT ac.id, ac.agency_id, ac.creator_id, ac.status, ac.invited_at, ac.joined_at,
		        u.name, u.username, u.image, u.creator_score
		 FROM agency_creators ac
		 JOIN users u ON u.id = ac.creator_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active'
		 ORDER BY ac.joined_at DESC NULLS LAST
		 LIMIT $2 OFFSET $3`,
		agencyID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []AgencyCreator
	for rows.Next() {
		var ac AgencyCreator
		if err := rows.Scan(&ac.ID, &ac.AgencyID, &ac.CreatorID, &ac.Status, &ac.InvitedAt, &ac.JoinedAt,
			&ac.CreatorName, &ac.CreatorUsername, &ac.CreatorImage, &ac.CreatorScore); err != nil {
			return nil, 0, err
		}
		results = append(results, ac)
	}
	if results == nil {
		results = []AgencyCreator{}
	}
	return results, total, nil
}

func (s *Store) ListCreatorInvites(ctx context.Context, creatorID string) ([]AgencyCreator, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT ac.id, ac.agency_id, ac.creator_id, ac.status, ac.invited_at, ac.joined_at,
		        a.name, a.description
		 FROM agency_creators ac
		 JOIN agencies a ON a.id = ac.agency_id
		 WHERE ac.creator_id = $1 AND ac.status = 'pending'
		 ORDER BY ac.invited_at DESC`,
		creatorID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AgencyCreator
	for rows.Next() {
		var ac AgencyCreator
		if err := rows.Scan(&ac.ID, &ac.AgencyID, &ac.CreatorID, &ac.Status, &ac.InvitedAt, &ac.JoinedAt,
			&ac.AgencyName, &ac.AgencyDescription); err != nil {
			return nil, err
		}
		results = append(results, ac)
	}
	if results == nil {
		results = []AgencyCreator{}
	}
	return results, nil
}

func (s *Store) FindAgencyInvite(ctx context.Context, inviteID string) (*AgencyCreator, error) {
	var ac AgencyCreator
	err := s.pool.QueryRow(ctx,
		`SELECT id, agency_id, creator_id, status, invited_at, joined_at
		 FROM agency_creators WHERE id = $1`, inviteID,
	).Scan(&ac.ID, &ac.AgencyID, &ac.CreatorID, &ac.Status, &ac.InvitedAt, &ac.JoinedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &ac, err
}

func (s *Store) CountActiveAgencyCreators(ctx context.Context, agencyID string) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM agency_creators WHERE agency_id = $1 AND status = 'active'`,
		agencyID,
	).Scan(&count)
	return count, err
}

func (s *Store) GetAgencyStats(ctx context.Context, agencyID string) (*AgencyStats, error) {
	stats := &AgencyStats{}

	// Total active creators
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM agency_creators WHERE agency_id = $1 AND status = 'active'`,
		agencyID,
	).Scan(&stats.TotalCreators)
	if err != nil {
		return nil, err
	}

	// Aggregate views across managed creators
	err = s.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(pv.ct), 0)
		 FROM agency_creators ac
		 JOIN (SELECT user_id, COUNT(*) AS ct FROM profile_views GROUP BY user_id) pv ON pv.user_id = ac.creator_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active'`,
		agencyID,
	).Scan(&stats.TotalViews)
	if err != nil {
		return nil, err
	}

	// Average score
	err = s.pool.QueryRow(ctx,
		`SELECT COALESCE(AVG(u.creator_score), 0)
		 FROM agency_creators ac
		 JOIN users u ON u.id = ac.creator_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active' AND u.creator_score IS NOT NULL`,
		agencyID,
	).Scan(&stats.AvgScore)
	if err != nil {
		return nil, err
	}

	// API calls in last 30 days (across all keys owned by the agency user)
	var agency Agency
	err = s.pool.QueryRow(ctx,
		`SELECT user_id FROM agencies WHERE id = $1`, agencyID,
	).Scan(&agency.UserID)
	if err != nil {
		return nil, err
	}

	err = s.pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM api_usage au
		 JOIN api_keys ak ON ak.id = au.api_key_id
		 WHERE ak.user_id = $1 AND au.created_at >= NOW() - INTERVAL '30 days'`,
		agency.UserID,
	).Scan(&stats.APICalls30d)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Store) RecordAPIUsage(ctx context.Context, apiKeyID, endpoint, method string, statusCode, responseTimeMs int) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO api_usage (api_key_id, endpoint, method, status_code, response_time_ms, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW())`,
		apiKeyID, endpoint, method, statusCode, responseTimeMs,
	)
	return err
}

func (s *Store) GetAPIUsageStats(ctx context.Context, apiKeyID string, days int) ([]APIUsageStat, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT TO_CHAR(created_at::date, 'YYYY-MM-DD') AS date, COUNT(*) AS count
		 FROM api_usage
		 WHERE api_key_id = $1 AND created_at >= NOW() - MAKE_INTERVAL(days => $2)
		 GROUP BY created_at::date
		 ORDER BY created_at::date`,
		apiKeyID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []APIUsageStat
	for rows.Next() {
		var s APIUsageStat
		if err := rows.Scan(&s.Date, &s.Count); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	if results == nil {
		results = []APIUsageStat{}
	}
	return results, nil
}

func (s *Store) GetAPIUsageSummary(ctx context.Context, userID string, days int) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	// Total calls
	var totalCalls int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM api_usage au
		 JOIN api_keys ak ON ak.id = au.api_key_id
		 WHERE ak.user_id = $1 AND au.created_at >= NOW() - MAKE_INTERVAL(days => $2)`,
		userID, days,
	).Scan(&totalCalls)
	if err != nil {
		return nil, err
	}
	result["totalCalls"] = totalCalls

	// By endpoint
	rows, err := s.pool.Query(ctx,
		`SELECT au.endpoint, COUNT(*) AS count
		 FROM api_usage au
		 JOIN api_keys ak ON ak.id = au.api_key_id
		 WHERE ak.user_id = $1 AND au.created_at >= NOW() - MAKE_INTERVAL(days => $2)
		 GROUP BY au.endpoint
		 ORDER BY count DESC`,
		userID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var byEndpoint []map[string]interface{}
	for rows.Next() {
		var endpoint string
		var count int
		if err := rows.Scan(&endpoint, &count); err != nil {
			return nil, err
		}
		byEndpoint = append(byEndpoint, map[string]interface{}{
			"endpoint": endpoint,
			"count":    count,
		})
	}
	if byEndpoint == nil {
		byEndpoint = []map[string]interface{}{}
	}
	result["byEndpoint"] = byEndpoint

	// By day
	rows2, err := s.pool.Query(ctx,
		fmt.Sprintf(`SELECT TO_CHAR(au.created_at::date, 'YYYY-MM-DD') AS date, COUNT(*) AS count
		 FROM api_usage au
		 JOIN api_keys ak ON ak.id = au.api_key_id
		 WHERE ak.user_id = $1 AND au.created_at >= NOW() - MAKE_INTERVAL(days => $2)
		 GROUP BY au.created_at::date
		 ORDER BY au.created_at::date`),
		userID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var byDay []map[string]interface{}
	for rows2.Next() {
		var date string
		var count int
		if err := rows2.Scan(&date, &count); err != nil {
			return nil, err
		}
		byDay = append(byDay, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	if byDay == nil {
		byDay = []map[string]interface{}{}
	}
	result["byDay"] = byDay

	return result, nil
}

func (s *Store) UpdateUserRole(ctx context.Context, userID string, role string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`,
		role, userID,
	)
	return err
}

func (s *Store) BulkSetVerified(ctx context.Context, creatorIDs []string, verified bool) error {
	for _, id := range creatorIDs {
		if err := s.AdminSetVerified(ctx, id, verified); err != nil {
			return err
		}
	}
	return nil
}

// GetAgencyAnalytics returns aggregate analytics across all managed creators
func (s *Store) GetAgencyAnalytics(ctx context.Context, agencyID string) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	// Total views
	var totalViews int
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(COUNT(*), 0)
		 FROM profile_views pv
		 JOIN agency_creators ac ON ac.creator_id = pv.user_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active'`,
		agencyID,
	).Scan(&totalViews)
	if err != nil {
		return nil, err
	}
	result["totalViews"] = totalViews

	// Total clicks
	var totalClicks int
	err = s.pool.QueryRow(ctx,
		`SELECT COALESCE(COUNT(*), 0)
		 FROM link_clicks lc
		 JOIN agency_creators ac ON ac.creator_id = lc.user_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active'`,
		agencyID,
	).Scan(&totalClicks)
	if err != nil {
		return nil, err
	}
	result["totalClicks"] = totalClicks

	// Views by day (last 30 days)
	rows, err := s.pool.Query(ctx,
		`SELECT TO_CHAR(pv.viewed_at::date, 'YYYY-MM-DD') AS date, COUNT(*) AS count
		 FROM profile_views pv
		 JOIN agency_creators ac ON ac.creator_id = pv.user_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active' AND pv.viewed_at >= NOW() - INTERVAL '30 days'
		 GROUP BY pv.viewed_at::date
		 ORDER BY pv.viewed_at::date`,
		agencyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var viewsByDay []map[string]interface{}
	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		viewsByDay = append(viewsByDay, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	if viewsByDay == nil {
		viewsByDay = []map[string]interface{}{}
	}
	result["viewsByDay"] = viewsByDay

	// Top creators by views
	rows2, err := s.pool.Query(ctx,
		`SELECT u.name, u.username, COUNT(pv.id) AS views
		 FROM agency_creators ac
		 JOIN users u ON u.id = ac.creator_id
		 LEFT JOIN profile_views pv ON pv.user_id = ac.creator_id
		 WHERE ac.agency_id = $1 AND ac.status = 'active'
		 GROUP BY u.id, u.name, u.username
		 ORDER BY views DESC
		 LIMIT 10`,
		agencyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	var topCreators []map[string]interface{}
	for rows2.Next() {
		var name, username *string
		var views int
		if err := rows2.Scan(&name, &username, &views); err != nil {
			return nil, err
		}
		topCreators = append(topCreators, map[string]interface{}{
			"name":     name,
			"username": username,
			"views":    views,
		})
	}
	if topCreators == nil {
		topCreators = []map[string]interface{}{}
	}
	result["topCreators"] = topCreators

	return result, nil
}
