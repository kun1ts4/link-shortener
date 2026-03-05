package usecase

import (
	"context"
	"errors"
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

func (s *ShortenerUseCase) CreateShort(ctx context.Context, original string) (*domain.Link, error) {
	existing, err := s.repo.FindByOriginal(ctx, original)
	if err == nil && existing != nil {
		return existing, nil
	}

	short := s.gen.Generate(original)
	link, err := domain.NewLink(original, short)
	if err != nil {
		return nil, fmt.Errorf("build link: %w", err)
	}

	err = s.repo.Create(ctx, link)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return s.repo.FindByShort(ctx, short)
		}
		return nil, fmt.Errorf("create link: %w", err)
	}

	return link, nil
}

func (s *ShortenerUseCase) GetOriginalLink(ctx context.Context, short string) (*domain.Link, error) {
	link, err := s.repo.FindByShort(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("find link: %w", err)
	}

	if err = s.repo.IncrementClicks(ctx, short); err != nil {
		return nil, fmt.Errorf("increment clicks: %w", err)
	}

	return link, nil
}
