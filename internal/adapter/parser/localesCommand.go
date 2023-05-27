package parser

import (
	"golang.org/x/exp/slices"
)

/*
	Locales:

	/locales key     show       group          [key]
	cmd	     subcmd  subcmd     required_args  optional_args

	/locales key     add        group, key, value
	cmd	     subcmd  subcmd     required_args      optional_args
*/

func (p *Parser) parseLocalesCommand(c *Cmd, str string) error {
	var subcmdsFirstList = []string{
		"key",
	}
	var subcmdsSecondList = []string{
		"show",
		"add",
	}

	t, nextLevel, tokenErr := p.nextToken(str)
	if tokenErr != nil {
		return ErrParsingSubCmd
	}
	if !slices.Contains(subcmdsFirstList, t.Value()) {
		return ErrNotKnownSubCommand
	}
	var firstSubCmd = NewSubCmd(t.Value())
	c.ChildSubcmd = &firstSubCmd

	t, nextLevel, tokenErr = p.nextToken(nextLevel)
	if tokenErr != nil {
		return ErrParsingSubCmd
	}
	if !slices.Contains(subcmdsSecondList, t.Value()) {
		return ErrNotKnownSubCommand
	}
	var secondSubCmd = NewSubCmd(t.Value())
	firstSubCmd.Child = &secondSubCmd
	err := p.localesKeySubcmd(&secondSubCmd, nextLevel)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) localesKeySubcmd(c *Subcmd, str string) error {
	if c.Name == "show" {
		group, nextLevel, tokenErr := p.nextToken(str)
		if tokenErr != nil {
			return ErrParsingRequiredArg
		}
		var groupArg RequiredArg
		groupArg.SetValue(group.Value())
		c.RequiredArgs[0] = groupArg
		if nextLevel != "" {
			key, _, tokenErr := p.nextToken(nextLevel)
			if tokenErr != nil {
				return ErrParsingOptionalArg
			}
			if key.IsEmpty() == nil {
				var keyArg OptionalArg
				keyArg.SetValue(key.Value())
				c.OptionalArgs[0] = keyArg
			}
		}
	} else if c.Name == "add" {
		group, nextLevel, tokenErr := p.nextToken(str)
		if tokenErr != nil {
			return ErrParsingRequiredArg
		}
		var groupArg RequiredArg
		groupArg.SetValue(group.Value())
		c.RequiredArgs[0] = groupArg

		var key Token
		key, nextLevel, tokenErr = p.nextToken(nextLevel)
		if tokenErr != nil {
			return ErrParsingRequiredArg
		}
		var keyArg RequiredArg
		keyArg.SetValue(key.Value())
		c.RequiredArgs[1] = keyArg

		var value Token
		value, _, tokenErr = p.nextToken(nextLevel)
		if tokenErr != nil {
			return ErrParsingRequiredArg
		}
		var valueArg RequiredArg
		valueArg.SetValue(value.Value())
		c.RequiredArgs[2] = valueArg
	}
	return nil
}
