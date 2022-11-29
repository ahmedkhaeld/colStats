# colStats
column stats an application that executes
statistical operations on a CSV file

---
#### Benchmarking the tool
`$ go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults03m.txt`
``` 
goos: linux
goarch: amd64
pkg: github.com/ahmedkhaeld/colStats
cpu: Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz
BenchmarkRun-4                10         355138870 ns/op        230413111 B/op      2528052 allocs/op
PASS
ok      github.com/ahmedkhaeld/colStats 3.984s


```
`benchcmp benchresults02m.txt benchresults03m.txt`
``` 
benchmark          old ns/op     new ns/op     delta
BenchmarkRun-4     359666931     351380942     -2.30%

benchmark          old allocs     new allocs     delta
BenchmarkRun-4     2530554        2528051        -0.10%

benchmark          old bytes     new bytes     delta
BenchmarkRun-4     230547308     230413101     -0.06%

```
as you can see this version runs over 2.3% faster than previous

Compare this result to the original version
`benchcmp benchresults00m.txt benchresults03m.txt`

``` 
benchmark          old ns/op     new ns/op     delta
BenchmarkRun-4     882514991     351380942     -60.18%

benchmark          old allocs     new allocs     delta
BenchmarkRun-4     5043040        2528051        -49.87%

benchmark          old bytes     new bytes     delta
BenchmarkRun-4     564339412     230413101     -59.17%

```
see how much we've improved the performance
this version is over almost three times faster than original.
it allocates nearly 60% less memory

`$go build`<br>
`$ time ./colStats -op avg -col 2 testdata/benchmark/*.csv`
```
50006.0653788

real    0m0.398s
user    0m1.241s
sys     0m0.059s

```
this time, the program processed all one thousand files in 0.38 seconds
compared to the original one second
