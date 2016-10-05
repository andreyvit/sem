# Go semaphore

Implements a semaphore in Go, allowing to limit the concurrency of goroutines.

Note that this builds on top of Go channels, which is a recommended pattern used inside stdlib, but may be somewhat heavyweight.

Installation:

    go get github.com/andreyvit/sem

[Docs on godoc.org](https://godoc.org/github.com/andreyvit/sem)


## Example

