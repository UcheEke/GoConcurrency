/*
    structchans.go
    Example of using structures that contain channels themselves, processed by a pool of concurrent worker
    functions without the use of waitGroups or Muxes

 */

package main

import (
    "fmt"
    "math/rand"
    "time"
)

const (
    maxRequestHandlers = 10
)

// Request: Basic data structure that comprises of a data payload, a function
// to process that data and a channel to receive the result of the processing
type Request struct {
    id int
    data []int
    nf NamedFunction
    results chan int
}

// Process: Processes the request and returns the result on the internal channel
func (r *Request) Process() {
    r.results <- r.nf.fn(r.data)
}

// GetResults: Decorative getter method
func (r *Request) GetResults() int {
    return <- r.results
}

// NamedFunction: Allows the function used by a Request to be named in a result
type NamedFunction struct {
    name string
    fn func([]int)int
}

// Basic functions used by Requests:
func sum(v []int) int {
    var s int
    for _, elem := range v {
        s += elem
    }
    return s
}

func prod(v []int) int {
    var p int = 1
    for _, elem := range v {
        p *= elem
    }
    return p
}

func max(v []int) int {
    m := v[0]
    for i, n := range v {
        if i == 0 {
            continue
        }
        if m < n {
            m = n
        }
    }
    return m
}

// createData: generates a random payload for a Request
func createData(v []int) []int {
    rand.Seed(time.Now().UnixNano())
    l := len(v)
    if l < 3 {
        return []int{0,0,0}
    }

    res := make([]int,3)
    for i,k := range rand.Perm(l)[0:3]{
        res[i] = v[k]
    }
    return res
}

// handle: Calls the Process method of each request. This will be run as a concurrent function
// on maxRequestHandlers go routines
func handle(q chan *Request) {
    for r := range q {
        r.Process()
    }
}

func main() {
    // Request creation variables
    numbers := []int{1,2,3,4,5,6,7,8,9}

    Sum := NamedFunction{
        name: "sum",
        fn: sum,
    }

    Prod := NamedFunction{
        name: "prod",
        fn : prod,
    }

    Max := NamedFunction{
        name: "max",
        fn: max,
    }

    funcs := []NamedFunction{
        Sum,
        Prod,
        Max,
    }

    // Create the requests
    var reqs = make([]*Request, 0)
    for i:=0; i<30; i++{
        r := &Request{
            id: rand.Intn(100),  // Give it a random ID
            data: createData(numbers), // add some data
            nf: funcs[rand.Intn(len(funcs))], // pick a function to use
            results: make(chan int), // create the results channel
        }
        reqs = append(reqs, r)
    }

    // Create a queue for the requests
    queue := make(chan *Request, maxRequestHandlers)

    // Create multiple threads to handle a maximum number of requests at a time
    for i:=0;i<maxRequestHandlers; i++{
        go handle(queue)
    }

    for _, r := range reqs {
        // Send the request for processing
        queue <- r
        // Wait for the results
        fmt.Printf("Result from Request ID[%d]: %d\t\t(Data: %v\tFunction: %q)\n",
            r.id,
            r.GetResults(), // This will block until it returns something
            r.data,
            r.nf.name,
        )
    }
}