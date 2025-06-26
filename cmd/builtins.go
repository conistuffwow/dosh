package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func StartShell() {
	reader := bufio.NewReader(os.Stdin)
	for {
		showPrompt()
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading input... ", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		
		args := strings.Fields(input)
		command := strings.ToUpper(args[0])
		restArgs := args[1:]
		switch command {
		case "EXIT":
			os.Exit(0)
		case "CLS":
			clearScreen()
		case "CD":
			changeDirectory(restArgs)
		case "DIR":
			listDirectory()
		default:
			if strings.HasSuffix(strings.ToLower(args[0]), ".bat") {
				err := runBatFile(args[0])
				if err != nil {
					fmt.Println("error running bat file:", err)
				}
			} else {
				runExternalCommand(args[0], restArgs)
			}
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear") // use cls on widnows at some point (ill probably never add windows support as windows already has a dos shell.)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func changeDirectory(args []string) {
	if len(args) == 0 {
		fmt.Println("usgage: CD <directory>")
		return
	}
	err := os.Chdir(args[0])
	if err != nil {
		fmt.Println("Invalid directory:", err)
	}
}

func listDirectory() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Println("[DIR] ", file.Name())
		} else {
			fmt.Println(file.Name())
		}
	}
}

func runExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running command:", err)
	}
}