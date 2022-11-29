package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

//statsFunc defines a generic statistical function to represent sum and avg signature
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

//csv2float read data from r [csv file], with column specified
// to return slice of float64 of that column
func csv2float(r io.Reader, column int) ([]float64, error) {
	//create the csv reader, used to read in data from csv files
	cr := csv.NewReader(r)

	//Adjusting for 0 based index
	column--

	//read in all CSV data
	records, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}
	var data []float64

	//loop through all records
	for i, row := range records {
		//skip the file header
		if i == 0 {
			continue
		}
		// checking number of col in csv file
		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		//try to convert data read into a float number
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}
	return data, nil
}
