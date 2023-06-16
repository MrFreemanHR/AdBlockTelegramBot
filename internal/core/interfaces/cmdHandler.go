package interfaces

import "adblock_bot/internal/adapter/parser"

type CmdHandler interface {
	ProcessCommand(cmd *parser.Cmd) string
}
