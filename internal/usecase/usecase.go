package usecase

import (
	"context"
	"errors"
	"fmt"
	"link-shortener/internal/domain"
	"strconv"
)

type ShortenerUseCase struct {
	gen  domain.LinkGen
	repo domain.Repository
}

func NewShortenerUseCase(gen domain.LinkGen, repo domain.Repository) *ShortenerUseCase {
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

	// проверка и устранение коллизий

	attempt := 0
	for {
		found, foundErr := s.repo.FindByShort(ctx, short)
		if foundErr != nil && !errors.Is(foundErr, domain.ErrNotFound) {
			return nil, fmt.Errorf("collision check: %w", foundErr)
		}

		if errors.Is(foundErr, domain.ErrNotFound) {
			break
		}

		if found.Original == original {
			return found, nil
		}

		short = s.gen.Generate(original + strconv.Itoa(attempt))
		attempt++
	}

	link, err := domain.NewLink(original, short)
	if err != nil {
		return nil, fmt.Errorf("build link: %w", err)
	}

	err = s.repo.Create(ctx, link)
	if err != nil {
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
