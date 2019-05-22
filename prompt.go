package bhlclone

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation(deflt bool) bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response == "" {
		return deflt
	}
	if response == "y" || response == "yes" {
		return true
	}
	if response == "n" || response == "no" {
		return false
	}
	fmt.Println(`Please type "yes" or "no" and then press enter:`)
	return askForConfirmation(deflt)
}
