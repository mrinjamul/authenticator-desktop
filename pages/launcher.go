package pages

import (
	"github.com/mrinjamul/authenticator-desktop/config"
	"github.com/mrinjamul/authenticator-desktop/constants"
	"github.com/mrinjamul/authenticator-desktop/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

	// create account list
	for _, account := range accounts {
		var card *widget.Card
		// var totp string
		totp := utils.GetTOTPToken(account.Secret)
		// item := widget.NewLabelWithStyle(totp, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		item := widget.NewButton("Copy", func() {
			utils.CopyToClipboard(totp)
		})
		card = widget.NewCard(account.Name+" | "+account.Username, totp, item)
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
