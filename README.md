# colStats
column stats an application that executes
statistical operations on a CSV file

---
#### Benchmarking the tool
`$ go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults02m.txt`
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         359666931 ns/op        230547308 B/op      2530554 allocs/op
PASS
ok      github.com/ahmedkhaeld/colStats 4.099s

```
`benchcmp benchresults01m.txt benchresults02m.txt`
``` 
benchmark          old ns/op     new ns/op     delta
BenchmarkRun-4     576020868     359666931     -37.56%

benchmark          old allocs     new allocs     delta
BenchmarkRun-4     2528036        2530554        +0.10%

benchmark          old bytes     new bytes     delta
BenchmarkRun-4     230411072     230547308     +0.06%
```
Another great improvement, it's almost twice as fast the previous version



---
### Profiling the tool
`$ go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu01.pprof`
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         580915778 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 6.645s
```

`$ go tool pprof cpu01.pprof`
``` 
File: colStats.test
Type: cpu
Time: Nov 29, 2022 at 3:54pm (EET)
Duration: 6.64s, Total samples = 7.12s (107.27%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top -cum
Showing nodes accounting for 2400ms, 33.71% of 7120ms total
Dropped 121 nodes (cum <= 35.60ms)
Showing top 10 nodes out of 79
      flat  flat%   sum%        cum   cum%
         0     0%     0%     6420ms 90.17%  github.com/ahmedkhaeld/colStats.BenchmarkRun
         0     0%     0%     6420ms 90.17%  github.com/ahmedkhaeld/colStats.run
         0     0%     0%     6420ms 90.17%  testing.(*B).runN
         0     0%     0%     6400ms 89.89%  github.com/ahmedkhaeld/colStats.consolidate
     230ms  3.23%  3.23%     5940ms 83.43%  github.com/ahmedkhaeld/colStats.csv2float
         0     0%  3.23%     5820ms 81.74%  testing.(*B).launch
     110ms  1.54%  4.78%     4220ms 59.27%  encoding/csv.(*Reader).Read
    1270ms 17.84% 22.61%     4110ms 57.72%  encoding/csv.(*Reader).readRecord
     190ms  2.67% 25.28%     1260ms 17.70%  runtime.slicebytetostring
     600ms  8.43% 33.71%     1210ms 16.99%  runtime.mallocgc
(pprof) 

```
The profiling has changed slightly. The top part is still the same, as the same functions are responsible for executing the program. The csv2float function is still there, which also makes sense. But in the bottom part of the output, the functions related to memory allocation and garbage collection are no longer in the top 10. 

---
### Tracing the tool
`$ go test -bench . -benchtime=10x -run ^$ -trace trace02.out`
```
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         408675980 ns/op
PASS
ok      github.com/ahmedkhaeld/colStats 4.584s
```
`go tool trace trace02.out`

view the trace link
>the program used all four CPUs, improving the speed of the tool. 