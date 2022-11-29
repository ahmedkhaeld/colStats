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
	//enable opt to reuse the same slice
	cr.ReuseRecord = true

	//Adjusting for 0 based index
	column--

	var data []float64

	//read CSV data
	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}

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
