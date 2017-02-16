package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"encoding/csv"
	"io"
	"bufio"
	"strconv"
)

// ##### Constants ###########################################################

const APP_NAME string = "csv-value-counter"
const APP_VERSION string = "1.0.1"

// ##### Variables ###########################################################

var (
	data map[string]int

	inputFilePath  = kingpin.Flag("input", "Input file containing the data").Short('i').Required().String()
	outputFilePath = kingpin.Flag("output", "Output file path for results").Short('o').Required().String()
	fieldIndex     = kingpin.Flag("field", "Field index (Defaults to 1)").Short('f').Required().Int()
	sorted		   = kingpin.Flag("sort", "Sort (Defaults to false").Short('s').Bool()
	delimiter	   = kingpin.Flag("delimiter", "Delimiter (Defaults to ,").Short('d').String()
)

// ##### Methods #############################################################

// Application entry point
func main() {

	fmt.Printf("\n%s %s\n\n", APP_NAME, APP_VERSION)

	kingpin.Parse()

	if len(*delimiter) > 1 {
		fmt.Println("Invalid delimiter, need to be a single character: " + *delimiter)
		return
	}

	inFile, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("Error opening input file: " + err.Error())
		return
	}
	defer inFile.Close()

	outFile, err := os.Create(*outputFilePath)
	if err != nil {
		fmt.Println("Error creating output file: " + err.Error())
		return
	}
	defer outFile.Close()

	comma := ';'
	if *delimiter != "" {
		comma = []rune(*delimiter)[0]
	}

	reader := csv.NewReader(inFile)
	reader.Comma = comma

	// Decrement the field index as the CSV reader uses a 0 based index
	if *fieldIndex != 0 {
		*fieldIndex--
	}

	data = make(map[string]int, 0)
	var record []string
	var count int
	for {
		// read just one record, but we could ReadAll() as well
		record, err = reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if *fieldIndex > (len(record) - 1)  {
			continue
		}

		count = data[record[*fieldIndex]]
		count++
		data[record[*fieldIndex]] = count
	}

	writer := bufio.NewWriter(outFile)

	if *sorted == true {
		pl := RankByWordCount(data)

		for _, v := range pl {
			_, _ = writer.WriteString(v.Key + "," + strconv.Itoa(v.Value) + "\n")
		}
	} else {
		for v := range data {
			_, _ = writer.WriteString(v + "," + strconv.Itoa(data[v]) + "\n")
		}
	}

	writer.Flush()
}


