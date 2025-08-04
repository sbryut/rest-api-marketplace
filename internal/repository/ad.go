package repository

import (
	"context"
	"database/sql"
	"fmt"
	"rest-api-marketplace/internal/entity"
	"strings"
)

type AdsRepo struct {
	db *sql.DB
}

func NewAdsRepo(db *sql.DB) *AdsRepo {
	return &AdsRepo{db: db}
}

func (r AdsRepo) Create(ctx context.Context, ad entity.Ad) (int64, error) {
	query := `INSERT INTO ads (user_id, title, description, image_url, price) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int64

	err := r.db.QueryRowContext(ctx, query, ad.UserID, ad.Title, ad.Description, ad.ImageURL, ad.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create new ad: %w", err)
	}

	return id, nil
}

func (r AdsRepo) Update(ctx context.Context, id int64, ad entity.Ad) error {
	query := `UPDATE ads SET title = $1, description = $2, image_url = $3, price = $4 WHERE id = $5`

	res, err := r.db.ExecContext(ctx, query, ad.Title, ad.Description, ad.ImageURL, ad.Price, id)
	if err != nil {
		return fmt.Errorf("failed to update ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return entity.ErrAdNotFound
	}

	return nil
}

func (r AdsRepo) GetById(ctx context.Context, id int64) (*entity.Ad, error) {
	query := `SELECT id, user_id, title, description, image_url, price, created_at FROM ads WHERE id = $1`

	var ad entity.Ad

	err := r.db.QueryRowContext(ctx, query, id).Scan(&ad.ID, &ad.UserID, &ad.Title, &ad.Description, &ad.ImageURL, &ad.Price, &ad.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrAdNotFound
		}
		return nil, fmt.Errorf("failed to get ad by id: %w", err)
	}

	return &ad, nil
}

func (r AdsRepo) GetAll(ctx context.Context, params entity.GetAdsQuery) ([]entity.AdWithAuthor, error) {
	baseQuery := `
    SELECT a.id, a.user_id, a.title, a.description, a.image_url, a.price, a.created_at, u.login
    FROM ads a
    JOIN users u ON a.user_id = u.id
  `

	var filters []string
	var args []interface{}
	argId := 1

	if params.MinPrice > 0 {
		filters = append(filters, fmt.Sprintf("a.price >= $%d", argId))
		args = append(args, params.MinPrice)
		argId++
	}
	if params.MaxPrice > 0 {
		filters = append(filters, fmt.Sprintf("a.price <= $%d", argId))
		args = append(args, params.MaxPrice)
		argId++
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

	baseQuery += fmt.Sprintf(" LIMIT $%d", argId)
	args = append(args, limit)
	argId++

	offset := 0
	if params.Page > 1 {
		offset = (params.Page - 1) * limit
	}
	baseQuery += fmt.Sprintf(" OFFSET $%d", argId)
	args = append(args, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ads: %w", err)
	}
	defer rows.Close()

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
			return nil, fmt.Errorf("failed to get all ads: %w", err)
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

func (r AdsRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM ads WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete ad: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return entity.ErrAdNotFound
	}

	return nil
}
