package random

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/kun1ts4/link-shortener/internal/config"
)

type LinkGen struct {
	alphabet []rune
	length   int
}

func NewLinkGen(cfg config.ShortenerConfig) *LinkGen {
	return &LinkGen{
		alphabet: []rune(cfg.Alphabet),
		length:   cfg.ShortLength,
	}
}

// алгоритм хеширует оригинальную ссылку, хеш конвертируется в 16-ричную систему
// 16-ричный хэш делится на двузначные числа которые по остатку от деления сопоставляются символам алфавита

func (g *LinkGen) Generate(original string) string {
	hash := sha256.Sum256([]byte(original))
	hashString := hex.EncodeToString(hash[:])
	result := make([]byte, g.length)

	for i := 0; i < g.length; i++ {
		pos := (i * 2) % len(hashString)

		firstVal := hexToInt(hashString[pos])
		secondVal := hexToInt(hashString[(pos+1)%len(hashString)])

		num := firstVal*16 + secondVal
		alphabetIndex := num % len(g.alphabet)
		result[i] = byte(g.alphabet[alphabetIndex])
	}
	return string(result)
}

func hexToInt(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	return int(c - 'a' + 10)
}
