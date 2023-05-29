package parser

import (
	"golang.org/x/exp/slices"
)

/*
	Verifier:

	/verifier add    name        author:nickname ...
	cmd       subcmd require_arg require_arg     optional_args

	/verifier remove name
	cmd       subcmd require_arg

	/verifier list
	cmd       subcmd
*/

func (p *Parser) parseVerifierCommand(c *Cmd, str string) error {
	var subcmdsFirstList = []string{
		"add",
		"remove",
		"list",
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

	if firstSubCmd.Name == "add" {
		return p.verifierAddSubcmd(&firstSubCmd, nextLevel)
	}
	if firstSubCmd.Name == "remove" {
		return p.verifierRemoveSubcmd(&firstSubCmd, nextLevel)
	}
	if firstSubCmd.Name == "list" {
		return nil
	}

	return ErrNotKnownSubCommand
}

func (p *Parser) verifierAddSubcmd(c *Subcmd, str string) error {
	name, nextLevel, tokenErr := p.nextToken(str)
	if tokenErr != nil {
		return ErrParsingRequiredArg
	}
	var nameArg RequiredArg
	nameArg.SetValue(name.Value())
	c.RequiredArgs[0] = nameArg

	author, nextLevel, tokenErr := p.nextToken(nextLevel)
	if tokenErr != nil {
		return ErrParsingRequiredArg
	}
	var authorArg RequiredArg
	authorArg.SetValue(author.Value())
	c.RequiredArgs[1] = authorArg

	var optionalArgCounter = 0
	for nextLevel != "" {
		if optionalArgCounter > 5 { // Max count of args
			break
		}
		var optionalArgToken Token
		optionalArgToken, nextLevel, tokenErr = p.nextToken(nextLevel)
		if tokenErr != nil {
			break // Only "empty token" error
		}
		var optionalArg OptionalArg
		optionalArg.SetValue(optionalArgToken.Value())
		c.OptionalArgs[byte(optionalArgCounter)] = optionalArg
		optionalArgCounter++
	}

	return nil
}

func (p *Parser) verifierRemoveSubcmd(c *Subcmd, str string) error {
	name, _, tokenErr := p.nextToken(str)
	if tokenErr != nil {
		return ErrParsingRequiredArg
	}
	var nameArg RequiredArg
	nameArg.SetValue(name.Value())
	c.RequiredArgs[0] = nameArg

	return nil
}
