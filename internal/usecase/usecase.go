package usecase

import (
	"fmt"
	"link-shortener/internal/domain"
	"link-shortener/internal/random"
)

type ShortenerUseCase struct {
	gen  random.LinkGen
	repo domain.Repository
}

func NewShortenerUseCase(gen random.LinkGen, repo domain.Repository) *ShortenerUseCase {
	return &ShortenerUseCase{
		gen:  gen,
		repo: repo,
	}
}

func (s *ShortenerUseCase) CreateShort(original string) (*domain.Link, error) {
	short := s.gen.Generate(original)
	link, err := domain.NewLink(original, short)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByOriginal(original)
	if err == nil && existing != nil {
		return existing, nil
	}

	err = s.repo.Create(link)
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return link, nil
}

//func (s *ShortenerUseCase) GetOriginalLink(key string) (*domain.Link, error) {
//
//}
