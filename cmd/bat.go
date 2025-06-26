package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func runBatFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Could not open file.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	labelMap := make(map[string]int)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines = append(lines, line)
		if strings.HasPrefix(line, ":") {
			label := strings.ToUpper(strings.TrimPrefix(line, ":"))
			labelMap[label] = len(lines) - 1
		}
	}

	variables := make(map[string]string)
	lineNum := 0

	for lineNum < len(lines) {
		rawLine := lines[lineNum]
		lineNum++

		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(strings.ToUpper(line), "REM") {
			continue // skip
		}

		// replace varz
		for k, v := range variables {
			line = strings.ReplaceAll(line, "%"+k+"%", v)
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])
		rest := args[1:]

		switch command {
		case "SET":
			if len(rest) >= 1 {
				parts := strings.SplitN(strings.Join(rest, " "), "=", 2)
				if len(parts) == 2 {
					key := strings.ToUpper(strings.TrimSpace(parts[0]))
					val := strings.TrimSpace(parts[1])
					variables[key] = val
				}
			}
		case "GOTO":
			if len(rest) == 1 {
				label := strings.ToUpper(rest[0])
				if pos, ok := labelMap[label]; ok {
					lineNum = pos
				} else {
					fmt.Println("label not found:", label)
				}
			}
		case "IF":
			if len(rest) >= 2 {
				condition := rest[0]
				if strings.Contains(condition, "==") {
					parts := strings.SplitN(condition, "==", 2)
					if len(parts) == 2 {
						// handle if variables, aaaaaaaaHHHHHHHHH THIS IS A DISGUSTING HACK
						leftRaw := strings.Trim(parts[0], "%")
						right := parts[1]
						left := strings.ToUpper((leftRaw))
						ifVal := variables[left]

						if ifVal == right && len(rest) > 1 {
							// execute rest of shit
							cmdParts := strings.Fields(strings.Join(rest[1:], " "))
							if len(cmdParts) == 0 {
								break
							}

							cmdName := strings.ToUpper(cmdParts[0])
							cmdArgs := cmdParts[1:]

							switch cmdName {
							case "ECHO":
								fmt.Println(strings.Join(cmdArgs, " "))
							case "CD":
								changeDirectory(cmdArgs)
							case "DIR":
								listDirectory()
							case "CLS":
								clearScreen()
							default:
								if strings.HasPrefix(cmdName, ":") {
									continue
								}
								runExternalCommand(cmdParts[0], cmdArgs)
							}
						}
					}
				}
			}
		case "ECHO":
			fmt.Println(strings.Join(rest, " "))
		case "EXIT":
			return nil
		case "CD":
			changeDirectory(rest)
		case "DIR":
			listDirectory()
		case "CLS":
			clearScreen()
		default:
			if strings.HasPrefix(command, ":") {
				continue
			}

			runExternalCommand(args[0], rest)
		}
	}
	return nil
}
