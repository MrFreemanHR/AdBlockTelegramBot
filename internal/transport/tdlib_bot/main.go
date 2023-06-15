package tdlibbot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"adblock_bot/internal/core/entity"
	"adblock_bot/internal/core/interfaces"
	"adblock_bot/internal/transport"
	"os"
	"strconv"

	"github.com/zelenin/go-tdlib/client"
)

type tdlibbot struct {
	client          *client.Client
	me              *client.User
	events          *client.Listener
	messageHandlers []interfaces.MessageHandler
}

func New() *tdlibbot {
	err := createDirs()
	if err != nil {
		logger.Logger().Error("Can't create cache dirs: %s", err.Error())
		return nil
	}

	authorizer := client.ClientAuthorizer()
	go client.CliInteractor(authorizer)

	apiid, err := strconv.Atoi(config.CurrentConfig.APIid)
	if err != nil {
		logger.Logger().Error("Can't convert APP_ID to int: %s", err.Error())
		return nil
	}
	authorizer.TdlibParameters <- &client.TdlibParameters{
		UseTestDc:              false,
		DatabaseDirectory:      "./tdlib-db",
		FilesDirectory:         "./tdlib-files",
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		ApiId:                  int32(apiid),
		ApiHash:                config.CurrentConfig.APIHash,
		SystemLanguageCode:     "en",
		DeviceModel:            "server",
		SystemVersion:          "1.0.0",
		ApplicationVersion:     "1.0.0",
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	}

	_, err = client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		logger.Logger().Error("Can't change verbosity level: %s", err.Error())
		return nil
	}

	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		logger.Logger().Error("Can't create new tdlib client: %s", err.Error())
		return nil
	}

	me, err := tdlibClient.GetMe()
	if err != nil {
		logger.Logger().Error("Can't get myself: %s", err.Error())
		return nil
	}

	events := tdlibClient.GetListener()
	return &tdlibbot{
		client: tdlibClient,
		me:     me,
		events: events,
	}
}

func createDirs() error {
	if _, err := os.Stat("./tdlib-db"); os.IsNotExist(err) {
		err = os.Mkdir("./tdlib-db", os.ModePerm)
		if err != nil {
			logger.Logger().Error("Can't create tdlib-db directory: %s", err.Error())
			return err
		}
	}

	if _, err := os.Stat("./tdlib-files"); os.IsNotExist(err) {
		err = os.Mkdir("./tdlib-files", os.ModePerm)
		if err != nil {
			logger.Logger().Error("Can't create tdlib-files directory: %s", err.Error())
			return err
		}
	}

	return nil
}

func (t *tdlibbot) SetMessageHandlers(messageHandlers []interfaces.MessageHandler) {
	t.messageHandlers = messageHandlers
}

func (t *tdlibbot) GetTelegramAPI() *transport.TelegramAPI {
	self, err := t.client.GetMe()
	if err != nil {
		logger.Logger().Fatal("[TD LIB] Can't get me: %s", err.Error())
		return nil
	}
	return transport.NewUnifiedTransportWithTDlib(
		t.client,
		entity.TelegramUser{
			ID:           self.Id,
			IsBot:        self.Type.UserTypeType() == client.TypeUserTypeBot,
			FirstName:    self.FirstName,
			LastName:     self.LastName,
			UserName:     self.Username,
			LanguageCode: self.LanguageCode,
		},
	)
}

func (t *tdlibbot) Run() {
	for update := range t.events.Updates {
		if update.GetClass() == client.ClassUpdate {
			t.ProcessEvents(update)
		}
	}
}
