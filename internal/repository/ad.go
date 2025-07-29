package repository

import (
	"context"
	"database/sql"
	"fmt"
	"rest-api-marketplace/internal/entity"
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

func (r AdsRepo) GetById(ctx context.Context, id string) (*entity.Ad, error) {
	query := `SELECT id, user_id, title, description, image_url, price, created_at FROM ads WHERE id $1`

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

func (r AdsRepo) GetAll(ctx context.Context, params service.GetAdsParams) ([]entity.AdWithAuthor, error) {

}

func (r AdsRepo) Delete(ad entity.Ad) error {
	return nil
}
