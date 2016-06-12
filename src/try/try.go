package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

const (
	i int = iota
	b int = iota
	c int = iota
)

func main() {
	fmt.Println(i, b, c)
	username, password := credentials()
	fmt.Printf("Username: %s, Password: %s\n", username, password)
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

codeToByte[IAC] = '\xff\xfb\x01'
codeToByte[WILL] = '\xfb'
codeToByte[WONT] = '\xfc'
codeToByte[ECHO] = '\x01'