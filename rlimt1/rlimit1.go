/*
    rlimit.go

    An example that illustrates how to limit the maximum number of go routines
    spawned by a program.

 */

package main

import (
    "fmt"
    "time"
    "math/rand"
)

// process takes in an int and prints its square
func process(ch chan int, id int, done chan int) {
    for x := range ch{
        t := time.Duration(rand.Intn(20))*time.Millisecond
        time.Sleep(t)
        fmt.Printf("%d: The square of %d is %d (%d ms delay)\n",id, x, x*x,int(t / 1000000))
        done <- 1
    }
}

func main(){
    ch := make(chan int)
    done := make(chan bool)
    var count int

    // Start sending data down the channel for processing
    go func(){
        for i := 1; i <= 105; i++ {
            ch <- i
        }
        done <- true
    }()

    // Create data processors on multiple threads
    n := make(chan int,101)
    for i:= 0; i<10; i++ {
        go func(ident int){
            process(ch, ident, n)
        }(i + 100)
    }
    <- done

    //
    for count < 105 {
        count += <- n
    }
}