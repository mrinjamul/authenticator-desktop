package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//PanicIf error is not nil
func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

// Ptrfix0 Append extra 0s if the length of otp is less than 6
//If otp is "1234", it will return it as "001234"
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
