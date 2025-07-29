package service

import (
	"context"
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

func (s adService) GetByID(ctx context.Context, id string) entity.Ad {
	return s.repo.GetOne(id)
}

func (s adService) GetAll(ctx context.Context) []entity.Ad {
	return s.repo.GetAll(ctx)
}

/*func (s adService) Delete(ctx context.Context) error {
	return s.storage.Delete()
}*/
