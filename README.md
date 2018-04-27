# FOIER

FOIER is a learning tool designed to illustrate how concurrency can be leveraged
in Golang for rapid reconnaissance and discovery. It is being presented for the
first time at the [Atlantic Security Conference](https://atlseccon.com) on April
26th, 2018.

## Getting started

* First you need the Go environment
  * [https://golang.org/doc/install](https://golang.org/doc/install)
  * It's available in `brew` and Linux package managers, which will be easier.
* Clone this repository(in `$GOPATH/src/github.com/invisiblethreat/`)
  * `git clone https://github.com/invisiblethreat/foier.git`
* Get `dep` so you can install your dependencies
  * `go get -u github.com/golang/dep/cmd/dep`
* Get all of the project dependencies
  * `dep ensure`
* Build the project!
  * `go build`
* Run the binary
  * `foier --help`


## Demo Success!

To give an idea of the performance of this project, I was able to download 7000
files in under two seconds. The two hosts that I used were a $5 instance from
Vultr and a $40 instance from Linode.

## Further Learning

* [Go By Example](https://gobyexample.com/)
* [A Tour of Go](https://tour.golang.org/welcome/1)
* [Concurrency is not Parallelism, Rob Pike](https://talks.golang.org/2012/waza.slide#1)
* [Visualizing Concurrency in Go, Ivan Daniluk](http://divan.github.io/posts/go_concurrency_visualize/)
* [7 common mistakes in Go and when to avoid them, Steve Francia](https://www.youtube.com/watch?v=29LLRKIL_TI)
