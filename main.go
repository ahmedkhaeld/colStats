package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
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

	//validate the operation and define the opFunc accordingly
	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return fmt.Errorf("%w: %s", ErrInvalidOperation, op)
	}

	err := consolidate(filenames, col, opFunc, dest)
	if err != nil {
		return fmt.Errorf("%w can not consolidate data", err)
	}

	return err

}

func consolidate(filenames []string, column int, opFunc statsFunc, out io.Writer) error {
	collected := make([]float64, 0)

	// Create the channel to receive results or errors of operations
	resCh := make(chan []float64)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	wg := sync.WaitGroup{}

	// Loop through all files and create a goroutine to process
	// each one concurrently
	for _, fn := range filenames {
		wg.Add(1)
		go func(fn string) {

			defer wg.Done()

			// Open the file for reading
			f, err := os.Open(fn)
			if err != nil {
				errCh <- fmt.Errorf("cannot open file: %w", err)
				return
			}

			// Parse the CSV into a slice of float64 numbers
			data, err := csv2float(f, column)
			if err != nil {
				errCh <- err
			}

			if err := f.Close(); err != nil {
				errCh <- err
			}

			resCh <- data
		}(fn)
	}

	//concurrently, close all file have been processed
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case data := <-resCh:
			collected = append(collected, data...)
		case <-doneCh:
			_, err := fmt.Fprintln(out, opFunc(collected))
			return err
		}
	}
}
