package locales

import (
	"encoding/json"
	"errors"
	"os"
)

type Locale struct {
	Name string                       `json:"name"`
	Keys map[string]map[string]string `json:"keys"`
}

var (
	ErrNotValidLocale = errors.New("not valid locale file")
	ErrNoGroupFound   = errors.New("no group found")
	ErrNoKeyFound     = errors.New("no value for this key")
)

func (l *Locale) Load(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, l)
	if err != nil {
		return err
	}

	if l.Name == "" || len(l.Keys) == 0 {
		return ErrNotValidLocale
	}

	return nil
}

func (l *Locale) Save(filename string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (l *Locale) GetByKey(groupName, key string) (string, error) {
	var group map[string]string
	var value string
	var ok bool

	if group, ok = l.Keys[groupName]; !ok {
		return "", ErrNoGroupFound
	}
	if value, ok = group[key]; !ok {
		return "", ErrNoKeyFound
	}

	return value, nil
}

func (l *Locale) SetByKey(group, key, value string) {
	if l.Keys[group] == nil {
		l.Keys[group] = make(map[string]string)
	}
	l.Keys[group][key] = value
}

func (l *Locale) DeleteByKey(group, key string) error {
	_, err := l.GetByKey(group, key)
	if err != nil {
		return err
	}

	delete(l.Keys[group], key)
	return nil
}

func (l *Locale) DeleteGroup(groupName string) error {
	if _, ok := l.Keys[groupName]; !ok {
		return ErrNoGroupFound
	}

	delete(l.Keys, groupName)
	return nil
}
