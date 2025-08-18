package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// WITH PIPED DATA
	if isInputFromPipe() {
		if len(os.Args) > 1 {
			CASE := os.Args[1]
			LABEL := os.Args[2]
			ensureCaseDir(CASE)
			dbpath := ensureDB(CASE)
			saveRecord(dbpath, LABEL, "", "FILE", "N")
			dataHandler(CASE, LABEL, os.Stdin, os.Stdout)
		} else {
			fmt.Println("USAGE: | keeptrak CASE LABEL")
		}
	} else {
		if len(os.Args) == 2 {
			// handle help
			if os.Args[1] == "--help" {
				fmt.Println("")
				fmt.Println("Run Nested Shell: keeptrak")
				fmt.Println("")
				fmt.Println("Save Record: keeptrak CASE LABEL VALUE DATA_TYPE CONFIRMED")
				fmt.Println("\tExample: keeptrak johndoe username jdoe credential Y")
				fmt.Println("")
				fmt.Println("Save Note: keeptrak note TEXT")
				fmt.Println("\tExample: keeptrak note \"This is useful information\"")
				fmt.Println("")
				fmt.Println("Pipe data: keeptrak CASE LABEL")
				fmt.Println("\tExample: nmap 192.168.88.1 | keeptrak case103 nmap")
				fmt.Println("")
			} else {
				// handle unknown command
				fmt.Println("Unknown command: " + os.Args[1])
			}
		} else if len(os.Args) > 2 {
			CASE := os.Args[1]
			ensureCaseDir(CASE)
			if len(os.Args) == 4 {
				// handle note command
				if os.Args[2] == "note" {
					saveLineToFile(filepath.Join(CASE, "notes"), getTime()+"\t"+os.Args[3])
				}
			} else if len(os.Args) < 6 {
				fmt.Println("Too few arguments. Run --help to see correct usage")
			} else {
				// save to csv db
				LABEL := os.Args[2]
				VALUE := os.Args[3]
				DATATYPE := os.Args[4]
				CONFIRMED := os.Args[5]
				dbpath := ensureDB(CASE)
				saveRecord(dbpath, LABEL, VALUE, DATATYPE, CONFIRMED)
			}
		} else {
			// WRAPPED SHELL
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("")
			fmt.Println("K E E P T R A K   " + getTime())
			fmt.Println("")
			fmt.Print("Enter Case Name: ")
			text, _ := reader.ReadString('\n')
			fmt.Println("")
			// convert CRLF to LF
			CASE := strings.Replace(text, "\n", "", -1)
			ensureCaseDir(CASE)
			ensureDB(CASE)

			for {
				fmt.Print("KEEPTRAK> ")
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				command := strings.Replace(text, "\n", "", -1)
				if command != "" {
					// handle exit command
					if command == "exit" {
						return
					}
					// handle note command
					if strings.HasPrefix(command, "note: ") {
						saveLineToFile(filepath.Join(CASE, "notes"), getTime()+"\t"+command[6:])
						continue
					}
					// handle nested bash
					saveLineToFile(filepath.Join(CASE, "history"), getTime()+"\t"+command)
					cmd := exec.Command("bash", "-c", command)
					pipe, _ := cmd.StdoutPipe()
					if err := cmd.Start(); err != nil && err != io.EOF {
						fmt.Println(err)
					}
					reader := bufio.NewReader(pipe)
					line, err := reader.ReadString('\n')
					LABEL := strings.Split(command, " ")[0]
					totalOutput := ""
					lineCount := 0
					// loop through each line of output
					for err == nil {
						totalOutput += line
						lineCount++
						line = strings.ReplaceAll(line, "\n", "")
						fmt.Println(line)
						saveLineToFile(filepath.Join(CASE, LABEL), line)
						saveLineToFile(filepath.Join(CASE, "dump"), line)
						line, err = reader.ReadString('\n')
					}
					saveCommandSignature(CASE, LABEL, command, totalOutput, lineCount)
					if err != nil && err != io.EOF {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func dataHandler(CASE string, LABEL string, r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		saveLineToFile(filepath.Join(CASE, LABEL), text)
		saveLineToFile(filepath.Join(CASE, "dump"), text)
	}
}

func saveLineToFile(filePath string, content string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Failed to open to file", filePath)
	} else {
		content = stripColorCodes(content)
		_, err := f.WriteString(content + "\n")
		if err != nil {
			fmt.Println("Failed to write to file", filePath)
		}
	}
	if err := f.Close(); err != nil {
		log.Fatal("Failed to close file: "+filePath, err)
	}
}

func ensureCaseDir(CASE string) {
	_, err := os.Stat(CASE)
	if os.IsNotExist(err) {
		os.Mkdir(CASE, os.FileMode(0777))
	} else if err != nil {
		panic(err)
	}
}

func ensureDB(CASE string) string {
	dbpath := filepath.Join(CASE, "db.csv")
	_, err := os.Stat(dbpath)
	if os.IsNotExist(err) {
		saveLineToFile(dbpath, "LABEL,VALUE,DATA_TYPE,CONFIRMED,DATE_ADDED")
	}
	return dbpath
}

func getTime() string {
	currentTime := time.Now()
	return currentTime.Format("2006.01.02 15:04:05")
}

func saveRecord(dbpath string, LABEL string, VALUE string, DATATYPE string, CONFIRMED string) {
	saveLineToFile(dbpath, LABEL+","+VALUE+","+DATATYPE+","+CONFIRMED+","+getTime())
}

func stripColorCodes(str string) string {
	codePrefix := "\033["
	codeSuffix := "m"
	for strings.Contains(str, codePrefix) {
		prefixIndex := strings.Index(str, codePrefix)
		sub := str[prefixIndex:]
		code := str[prefixIndex : prefixIndex+strings.Index(sub, codeSuffix)+1]
		str = strings.ReplaceAll(str, code, "")
	}
	return str
}

func getFileContents(filePath string) string {
	b, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read file: " + filePath)
		fmt.Println(err)
	}
	return string(b)
}

func saveCommandSignature(CASE string, LABEL string, command string, output string, numOflines int) {
	fileContent := ""
	filePath := filepath.Join(CASE, ".trustchain")
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		fileContent = getFileContents(filePath)
	}
	data := fileContent + output
	line := LABEL + "\t" + strconv.Itoa(numOflines) + "\t" + command + "\t"
	line += generateSignature(data + line)
	saveLineToFile(filePath, line)
}

func generateSignature(str string) string {
	timestamp := time.Now().Format("D20060102T150405")
	hash := sha256.New()
	hash.Write([]byte(str + timestamp))
	return timestamp + "H" + hex.EncodeToString(hash.Sum(nil))
}
