package main

import (
	"fmt"
	"io/ioutil"

	"github.com/mrinjamul/authenticator-cli/utils"
)

func main() {
	//Read the secret token from file system
	data, err := ioutil.ReadFile("secret.pem")
	utils.PanicIf(err)
	secret := string(data)
	otp := utils.GetTOTPToken(secret)
	fmt.Println(otp)
	utils.CopyToClipboard(otp)
}
