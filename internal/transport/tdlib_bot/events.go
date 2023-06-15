package tdlibbot

import (
	"adblock_bot/internal/adapter/logger"

	"github.com/Arman92/go-tdlib"
)

func (b *tdlibbot) ProcessEvents(event tdlib.UpdateMsg) {
	logger.Logger().UselessInfo("Data: %#+v\n", event.Data)
}
