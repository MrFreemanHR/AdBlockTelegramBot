package parser

import (
	"errors"
	"strings"
)

type Token struct {
	chars []Char
}

var (
	ErrEmptyToken = errors.New("empty token")
)

func (t Token) IsEmpty() error {
	if t.Value() == "" {
		return ErrEmptyToken
	}
	return nil
}

func (t *Token) Value() string {
	var value string
	for _, c := range t.chars {
		value += string(c)
	}
	return strings.TrimSpace(value)
}

func (t *Token) InsertChar(c Char) {
	t.chars = append(t.chars, c)
}
