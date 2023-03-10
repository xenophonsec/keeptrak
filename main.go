package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
			if os.Args[1] == "--help" {
				fmt.Println("USAGE: keeptrak CASE LABEL VALUE DATA_TYPE CONFIRMED")
			}
		} else if len(os.Args) > 2 {
			if len(os.Args) < 6 {
				fmt.Println("Too few arguments. Run --help to see correct usage")
			} else {
				// WITH ARGUMENTS
				CASE := os.Args[1]
				LABEL := os.Args[2]
				VALUE := os.Args[3]
				DATATYPE := os.Args[4]
				CONFIRMED := os.Args[5]
				ensureCaseDir(CASE)
				dbpath := ensureDB(CASE)
				saveRecord(dbpath, LABEL, VALUE, DATATYPE, CONFIRMED)
			}
		} else {
			// WRAPPED SHELL
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("")
			fmt.Println("K E E P T R A K \t" + getTime())
			fmt.Println("")
			fmt.Print("Enter Case Name: ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			CASE := strings.Replace(text, "\n", "", -1)
			ensureCaseDir(CASE)
			dbpath := ensureDB(CASE)

			for {
				fmt.Print("KEEPTRAK> ")
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				command := strings.Replace(text, "\n", "", -1)
				if command != "" {
					saveLineToFile(CASE+"/history", command)
					out, err := exec.Command("bash", "-c", command).Output()
					// TODO: https://pkg.go.dev/os/exec#Cmd.StdoutPipe
					if err != nil {
						fmt.Println(err)
					} else {
						LABEL := strings.Split(command, " ")[0]
						text := string(out)
						fmt.Println(text)
						saveRecord(dbpath, LABEL, "", "FILE", "N")
						saveLineToFile(CASE+"/"+LABEL, text)
						saveLineToFile(CASE+"/dump", text)

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
