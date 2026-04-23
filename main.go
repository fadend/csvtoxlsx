package main

import (
	xlsx "codeberg.org/tealeg/xlsx/v4"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	inputCSV   = flag.String("input_csv", "", "Path to input CSV file.")
	sheetName  = flag.String("sheet_name", "Sheet 1", "Name for the sheet created with the data.")
	outputXLSX = flag.String("output_xlsx", "", "Path to the output .xlsx file.")
)

func checkNonemptyFlag(flagValue *string, name string) {
	if *flagValue == "" {
		fmt.Fprintf(os.Stderr, "Non-empty value required for flag: %s\n", name)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
	checkNonemptyFlag(inputCSV, "--input_csv")
	checkNonemptyFlag(sheetName, "--sheet_name")
	checkNonemptyFlag(outputXLSX, "--output_xlsx")
	csvFile, err := os.Open(*inputCSV)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem opening CSV %s: %v\n", *inputCSV, err)
		os.Exit(1)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)
	xlsxFile := xlsx.NewFile()
	sheet, err := xlsxFile.AddSheet(*sheetName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem creating sheet <<%s>>: %v\n", *sheetName, err)
		os.Exit(1)
	}
	for {
		csvRow, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading row from %s: %v\n", *inputCSV, err)
			os.Exit(1)
		}
		row := sheet.AddRow()
		for _, value := range csvRow {
			row.AddCell().SetValue(value)
		}
	}
	if err = xlsxFile.Save(*outputXLSX); err != nil {
		fmt.Fprintf(os.Stderr, "Problem saving .xlsx file %s: %v\n", *outputXLSX, err)
		os.Exit(1)
	}
}
