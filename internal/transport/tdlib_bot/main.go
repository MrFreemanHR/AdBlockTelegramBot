package tdlibbot

import (
	"adblock_bot/internal/adapter/logger"
	"adblock_bot/internal/config"
	"fmt"
	"os"

	"github.com/Arman92/go-tdlib"
)

type tdlibbot struct {
	client *tdlib.Client
}

func New() *tdlibbot {
	tdlib.SetLogVerbosityLevel(1)
	tdlib.SetFilePath("./tdlib_errors.log")

	if createDirs() != nil {
		logger.Logger().Fatal("Can't start user-bot due to previous errors!")
		return nil
	}

	// Create new instance of client
	client := tdlib.NewClient(tdlib.Config{
		APIID:               config.CurrentConfig.APIid,
		APIHash:             config.CurrentConfig.APIHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
		UseMessageDatabase:  true,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseTestDataCenter:   false,
		DatabaseDirectory:   "./tdlib-db",
		FileDirectory:       "./tdlib-files",
		IgnoreFileNames:     false,
	})

	var td = tdlibbot{
		client: client,
	}

	if td.Auth() != nil {
		logger.Logger().Fatal("Can't start user-bot due to previous errors!")
		return nil
	}

	return &td
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

func (t *tdlibbot) Auth() error {
	logger.Logger().Info("=== Authorization for user-bot ===")
	fmt.Printf("Pointer: %#+v\n", t.client)
	for {
		currentState, _ := t.client.Authorize()
		if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPhoneNumberType {
			fmt.Print("> Enter phone: ")
			var number string
			fmt.Scanln(&number)
			_, err := t.client.SendPhoneNumber(number)
			if err != nil {
				logger.Logger().Error("Error while sending phone number: %s", err)
				return err
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("> Enter code: ")
			var code string
			fmt.Scanln(&code)
			_, err := t.client.SendAuthCode(code)
			if err != nil {
				logger.Logger().Error("Error while sending auth code: %s", err)
				return err
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("> Enter Password: ")
			var password string
			fmt.Scanln(&password)
			_, err := t.client.SendAuthPassword(password)
			if err != nil {
				logger.Logger().Error("Error while sending password: %s", err)
				return err
			}
		} else if currentState.GetAuthorizationStateEnum() == tdlib.AuthorizationStateReadyType {
			logger.Logger().Info("Authorization success!")
			break
		}
	}
	return nil
}
