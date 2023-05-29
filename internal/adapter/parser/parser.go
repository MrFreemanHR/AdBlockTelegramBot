package parser

import (
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type Parser struct {
}

var (
	ErrNotKnownCommand    = errors.New("not known command")
	ErrNoCommand          = errors.New("no command")
	ErrParsingSubCmd      = errors.New("can't parsing subcommands")
	ErrNotKnownSubCommand = errors.New("not know subcommand")
	ErrParsingRequiredArg = errors.New("can't parsing required argument")
	ErrParsingOptionalArg = errors.New("can't parsing optional argument")
)

func (p *Parser) Parse(command string) (*Cmd, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, ErrNoCommand
	}
	if command[0] == '/' {
		command = command[1:]
	}
	return p.parseCmd(command)
}

func (p *Parser) commandsParsers() map[string]func(*Cmd, string) error {
	return map[string]func(*Cmd, string) error{
		"locales": p.parseLocalesCommand,
	}
}

func (p *Parser) parseCmd(command string) (*Cmd, error) {
	token, str, err := p.nextToken(command)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(GetCommandNames(), token.Value()) {
		return nil, ErrNotKnownCommand
	}
	cmd := NewCmd(token.Value())
	err = p.commandsParsers()[token.Value()](&cmd, str)
	if err != nil {
		return nil, err
	}
	return &cmd, nil
}

func (p *Parser) nextToken(str string) (Token, string, error) {
	var token Token
	var char rune
	var i int
	var inQuotes bool

	if len(str) > 0 && Char(str[0]).IsQuote() {
		inQuotes = true
	}
	var previousChar = Char(' ')
	for i, char = range str {
		currentChar := Char(char)
		if currentChar.IsWhitespace() && !inQuotes {
			break
		}
		if currentChar.IsQuote() && inQuotes {
			if i == 0 {
				continue
			}
			if previousChar.IsBackSlash() {
				token.InsertChar(currentChar)
				continue
			}
			break
		}
		if !currentChar.IsBackSlash() {
			token.InsertChar(currentChar)
		}
		previousChar = currentChar
	}
	if empty := token.IsEmpty(); empty != nil {
		return token, "", empty
	}
	if len(str) == i+1 {
		return token, "", nil
	}
	return token, strings.TrimSpace(str[i:]), nil
}
