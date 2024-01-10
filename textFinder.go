package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {

	log.Println("Text finder started.")

	logfile, err := os.OpenFile("app.log", os.O_APPEND, 0)

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logfile)
	defer logfile.Close()

	var wg sync.WaitGroup

	countLines := flag.Bool("c", false, "This is used to display the count of repeated characters lines")
	caseSensitive := flag.Bool("i", false, "Case sensitive search")
	invertMatching := flag.Bool("v", false, "This prints out all the lines that do not matches the pattern")
	displayLineNumbers := flag.Bool("n", false, "Display the matched lines and their line numbers")

	log.Println("Waiting for a flag to proceed")

	flag.Parse()

	searchString := flag.Arg(0)

	fileNames := flag.Args()[1:]

	if len(fileNames) == 0 || searchString == "" {

		log.Fatal("Usage: ./textFinder flag[c, i, v, n] <search string> <filename/s>")
		log.Println("Error reading from the terminal.")
	}

	errChannel := make(chan error)

	for _, file := range fileNames {

		wg.Add(1)

		go func(file string) {

			defer wg.Done()

			err := searchingFunc(file, searchString, *countLines, *displayLineNumbers, *caseSensitive, *invertMatching)

			if err != nil {

				errChannel <- fmt.Errorf("error in file %s: %v", file, err)
			}
		}(file)
	}

	// Close the error channel after all the go routines are finished

	go func() {

		wg.Wait()
		close(errChannel)
	}()

	for err := range errChannel {

		log.Println(err)
	}
}

func searchingFunc(fileName, searchString string, countLines, displayLineNumbers, caseSensitive, invertMatching bool) error {

	file, err := os.Open(fileName)

	if err != nil {

		return err
	}

	log.Println("Searching function call.")

	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 0

	log.Println("Scanning the text")

	for scanner.Scan() {

		line := scanner.Text()

		if countLines {
			if strings.Contains(line, searchString) {

				// log.Println("Counting lines")

				lineNumber++
				continue
			}
		}

		if displayLineNumbers {

			if strings.Contains(line, searchString) {

				fmt.Println(line)
				// log.Println("Displaying lines")
				lineNumber++
				continue
			}
		}

		if caseSensitive {

			// log.Println("To convert lower string")
			lowerString := strings.ToLower(searchString)

			// log.Println("To convert upper string")
			upperString := strings.ToUpper(searchString)

			camleCaseString := strings.Title(lowerString)

			if strings.Contains(line, lowerString) || strings.Contains(line, upperString) || strings.Contains(line, camleCaseString) {

				fmt.Println(line)
				continue
			}
		}

		if invertMatching {

			matched := strings.Contains(line, searchString)

			// log.Println("invert match log")

			if matched == false {
				fmt.Println(line)
				continue
			}
			continue
		}

		if strings.Contains(line, searchString) {

			fmt.Println(line)
		}

	}

	if countLines || displayLineNumbers {
		fmt.Println(lineNumber)
	}

	if err := scanner.Err(); err != nil {

		return err
	}

	log.Println("Output is on the console.")

	return nil
}
