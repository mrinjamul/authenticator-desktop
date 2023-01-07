package pages

import (
	"time"

	"github.com/mrinjamul/authenticator-desktop/config"
	"github.com/mrinjamul/authenticator-desktop/constants"
	"github.com/mrinjamul/authenticator-desktop/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type launcher struct {
	config config.Config
}

func (page *launcher) HashCode() string {
	return constants.PAGE_LAUNCHER_KEY
}

func (page *launcher) Render() {
	items := []fyne.CanvasObject{}
	// load config from file
	config := utils.GetConfig()
	// load accounts from file
	accounts, err := utils.ReadAccounts(config)
	if err != nil {
		panic(err)
	}

	// dynamically generate OTP every 30 seconds using data/binding fyne
	// create account list
	for _, account := range accounts {
		// var totp string
		totp := binding.NewString()
		go func() {
			// Generate OTP when time second is 00 and 30
			totp.Set(utils.GetTOTPToken(account.Secret))
			for {
				// Get current second
				second := time.Now().UTC().Second()
				// If second is 00 or 30 then generate OTP
				if second == 0 || second == 30 {
					// Generate OTP
					totp.Set(utils.GetTOTPToken(account.Secret))
				}
				// Wait for 1 second
				time.Sleep(time.Second)
			}
		}()
		// item := widget.NewLabelWithStyle(totp, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		item := widget.NewButton("Copy", func() {
			utils.CopyToClipboardWithBinding(totp)
		})
		// Create a card with account name, totp code and the button
		otp := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		otp.Bind(totp)
		card := container.NewVBox(
			widget.NewLabelWithStyle(account.Name+" | "+account.Username, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			otp,
			item,
		)
		items = append(items, card)
	}

	// Add new account button
	addAccountButton := widget.NewButton("Add Account", func() {
		page.config.Launch()(constants.PAGE_ACCOUNT_PAGE_KEY)
	})
	// create centered container for button
	addAccountButtonContainer := container.NewCenter(addAccountButton)
	accountList := container.NewVScroll(container.NewVBox(items...))
	accountList.SetMinSize(fyne.NewSize(400, 400))
	// Add create a container and add button to it
	c := container.NewVBox(accountList, addAccountButtonContainer)
	cc := container.NewHBox(c)
	page.config.GetWindow().SetContent(cc)
}

func NewLauncher(conf config.Config) Page {
	return &launcher{
		config: conf,
	}
}
