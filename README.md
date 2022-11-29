# colStats
column stats an application that executes
statistical operations on a CSV file

1. `-col`: flag on which col to execute the operation. It defaults to 1.
2. `-op`: flag of The operation to execute on the selected column. 

`$ go build`
* execute on single file
<br> `$  ./colStats -op avg -col 3 testdata/example.csv`

* execute on multiple files
<br> `$ ./colStats -op avg -col 3 testdata/example.csv testdata/example2.csv`

---
### Current status
* Current status: How fast it takes to precess the workload<br>
`$ time ./colStats -op avg -col 3 testdata/benchmark/*.csv`

```` 
real    0m0.995s
user    0m1.037s
sys     0m0.124s
````
As you can see it takes nearly 1 second, still small amount of time

---
#### Benchmarking the tool

``$ go test -bench . -run ^$`` <br>
using go test tool with `-bench regexp` parameter
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                 2         902543468 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 2.796s

```
The benchmark executed two times, this because it processed a 
thousand files, and it took more than 2 seconds to complete
. The benchmark tool also prints the avg time per operation in nano sec 

force benchmark execution for 10 times, to get rid of potential noise<br>
`$ go test -bench . -benchtime=10x -run ^$`

``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         888847335 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 9.864s

```
---
### Profiling the tool
shows a breakdown of where the program spends its execution time
<br> enable the CPU profiler:<br>
`$ go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu00.pprof`
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         884717278 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 10.064s

```

Analyze the profiling results by <br>
`$ go tool pprof cpu00.pprof`
``` 
File: colStats.test
Type: cpu
Time: Nov 29, 2022 at 11:52am (EET)
Duration: 10.05s, Total samples = 10.99s (109.39%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof)
```
when the profiler is enabled, it stops the program execution every 10 ms
and takes a sample of the function stack [This sample contains all functions that are executing or waiting to execute at that time. T]

use `top` to see where the program is spending most of its time
``` 
(pprof) top
Showing nodes accounting for 7180ms, 65.33% of 10990ms total
Dropped 149 nodes (cum <= 54.95ms)
Showing top 10 nodes out of 88
      flat  flat%   sum%        cum   cum%
    1440ms 13.10% 13.10%     1530ms 13.92%  runtime.heapBitsSetType
    1100ms 10.01% 23.11%     6000ms 54.60%  encoding/csv.(*Reader).readRecord
     900ms  8.19% 31.30%      900ms  8.19%  runtime.memclrNoHeapPointers
     880ms  8.01% 39.31%     3940ms 35.85%  runtime.mallocgc
     790ms  7.19% 46.50%      790ms  7.19%  strconv.readFloat
     760ms  6.92% 53.41%      760ms  6.92%  indexbytebody
     540ms  4.91% 58.33%      540ms  4.91%  runtime.memmove
     280ms  2.55% 60.87%      280ms  2.55%  runtime.procyield
     270ms  2.46% 63.33%     1350ms 12.28%  strconv.atof64
     220ms  2.00% 65.33%      220ms  2.00%  internal/bytealg.IndexByte
(pprof) 
```
flat time: the time the function spends executing on the CPU.
Here the program spends 13.92% of its CPU time executing runtime.heapBitsSetType

sorting based on the cumulative time
``` 
(pprof) top -cum
Showing nodes accounting for 2.37s, 21.57% of 10.99s total
Dropped 149 nodes (cum <= 0.05s)
Showing top 10 nodes out of 88
      flat  flat%   sum%        cum   cum%
         0     0%     0%      9.43s 85.81%  github.com/ahmedkhaeld/colStats.BenchmarkRun
         0     0%     0%      9.43s 85.81%  github.com/ahmedkhaeld/colStats.run
         0     0%     0%      9.43s 85.81%  testing.(*B).runN
         0     0%     0%      9.42s 85.71%  github.com/ahmedkhaeld/colStats.consolidate
     0.18s  1.64%  1.64%      8.99s 81.80%  github.com/ahmedkhaeld/colStats.csv2float
         0     0%  1.64%      8.57s 77.98%  testing.(*B).launch
     0.19s  1.73%  3.37%      7.16s 65.15%  encoding/csv.(*Reader).ReadAll
     1.10s 10.01% 13.38%         6s 54.60%  encoding/csv.(*Reader).readRecord
     0.88s  8.01% 21.38%      3.94s 35.85%  runtime.mallocgc
     0.02s  0.18% 21.57%      2.31s 21.02%  runtime.makeslice
(pprof) 
```
The cumulative time accounts for the time the function was 
executing or waiting for a called function to return.

Here, the program spends most of its time on functions related to the benchmark functionality, which are irrelevant to us. 
You can see from this output that your program is spending over 85% of its time on the consolidate function. 
This is important since this is a function that you wrote.

a deeper look at how this function is spending its time by using the list subcommand. This subcommand displays the source code of the function annotated with the time spent to run each line of code.
```  
(pprof) list consolidate
Total: 10.99s
ROUTINE ======================== github.com/ahmedkhaeld/colStats.consolidate in /home/ahmed/colStats/main.go
         0      9.42s (flat, cum) 85.71% of Total
         .          .     52:func consolidate(fileNames []string, col int) ([]float64, error) {
         .          .     53:   cons := make([]float64, 0)
         .          .     54:   // Loop through all files adding their data to consolidate
         .          .     55:   for _, fn := range fileNames {
         .          .     56:           // Open the file for reading
         .       40ms     57:           f, err := os.Open(fn)
         .          .     58:           if err != nil {
         .          .     59:                   return nil, fmt.Errorf("cannot open file: %w", err)
         .          .     60:           }
         .          .     61:
         .          .     62:           // Parse the CSV into a slice of float64 numbers
         .      8.99s     63:           data, err := csv2float(f, col)
         .          .     64:           if err != nil {
         .          .     65:                   return nil, err
         .          .     66:           }
         .          .     67:
         .       50ms     68:           if err := f.Close(); err != nil {
         .          .     69:                   return nil, err
         .          .     70:           }
         .          .     71:
         .          .     72:           // Append the data to consolidate
         .      340ms     73:           cons = append(cons, data...)
         .          .     74:   }
         .          .     75:   return cons, nil
         .          .     76:}
(pprof) 

```
`8.99s     63:           data, err := csv2float(f, col)`
This output shows that the program is spending  8.99s to complete
the `csv2float` function

```
(pprof) list csv2float
Total: 10.99s
ROUTINE ======================== github.com/ahmedkhaeld/colStats.csv2float in /home/ahmed/colStats/csv.go
     180ms      8.99s (flat, cum) 81.80% of Total
         .          .     24:
         .          .     25://csv2float read data from r [csv file], with column specified
         .          .     26:// to return slice of float64 of that column
         .          .     27:func csv2float(r io.Reader, column int) ([]float64, error) {
         .          .     28:   //create the csv reader, used to read in data from csv files
         .       30ms     29:   cr := csv.NewReader(r)
         .          .     30:
         .          .     31:   //Adjusting for 0 based index
         .          .     32:   column--
         .          .     33:
         .          .     34:   //read in all CSV data
         .      7.16s     35:   records, err := cr.ReadAll()
      10ms       10ms     36:   if err != nil {
         .          .     37:           return nil, fmt.Errorf("cannot read data from file: %w", err)
         .          .     38:   }
         .          .     39:   var data []float64
         .          .     40:
         .          .     41:   //loop through all records
      90ms       90ms     42:   for i, row := range records {
         .          .     43:           //skip the file header
      10ms       10ms     44:           if i == 0 {
         .          .     45:                   continue
         .          .     46:           }
         .          .     47:           // checking number of col in csv file
         .          .     48:           if len(row) <= column {
         .          .     49:                   return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
         .          .     50:           }
         .          .     51:
         .          .     52:           //try to convert data read into a float number
      40ms      1.41s     53:           v, err := strconv.ParseFloat(row[column], 64)
         .          .     54:           if err != nil {
         .          .     55:                   return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
         .          .     56:           } 7.16s     35:   records, err := cr.ReadAll()
      30ms      280ms     57:           data = append(data, v)
         .          .     58:   }
         .          .     59:   return data, nil
         .          .     60:}
(pprof) 

```
` 7.16s     35:   records, err := cr.ReadAll()`

---
### How much memory the program is allocating
to create a memory profile use:<br>
`$ go test -bench . -benchtime=10x -run ^$ -memprofile mem00.pprof`
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         897561793 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 9.936s
```
to view the results use<br>
`$ go tool pprof -alloc_space mem00.pprof`
``` 
File: colStats.test
Type: alloc_space
Time: Nov 29, 2022 at 2:04pm (EET)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top -cum
Showing nodes accounting for 5.78GB, 99.10% of 5.84GB total
Dropped 28 nodes (cum <= 0.03GB)
Showing top 10 nodes out of 12
      flat  flat%   sum%        cum   cum%
         0     0%     0%     5.83GB   100%  github.com/ahmedkhaeld/colStats.BenchmarkRun
         0     0%     0%     5.83GB   100%  testing.(*B).runN
    1.05GB 18.00% 18.00%     5.83GB 99.94%  github.com/ahmedkhaeld/colStats.consolidate
         0     0% 18.00%     5.83GB 99.94%  github.com/ahmedkhaeld/colStats.run
         0     0% 18.00%     5.26GB 90.23%  testing.(*B).launch
    0.87GB 14.91% 32.91%     4.78GB 81.90%  github.com/ahmedkhaeld/colStats.csv2float
    2.60GB 44.47% 77.39%     3.86GB 66.19%  encoding/csv.(*Reader).ReadAll
    1.27GB 21.71% 99.10%     1.27GB 21.71%  encoding/csv.(*Reader).readRecord
         0     0% 99.10%     0.57GB  9.74%  testing.(*B).run1.func1
         0     0% 99.10%     0.05GB  0.81%  bufio.NewReader (inline)
(pprof) As we suspected, the ReadAll function from the encoding/csv package is responsible for the allocation of almost 4GB of memory, which corresponds to 67% of all the memory allocation for this program. The more memory we allocate, the more garbage collection has to run, increasing the time it takes to run the program.

```
As we suspected, the ReadAll function from the encoding/csv package is responsible for the allocation of almost 4GB of memory, which corresponds to 67% of all the memory allocation for this program. The more memory we allocate, the more garbage collection has to run, increasing the time it takes to run the program.
`    2.60GB 44.47% 77.39%     3.86GB 66.19%  encoding/csv.(*Reader).ReadAll`