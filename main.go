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
	kmode := false
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
				fmt.Println("")
				fmt.Println("Run Nested Shell: keeptrak")
				fmt.Println("")
				fmt.Println("    Run Nested Shell in K Mode: keeptrak -k")
				fmt.Println("    K Mode reduces the terminal prefix from KEEPTRAK> to K>")
				fmt.Println("")
				fmt.Println("    WITHIN NESTED MODE")
				fmt.Println("")
				fmt.Println("    You can run any keeptrak commands normally or use nested commands.")
				fmt.Println("    Nested commands are triggered by using a colon followed by a space after the command.")
				fmt.Println("    Example: you can store a note by running KEEPTRAK> note: this is my note")
				fmt.Println("    Instead of KEEPTRAK> keeptrak note \"this is my note\"")
				fmt.Println("    - note: TEXT")
				fmt.Println("        Any and all text following the nested note command will be saved as the note.")
				fmt.Println("    - record: LABEL VALUE DATA_TYPE CONFIRMED")
				fmt.Println("        The nested record command requires at least the first two parameters.")
				fmt.Println("        DATA_TYPE and CONFIRMED are optional with the default values of null and N.")
				fmt.Println("    - ghost: BASH")
				fmt.Println("        The ghost command will run whatever follows without recording it anywhere.")
				fmt.Println("        Example: ghost: ls")
				fmt.Println("        For short hand you can use g: ")
				fmt.Println("")
				fmt.Println("Save Record: keeptrak CASE LABEL VALUE DATA_TYPE CONFIRMED")
				fmt.Println("\tExample: keeptrak johndoe username jdoe credential Y")
				fmt.Println("")
				fmt.Println("Save Note: keeptrak note TEXT")
				fmt.Println("\tExample: keeptrak note \"This is useful information\"")
				fmt.Println("")
				return
			} else if os.Args[1] == "-k" {
				kmode = true
			} else {
				fmt.Println("Unknown command: " + os.Args[1])
				return
			}
		}
		if len(os.Args) > 2 {
			CASE := os.Args[1]
			ensureCaseDir(CASE)
			if len(os.Args) == 4 {
				if os.Args[2] == "note" {
					saveLineToFile(filepath.Join(CASE, "notes"), getTime()+"\t"+os.Args[3])
				}
			} else if len(os.Args) < 6 {
				fmt.Println("Too few arguments. Run --help to see correct usage")
			} else if len(os.Args) == 6 {
				// WITH ARGUMENTS
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
				if kmode {
					fmt.Print("K> ")
				} else {
					fmt.Print("KEEPTRAK> ")
				}
				text, _ := reader.ReadString('\n')
				// convert CRLF to LF
				command := strings.Replace(text, "\n", "", -1)
				if command != "" {
					if command == "exit" {
						return
					}
					if strings.HasPrefix(command, "note: ") {
						saveLineToFile(filepath.Join(CASE, "notes"), getTime()+"\t"+command[6:])
						continue
					}
					if strings.HasPrefix(command, "record: ") {
						Args := strings.Split(command[8:], " ")
						if len(Args) < 2 {
							fmt.Println("Record command requires at least two arguements (record: LABEL VALUE)")
							continue
						}
						LABEL := Args[0]
						VALUE := Args[1]
						DATATYPE := "null"
						CONFIRMED := "N"
						if len(Args) > 2 {
							DATATYPE = Args[2]
						}
						if len(Args) > 3 {
							CONFIRMED = Args[3]
						}
						dbpath := ensureDB(CASE)
						saveRecord(dbpath, LABEL, VALUE, DATATYPE, CONFIRMED)
						continue
					}
					ghostMode := false
					if strings.HasPrefix(command, "ghost: ") {
						ghostMode = true
						command = command[7:]
					}
					if strings.HasPrefix(command, "g: ") {
						ghostMode = true
						command = command[3:]
					}
					if !ghostMode {
						saveLineToFile(CASE+"/history", getTime()+"\t"+command)
					}
					cmd := exec.Command("bash", "-c", command)
					pipe, _ := cmd.StdoutPipe()
					if err := cmd.Start(); err != nil && err != io.EOF {
						fmt.Println(err)
					}
					reader := bufio.NewReader(pipe)
					line, err := reader.ReadString('\n')
					LABEL := strings.Split(command, " ")[0]
					for err == nil {
						line = strings.ReplaceAll(line, "\n", "")
						fmt.Println(line)
						if !ghostMode {
							saveLineToFile(filepath.Join(CASE, LABEL), line)
							saveLineToFile(filepath.Join(CASE, "dump"), line)
						}
						line, err = reader.ReadString('\n')
					}
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
