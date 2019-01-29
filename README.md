# bigopool

`bigopool` is a small library that implements high performance worker pool in Golang and allows `error`/`result` handling in the main thread.

It also provides a thread safe parallel processing abstraction

## Quickstart

install:
`go get -u github.com/bigodines/bigopool`

### Worker pool

Use this whenever the number of jobs you have is too large to run as goroutines at the same time.

implement this simple interface:
```golang
type TestJob {
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

## Motivation

This is my take on the approach outlined on this classic blog post: http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/

I've added support to errors and results

## Contributing

If you can fix a bug or make it faster, I'll buy you coffee. PRs that drop code coverage will not be merged.

## Benchmarks
```
box specs:

MacBook Pro (Retina, 15-inch, Mid 2015)
Processor 2,5 GHz Intel Core i7
Memory 16 GB 1600 MHz DDR3
```
```bash
➜  bigopool git:(master) go test -bench=.  -benchmem=true -cpu=1,2,4,8,16                                                                                                                                                                    
goos: darwin                                                                                                                                                                                                                               
goarch: amd64
pkg: github.com/bigodines/bigopool
Benchmark1Workers1Queue                  1000000              4202 ns/op             256 B/op          0 allocs/op
Benchmark1Workers1Queue-2                1000000              2139 ns/op             176 B/op          0 allocs/op
Benchmark1Workers1Queue-4                1000000              1965 ns/op             160 B/op          0 allocs/op
Benchmark1Workers1Queue-8                1000000              1896 ns/op             152 B/op          0 allocs/op
Benchmark1Workers1Queue-16                500000              2004 ns/op             143 B/op          0 allocs/op
Benchmark5Workers1000Queue               1000000              4833 ns/op             323 B/op          1 allocs/op
Benchmark5Workers1000Queue-2             1000000              1111 ns/op             136 B/op          0 allocs/op
Benchmark5Workers1000Queue-4             1000000              1204 ns/op             163 B/op          0 allocs/op
Benchmark5Workers1000Queue-8             1000000              1628 ns/op             169 B/op          0 allocs/op
Benchmark5Workers1000Queue-16            1000000              1485 ns/op             160 B/op          0 allocs/op
Benchmark10Workers100Queue               2000000               964 ns/op              86 B/op          0 allocs/op
Benchmark10Workers100Queue-2             1000000              1079 ns/op             137 B/op          0 allocs/op
Benchmark10Workers100Queue-4             2000000               878 ns/op             145 B/op          0 allocs/op
Benchmark10Workers100Queue-8             1000000              1154 ns/op             125 B/op          0 allocs/op
Benchmark10Workers100Queue-16            1000000              1005 ns/op             110 B/op          0 allocs/op
Benchmark20Workers200Queue               2000000              1210 ns/op             109 B/op          0 allocs/op
Benchmark20Workers200Queue-2             1000000              1066 ns/op             136 B/op          0 allocs/op
Benchmark20Workers200Queue-4             2000000               796 ns/op             161 B/op          0 allocs/op
Benchmark20Workers200Queue-8             1000000              1006 ns/op             115 B/op          0 allocs/op
Benchmark20Workers200Queue-16            1000000              1118 ns/op             132 B/op          0 allocs/op
Benchmark20Workers10000Queue             1000000              1738 ns/op             111 B/op          0 allocs/op
Benchmark20Workers10000Queue-2           1000000              1186 ns/op             123 B/op          0 allocs/op
Benchmark20Workers10000Queue-4           2000000               872 ns/op             168 B/op          0 allocs/op
Benchmark20Workers10000Queue-8           1000000              1479 ns/op             162 B/op          0 allocs/op
Benchmark20Workers10000Queue-16          2000000               921 ns/op              94 B/op          0 allocs/op
Benchmark100Workers10000Queue            1000000              2155 ns/op             176 B/op          0 allocs/op
Benchmark100Workers10000Queue-2          1000000              1189 ns/op             126 B/op          0 allocs/op
Benchmark100Workers10000Queue-4          2000000               816 ns/op             137 B/op          0 allocs/op
Benchmark100Workers10000Queue-8          1000000              1172 ns/op             152 B/op          0 allocs/op
Benchmark100Workers10000Queue-16         2000000              1107 ns/op             128 B/op          0 allocs/op
PASS
ok      github.com/bigodines/bigopool     86.650s

```

## License

MIT
