package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

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
