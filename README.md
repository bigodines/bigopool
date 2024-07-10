# bigopool [![codecov](https://codecov.io/gh/bigodines/bigopool/branch/master/graph/badge.svg)](https://codecov.io/gh/bigodines/bigopool) [![Build Status](https://travis-ci.org/bigodines/bigopool.png)](https://travis-ci.org/bigodines/bigopool) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/bigodines/bigopool)


`bigopool` is a small library that implements high performance worker pool in Golang and allows `error`/`result` handling in the main thread.

It also provides a thread safe parallel processing abstraction

## Quickstart

install:
`go get -u github.com/bigodines/bigopool`

### Worker pool

Use this whenever the number of jobs you have is too large to run as goroutines at the same time.

implement this simple interface:
```golang
type TestJob struct {
    // your properties go here
}
func (j TestJob) Execute() (bigopool.Result, error) {
    // your logic here.
    // Result is an interface{}
    return "anything", nil
}
```

add to your code:
```golang
// configure dispatcher to run 5 workers with a queue of capacity 100
dispatcher, err := bigopool.NewDispatcher(5, 100)
if err != nil {
    panic(err)
}
// spawn workers
dispatcher.Run()
// send work items
dispatcher.Enqueue(TestJob{}) // <-- add one job
dispatcher.Enqueue(TestJob{}, TestJob{}) // <-- add multiple jobs

// wait for workers to finish (this is a blocking call)
results, errs := dispatcher.Wait()
// Note we've opted to use our own thread safe error module (it still implements the `error` interface but if you want to get a []error,
//  use errs.All())
```

:boom:

### Parallel processing

Use this if you don't need worker pools and just want to execute tasks in parallel

run multiple functions in parallel:
```go
func UploadAndDownload() error {
  var email string
  errs := bigopool.Parallel(
    func() error {
      return api.Post(Request{Order: 123123})
    },

    func() error {
      user, err := api.Get(Request{ID: 1234})
      if err != nil {
        return err
      }

      email = user.Email
      return nil
    },
  )

  fmt.Println("email:", email)
  return errs.ToError()
}
```

## Contributing

If you can fix a bug or make it faster, I'll buy you coffee. PRs that drop code coverage will not be merged.

## Benchmarks
```
box specs:

MacBook Pro (Retina, 15-inch, 2016)
Processor 2,7 GHz QUad-Core Intel Core i7
Memory 16 GB 2133 MHz LPDDR3
```
```bash
$ go test -bench=.  -benchmem=true -cpu=1,2,4,8 -run Bench
goos: darwin
goarch: amd64
pkg: github.com/bigodines/bigopool
Benchmark1Workers1Queue                  1756252               687 ns/op              98 B/op          0 allocs/op
Benchmark1Workers1Queue-2                1209124              1041 ns/op              91 B/op          0 allocs/op
Benchmark1Workers1Queue-4                1242471               953 ns/op              89 B/op          0 allocs/op
Benchmark1Workers1Queue-8                1151955              1023 ns/op              96 B/op          0 allocs/op
Benchmark5Workers1000Queue               2190982               566 ns/op              98 B/op          0 allocs/op
Benchmark5Workers1000Queue-2             1960696               610 ns/op              88 B/op          0 allocs/op
Benchmark5Workers1000Queue-4             1727991               692 ns/op              80 B/op          0 allocs/op
Benchmark5Workers1000Queue-8             1508758               773 ns/op              91 B/op          0 allocs/op
Benchmark10Workers100Queue               2168791               559 ns/op              79 B/op          0 allocs/op
Benchmark10Workers100Queue-2             1961818               624 ns/op              88 B/op          0 allocs/op
Benchmark10Workers100Queue-4             1672346               702 ns/op              82 B/op          0 allocs/op
Benchmark10Workers100Queue-8             1382486               864 ns/op              80 B/op          0 allocs/op
Benchmark20Workers200Queue               2159338               572 ns/op              80 B/op          0 allocs/op
Benchmark20Workers200Queue-2             1967389               608 ns/op              88 B/op          0 allocs/op
Benchmark20Workers200Queue-4             1765658               711 ns/op              98 B/op          0 allocs/op
Benchmark20Workers200Queue-8             1445559               838 ns/op              95 B/op          0 allocs/op
Benchmark20Workers10000Queue             2220480               551 ns/op              97 B/op          0 allocs/op
Benchmark20Workers10000Queue-2           1990562               629 ns/op              87 B/op          0 allocs/op
Benchmark20Workers10000Queue-4           1721840               684 ns/op              80 B/op          0 allocs/op
Benchmark20Workers10000Queue-8           1417897               835 ns/op              97 B/op          0 allocs/op
Benchmark100Workers10000Queue            2206069               556 ns/op              98 B/op          0 allocs/op
Benchmark100Workers10000Queue-2          2074681               586 ns/op              83 B/op          0 allocs/op
Benchmark100Workers10000Queue-4          1824800               685 ns/op              95 B/op          0 allocs/op
Benchmark100Workers10000Queue-8          1537908               785 ns/op              90 B/op          0 allocs/op
PASS
ok      github.com/bigodines/bigopool   46.191s
```

## License

MIT
