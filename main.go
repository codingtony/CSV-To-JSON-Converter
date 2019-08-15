package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	log.Println(strings.Repeat("=", 10), "Starting", strings.Repeat("=", 10))
	path := flag.String("path", "./data.csv", "Path of the file")
	flag.Parse()
	log.Printf("Opening %s\n", *path)
	csvFile, err := os.Open(*path)
	if err != nil {
		log.Fatal("The file is not found || wrong root")
	}
	defer csvFile.Close()

	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	log.Printf("Writing to %s\n", newFileName)
	jsonFile, err := os.Create(newFileName)
	if err != nil {
		log.Fatal("Unable to create json file")
	}
	defer jsonFile.Close()
	lineCount := ConvertCSVToJson(csvFile, jsonFile)
	log.Printf("Processed %d rows\n", lineCount)
	log.Println(strings.Repeat("=", 12), "Done", strings.Repeat("=", 12))
}

// ConvertCSVToJson to convert the content of CSV File to a JSON File
func ConvertCSVToJson(csvFile *os.File, jsonFile *os.File) int64 {

	reader := csv.NewReader(csvFile)
	var lineNumber int64
	headersArr := make([]string, 0)
	var buffer bytes.Buffer
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if lineNumber == 0 {
			for _, headE := range line {
				headersArr = append(headersArr, headE)
			}
		} else {
			buffer.WriteString("{")
			for j, y := range line {
				buffer.WriteString(`"` + headersArr[j] + `":`)
				_, fErr := strconv.ParseFloat(y, 32)
				_, bErr := strconv.ParseBool(y)
				if fErr == nil {
					buffer.WriteString(y)
				} else if bErr == nil {
					buffer.WriteString(strings.ToLower(y))
				} else {
					buffer.WriteString((`"` + y + `"`))
				}
				//end of property
				if j < len(line)-1 {
					buffer.WriteString(",")
				}

			}
			//end of object of the array
			buffer.WriteString("}")
			buffer.WriteString("\n")
			jsonFile.WriteString(buffer.String())
			buffer.Reset()
		}
		lineNumber++

	}
	return lineNumber

}
