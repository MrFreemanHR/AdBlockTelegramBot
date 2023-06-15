package parser

import "errors"

type RequiredArg struct {
	value string
}

var ErrEmptyRequiredValue = errors.New("empty required value")

func (a *RequiredArg) SetValue(value string) error {
	if value == "" {
		return ErrEmptyRequiredValue
	}
	a.value = value
	return nil
}

func (a *RequiredArg) GetValue() string {
	return a.value
}
