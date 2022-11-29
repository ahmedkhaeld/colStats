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
