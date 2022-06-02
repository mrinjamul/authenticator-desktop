package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mrinjamul/authenticator-desktop/models"
	"github.com/rs/xid"
)

// PanicIf error is not nil
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// Ptrfix0 Append extra 0s if the length of otp is less than 6
// If otp is "1234", it will return it as "001234"
func prefix0(otp string) string {
	if len(otp) == 6 {
		return otp
	}
	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}
	return otp
}

// GetHOTPToken returns the HOTP token
func GetHOTPToken(secret string, interval int64) string {

	//Converts secret to base32 Encoding. Base32 encoding desires a 32-character
	//subset of the twenty-six letters A–Z and ten digits 0–9
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	PanicIf(err)
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(interval))

	//Signing the value using HMAC-SHA1 Algorithm
	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)

	PanicIf(err)
	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000

	//Converts number as a string
	otp := strconv.Itoa(int(h12))

	return prefix0(otp)
}

// GetTOTPToken returns the TOTP token
func GetTOTPToken(secret string) string {
	//The TOTP token is just a HOTP token seeded with every 30 seconds.
	interval := time.Now().Unix() / 30
	return GetHOTPToken(secret, interval)
}

// Copy to clipboard function for Linux, Windows, MacOS
func CopyToClipboard(text string) {
	if runtime.GOOS == "windows" {
		copyWindows(text)
	}
	if runtime.GOOS == "linux" {
		copyLinux(text)
	}
	if runtime.GOOS == "darwin" {
		copyMac(text)
	}
}

func copyLinux(text string) {
	cmd := exec.Command("xsel", "-i", "-b")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

func copyWindows(text string) {
	cmd := exec.Command("clip")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

func copyMac(text string) {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	cmd.Run()
}

// GetConfigFile returns the config file path
func GetConfig() string {
	// Get Home directory
	home, err := os.UserHomeDir()
	PanicIf(err)

	configDir := filepath.Join(home, ".config/authenticator")

	configFile := filepath.Join(configDir, "config.json")

	// check if configFile exists or not
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// create configFile
		err = os.MkdirAll(configDir, os.ModePerm)
		PanicIf(err)
		// create configFile
		f, err := os.Create(configFile)
		PanicIf(err)
		defer f.Close()
	}
	return configFile
}

// SaveAccounts save all accounts
func SaveAccounts(accounts []models.Account, configFile string) error {
	ByteAccounts, err := json.Marshal(accounts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, ByteAccounts, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// ReadAccounts decrypt all accounts
func ReadAccounts(configFile string) ([]models.Account, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return []models.Account{}, err
	}
	var accounts []models.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return []models.Account{}, err
	}

	return accounts, nil
}

// GenTips generates random tips
func GenTips() string {
	tips := []string{
		"Don't use your authenticator to login to your account",
	}
	rand.Seed(time.Now().UnixNano())
	return tips[rand.Intn(len(tips))]
}

// GenerateID generates a random ID
func GenerateID() string {
	id := xid.New().String()
	return id
}

// ConfirmPrompt will prompt to user for yes or no
func ConfirmPrompt(message string) bool {
	var response string
	fmt.Print(message + " (yes/no) :")
	fmt.Scanln(&response)

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return false
	}
}

// AddAccount add new account to the db
func AddAccount(name, username, email, secret string) error {
	accounts, err := ReadAccounts(GetConfig())
	if err != nil {
		return err
	}
	// check if valid base32 secret
	if !IsBase32(secret) {
		log.Println("invalid secret")
		return errors.New("invalid secret")
	}
	accounts = append(accounts, models.Account{
		ID:       GenerateID(),
		Name:     name,
		Username: username,
		Email:    email,
		Secret:   secret,
	})
	err = SaveAccounts(accounts, GetConfig())
	return err
}

// RemoveAccount remove account from the db
func RemoveAccount(id string) {
	accounts, err := ReadAccounts(GetConfig())
	PanicIf(err)
	for i, account := range accounts {
		if account.ID == id {
			accounts = append(accounts[:i], accounts[i+1:]...)
			break
		}
	}
	err = SaveAccounts(accounts, GetConfig())
	PanicIf(err)
}

// IsBase32 checks if a string is a valid base32 string
func IsBase32(str string) bool {
	_, err := base32.StdEncoding.DecodeString(str)
	return err == nil
}
