package pages

import (
	"log"

	"github.com/mrinjamul/authenticator-desktop/constants"

	"github.com/mrinjamul/authenticator-desktop/config"
	"github.com/mrinjamul/authenticator-desktop/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type editPage struct {
	config config.Config
}

func (page *editPage) HashCode() string {
	return constants.PAGE_EDIT_ACCOUNT_KEY
}

func (page *editPage) Render() {
	content := []fyne.CanvasObject{}

	// Create Title Label
	titleLabel := widget.NewLabelWithStyle("Edit Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	// Create container for title label
	titleLabelContainer := container.NewCenter(titleLabel)

	// create add account form fields
	nameField := widget.NewEntry()
	nameField.SetPlaceHolder("John Doe")
	usernameField := widget.NewEntry()
	usernameField.SetPlaceHolder("john")
	emailField := widget.NewEntry()
	emailField.SetPlaceHolder("john@example.com")
	secretField := widget.NewEntry()
	secretField.SetPlaceHolder("Enter base32 encoded secret")

	// create add account button
	addAccountButton := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		// check if all fields are filled
		if nameField.Text == "" || (usernameField.Text == "" && emailField.Text == "") || secretField.Text == "" {
			log.Println("Form not filled")
			return
		}
		// save account to file
		err := utils.AddAccount(
			nameField.Text,
			usernameField.Text,
			emailField.Text,
			secretField.Text,
		)
		if err != nil {
			log.Println("Error saving account:", err)
			return
		}
		// go back to accounts page
		page.config.Launch()(constants.PAGE_LAUNCHER_KEY)
	})
	// create back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		page.config.Launch()(constants.PAGE_LAUNCHER_KEY)
	})

	// create container for button
	buttons := container.NewHBox(addAccountButton, backButton)
	// create centered container for button
	buttonContainer := container.NewCenter(buttons)

	// create form container
	form := container.NewVBox(
		titleLabelContainer,
		widget.NewLabel("Name"),
		nameField,
		widget.NewLabel("Username"),
		usernameField,
		widget.NewLabel("Email"),
		emailField,
		widget.NewLabel("Secret"),
		secretField,
		buttonContainer,
	)

	content = append(content, form)

	c := container.NewGridWrap(fyne.NewSize(400, 100), content...)
	// c := container.NewGridWithRows(2, items...)
	cc := container.NewHBox(c)
	page.config.GetWindow().SetContent(cc)
}

func NewEditAccountPage(conf config.Config) Page {
	return &editPage{
		config: conf,
	}
}
