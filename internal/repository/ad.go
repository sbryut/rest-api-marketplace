package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"rest-api-marketplace/internal/entity"
)

// AdsRepo provides DB operations for ads
type AdsRepo struct {
	db *sql.DB
}

// NewAdsRepo creates a new AdsRepo instance
func NewAdsRepo(db *sql.DB) *AdsRepo {
	return &AdsRepo{db: db}
}

// Create inserts a new ad and returns its ID
func (r AdsRepo) Create(ctx context.Context, ad entity.Ad) (int64, error) {
	const op = "repository.AdsRepo.Create"

	query := `INSERT INTO ads (user_id, title, description, image_url, price) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, ad.UserID, ad.Title, ad.Description, ad.ImageURL, ad.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// Update modifies an existing ad
func (r AdsRepo) Update(ctx context.Context, id int64, ad entity.Ad) error {
	const op = "repository.AdsRepo.Update"

	query := `UPDATE ads SET title = $1, description = $2, image_url = $3, price = $4 WHERE id = $5`

	res, err := r.db.ExecContext(ctx, query, ad.Title, ad.Description, ad.ImageURL, ad.Price, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: check rows affected: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrAdNotFound)
	}

	return nil
}

// GetByID retrieves an ad by its ID
func (r AdsRepo) GetByID(ctx context.Context, id int64) (*entity.Ad, error) {
	const op = "repository.AdsRepo.GetById"

	query := `SELECT id, user_id, title, description, image_url, price, created_at FROM ads WHERE id = $1`

	var ad entity.Ad

	err := r.db.QueryRowContext(ctx, query, id).Scan(&ad.ID, &ad.UserID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrAdNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ad, nil
}

// GetByIDWithAuthor retrieves an ad with author info
func (r AdsRepo) GetByIDWithAuthor(ctx context.Context, id int64) (*entity.AdWithAuthor, error) {
	const op = "repository.AdsRepo.GetByIdWithAuthor"

	query := `SELECT a.id, a.user_id, a.title, a.description, a.image_url, a.price, a.created_at, u.login
			  FROM ads a
			  JOIN users u ON a.user_id = u.id
			  WHERE a.id = $1`

	var ad entity.AdWithAuthor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ad.ID,
		&ad.UserID,
		&ad.Title,
		&ad.Description,
		&ad.ImageURL,
		&ad.Price,
		&ad.CreatedAt,
		&ad.AuthorLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, entity.ErrAdNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &ad, nil
}

// GetAll returns a list of ads with optional filters, sorting, and pagination
func (r AdsRepo) GetAll(ctx context.Context, params entity.GetAdsQuery) ([]entity.AdWithAuthor, error) {
	const op = "repository.AdsRepo.GetAll"

	baseQuery := `
    SELECT a.id, a.user_id, a.title, a.description, a.image_url, a.price, a.created_at, u.login
    FROM ads a
    JOIN users u ON a.user_id = u.id
  `

	var filters []string
	var args []interface{}
	argID := 1

	if params.MinPrice > 0 {
		filters = append(filters, fmt.Sprintf("a.price >= $%d", argID))
		args = append(args, params.MinPrice)
		argID++
	}
	if params.MaxPrice > 0 {
		filters = append(filters, fmt.Sprintf("a.price <= $%d", argID))
		args = append(args, params.MaxPrice)
		argID++
	}

	if len(filters) > 0 {
		baseQuery += " WHERE " + strings.Join(filters, " AND ")
	}

	orderBy := "a.created_at"
	switch params.SortBy {
	case "price":
		orderBy = "a.price"
	case "date":
		orderBy = "a.created_at"
	}

	orderDirection := "DESC"
	if strings.ToUpper(params.SortDir) == "ASC" {
		orderDirection = "ASC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDirection)

	limit := 10
	if params.Limit > 0 {
		limit = params.Limit
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d", argID)
	args = append(args, limit)
	argID++

	offset := 0
	if params.Page > 1 {
		offset = (params.Page - 1) * limit
	}
	baseQuery += fmt.Sprintf(" OFFSET $%d", argID)
	args = append(args, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query execution: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var ads []entity.AdWithAuthor
	for rows.Next() {
		var ad entity.AdWithAuthor
		if err := rows.Scan(
			&ad.ID,
			&ad.UserID,
			&ad.Title,
			&ad.Description,
			&ad.ImageURL,
			&ad.Price,
			&ad.CreatedAt,
			&ad.AuthorLogin,
		); err != nil {
			return nil, fmt.Errorf("%s: row scan: %w", op, err)
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

// Delete removes an ad by its ID
func (r AdsRepo) Delete(ctx context.Context, id int64) error {
	const op = "repository.AdsRepo.Delete"

	query := `DELETE FROM ads WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: check rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrAdNotFound)
	}

	return nil
}
