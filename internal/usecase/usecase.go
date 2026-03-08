package usecase

import (
	"context"
	"errors"
	"fmt"
	"link-shortener/internal/domain"
	"log/slog"
	"strconv"
)

type ShortenerUseCase struct {
	gen  domain.LinkGen
	repo domain.Repository
	log  *slog.Logger
}

func NewShortenerUseCase(gen domain.LinkGen, repo domain.Repository, log *slog.Logger) *ShortenerUseCase {
	return &ShortenerUseCase{
		gen:  gen,
		repo: repo,
		log:  log,
	}
}

func (s *ShortenerUseCase) CreateShort(ctx context.Context, original string) (*domain.Link, error) {
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

		s.log.Debug("collision found", "short", short, "attempt", attempt)
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

	s.log.Info("short link created", "short", link.Short)
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

	s.log.Debug("redirect", "short", short, "original", link.Original)
	return link, nil
}
