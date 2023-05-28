package localeshandler

import (
	"adblock_bot/internal/adapter/parser"
	"adblock_bot/internal/core/interfaces"
)

type cmdhandler struct {
}

func New() interfaces.CmdHandler {
	return &cmdhandler{}
}

func (h *cmdhandler) ProcessCommand(cmd *parser.Cmd) string {
	return "result"
}
