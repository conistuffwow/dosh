package cmd

import (
	"fmt"
	"os"
	"strings"
)

func showPrompt() {
	cwd, _ := os.Getwd()
	fmt.Printf("C:%s> ", strings.ReplaceAll(cwd, "/", "\\"))
}