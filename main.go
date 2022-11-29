package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	//verify and parse args
	op := flag.String("op", "sum", "Operation to be executed")
	col := flag.Int("col", 1, "CSV column on which to execute operation")
	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filenames []string, op string, col int, dest io.Writer) error {
	var opFunc statsFunc

	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if col < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidColumn, col)
	}

	cons, err := consolidate(filenames, col)
	if err != nil {
		return fmt.Errorf("%w can not consolidate data", err)
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	_, err = fmt.Fprintln(dest, opFunc(cons))
	return err

}

func consolidate(fileNames []string, col int) ([]float64, error) {
	cons := make([]float64, 0)
	// Loop through all files adding their data to consolidate
	for _, fn := range fileNames {
		// Open the file for reading
		f, err := os.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("cannot open file: %w", err)
		}

		// Parse the CSV into a slice of float64 numbers
		data, err := csv2float(f, col)
		if err != nil {
			return nil, err
		}

		if err := f.Close(); err != nil {
			return nil, err
		}

		// Append the data to consolidate
		cons = append(cons, data...)
	}
	return cons, nil
}
