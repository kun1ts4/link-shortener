package domain

type LinkGen interface {
	Generate(original string) string
}
