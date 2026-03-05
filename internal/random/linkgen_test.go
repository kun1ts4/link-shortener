package random

import (
	"link-shortener/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	gen := NewLinkGen(config.ShortenerConfig{
		Alphabet:    "abc012",
		ShortLength: 6,
	})

	url := "https://google.com"
	r1 := gen.Generate(url)
	r2 := gen.Generate(url)
	assert.Equal(t, r1, r2, "should be same result on same link")

	assert.NotEqual(t,
		gen.Generate("https://google.com"),
		gen.Generate("https://yandex.ru"),
		"should br different results on different link",
	)

	assert.Len(t, gen.Generate("https://google.com"), 6, "wrong length")
}

func TestAlphabet(t *testing.T) {
	alphabet := "abc012"
	gen := NewLinkGen(config.ShortenerConfig{
		Alphabet:    alphabet,
		ShortLength: 4,
	})

	result := gen.Generate("https://google.com")

	for _, c := range result {
		assert.Contains(t, alphabet, string(c), "wrong character: %c", c)
	}
}

func TestHexToInt(t *testing.T) {
	tests := []struct {
		name string
		c    byte
		want int
	}{
		{"hex 0", '0', 0},
		{"hex 5", '5', 5},
		{"hex 9", '9', 9},
		{"hex a", 'a', 10},
		{"hex f", 'f', 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hexToInt(tt.c)
			assert.Equal(t, tt.want, got)
		})
	}
}
