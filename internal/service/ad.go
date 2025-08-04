package service

import (
	"context"
	"fmt"
	"net/url"
	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/repository"
)

type adService struct {
	repo repository.Ads
}

func NewAdService(repo repository.Ads) *adService {
	return &adService{
		repo: repo,
	}
}

func (s adService) Create(ctx context.Context, input CreateAdInput, userId int64) (*entity.Ad, error) {
	if err := validateInput(input.Title, input.Description, input.ImageURL, input.Price); err != nil {
		return nil, err
	}

	ad := entity.Ad{
		UserID:      userId,
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	}

	adId, err := s.repo.Create(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("failed to create ad: %w", err)
	}

	return s.repo.GetById(ctx, adId)
}

func (s adService) Update(ctx context.Context, adId, userId int64, input UpdateAdInput) (*entity.Ad, error) {
	originalAd, err := s.repo.GetById(ctx, adId)
	if err != nil {
		return nil, fmt.Errorf("failed to get original ad: %w", err)
	}

	if originalAd.UserID != userId {
		return nil, entity.ErrForbidden
	}

	updatedAd := *originalAd

	if input.Title != nil {
		updatedAd.Title = *input.Title
	}
	if input.Description != nil {
		updatedAd.Description = *input.Description
	}
	if input.ImageURl != nil {
		updatedAd.ImageURL = *input.ImageURl
	}
	if input.Price != nil {
		updatedAd.Price = *input.Price
	}

	if err := validateInput(updatedAd.Title, updatedAd.Description, updatedAd.ImageURL, updatedAd.Price); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, adId, updatedAd); err != nil {
		return nil, fmt.Errorf("failed to update ad: %w", err)
	}

	return &updatedAd, nil
}

func (s adService) GetByID(ctx context.Context, id int64) (*entity.Ad, error) {
	ad, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ad by id: %w", err)
	}

	return ad, nil
}

func (s adService) GetAll(ctx context.Context, params entity.GetAdsQuery, currentUserId *int64) ([]entity.AdResponse, error) {
	adsWithAuthor, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get all ads: %w", err)
	}

	response := make([]entity.AdResponse, len(adsWithAuthor))

	for i, ad := range adsWithAuthor {
		res := entity.AdResponse{
			AdWithAuthor: ad,
		}

		if currentUserId != nil {
			isOwner := ad.UserID == *currentUserId
			res.IsOwner = &isOwner
		}

		response[i] = res

	}

	return response, nil
}

func (s adService) Delete(ctx context.Context, adId, userId int64) error {
	ad, err := s.repo.GetById(ctx, adId)
	if err != nil {
		return fmt.Errorf("failed to get by id for deleting: %w", err)
	}

	if ad.UserID != userId {
		return entity.ErrForbidden
	}
	if err := s.repo.Delete(ctx, adId); err != nil {
		return fmt.Errorf("failed to delete ad: %w", err)
	}

	return nil
}

func validateInput(title, description, imageURL string, price float64) error {
	if len(title) == 0 || len(title) > 100 {
		return fmt.Errorf("invalid title length")
	}
	if len(description) > 1000 {
		return fmt.Errorf("invalid description length")
	}
	if imageURL != "" {
		if _, err := url.ParseRequestURI(imageURL); err != nil {
			return fmt.Errorf("invalid umage url format")
		}
	}
	if price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return nil
}
