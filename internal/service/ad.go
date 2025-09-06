package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"rest-api-marketplace/internal/entity"
	"rest-api-marketplace/internal/repository"
)

// AdService provides operations to manage ads
type AdService struct {
	repo   repository.Ads
	logger *slog.Logger
}

// NewAdService creates a new AdService instance
func NewAdService(repo repository.Ads, logger *slog.Logger) *AdService {
	return &AdService{
		repo:   repo,
		logger: logger,
	}
}

// Create validates input and creates a new ad
func (s AdService) Create(ctx context.Context, input CreateAdInput, userID int64) (*entity.Ad, error) {
	const op = "service.AdService.Create"

	if err := validateInput(input.Title, input.Description, input.ImageURL, input.Price); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ad := entity.Ad{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
	}

	adID, err := s.repo.Create(ctx, ad)
	if err != nil {
		s.logger.Error("failed to create ad", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s.repo.GetByID(ctx, adID)
}

// Update modifies an existing ad with new data
func (s AdService) Update(ctx context.Context, adID, userID int64, input UpdateAdInput) (*entity.Ad, error) {
	const op = "service.AdService.Update"

	originalAd, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to get original ad", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if originalAd.UserID != userID {
		return nil, entity.ErrForbidden
	}

	updatedAd := *originalAd

	if input.Title != nil {
		updatedAd.Title = *input.Title
	}
	if input.Description != nil {
		updatedAd.Description = *input.Description
	}
	if input.ImageURL != nil {
		updatedAd.ImageURL = *input.ImageURL
	}
	if input.Price != nil {
		updatedAd.Price = *input.Price
	}

	if err := validateInput(updatedAd.Title, updatedAd.Description, updatedAd.ImageURL, updatedAd.Price); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := s.repo.Update(ctx, adID, updatedAd); err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to update ad", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedAd, nil
}

// GetByID retrieves an ad by its ID
func (s AdService) GetByID(ctx context.Context, id int64) (*entity.Ad, error) {
	const op = "service.AdService.GetByID"

	ad, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to get ad by id", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ad, nil
}

// GetByIDWithAuthor retrieves an ad with author information and optional ownership info
func (s AdService) GetByIDWithAuthor(ctx context.Context, id int64, currentUserID *int64) (*entity.AdResponse, error) {
	const op = "service.AdService.GetByIDWithAuthor"

	ad, err := s.repo.GetByIDWithAuthor(ctx, id)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to get ad with author by id", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res := &entity.AdResponse{
		AdWithAuthor: *ad,
	}

	if currentUserID != nil {
		isOwner := ad.UserID == *currentUserID
		res.IsOwner = &isOwner
	}

	return res, nil
}

// GetAll returns all ads with author info and optional ownership info
func (s AdService) GetAll(ctx context.Context, params entity.GetAdsQuery, currentUserID *int64) ([]entity.AdResponse, error) {
	const op = "service.AdService.GetAll"

	adsWithAuthor, err := s.repo.GetAll(ctx, params)
	if err != nil {
		s.logger.Error("failed to get all ads", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	response := make([]entity.AdResponse, len(adsWithAuthor))

	for i, ad := range adsWithAuthor {
		res := entity.AdResponse{
			AdWithAuthor: ad,
		}

		if currentUserID != nil {
			isOwner := ad.UserID == *currentUserID
			res.IsOwner = &isOwner
		}

		response[i] = res
	}

	return response, nil
}

// Delete removes an ad if the user is authorized
func (s AdService) Delete(ctx context.Context, adID, userID int64) error {
	const op = "service.AdService.Delete"

	ad, err := s.repo.GetByID(ctx, adID)
	if err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to get ad by id for delete", slog.String("op", op), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if ad.UserID != userID {
		return fmt.Errorf("%s: %w", op, entity.ErrForbidden)
	}
	if err := s.repo.Delete(ctx, adID); err != nil {
		if errors.Is(err, entity.ErrAdNotFound) {
			return fmt.Errorf("%s: %w", op, err)
		}
		s.logger.Error("failed to delete ad", slog.String("op", op), slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// validateInput checks if ad fields are correct
func validateInput(title, description, imageURL string, price float64) error {
	if len(title) < 1 || len(title) > 100 {
		return fmt.Errorf("title length must be between 1 and 100: %w", entity.ErrInvalidInput)
	}
	if len(description) > 1000 {
		return fmt.Errorf("description length must be less than 1000: %w", entity.ErrInvalidInput)
	}
	if imageURL != "" {
		if _, err := url.ParseRequestURI(imageURL); err != nil {
			return fmt.Errorf("invalid umage url format: %w", entity.ErrInvalidInput)
		}
	}
	if price < 0 {
		return fmt.Errorf("price cannot be negative: %w", entity.ErrInvalidInput)
	}
	return nil
}
