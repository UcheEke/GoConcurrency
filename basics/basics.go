package main

import (
    "fmt"
    "time"
)

func delays() {
    message := "This comes from a go routine!\n"
    go func(){
        fmt.Printf(message)
    }()  // Always execute the closure after its definition
    fmt.Println("\ndelays()")
    time.Sleep(1 * time.Second)
    fmt.Printf("This comes from the 'delays()' function thread\n")
}

func blockingChannels() {
    ch := make(chan int)
    go func(){
       fmt.Println("This is a messsage from a go routine")
       ch <- 1
    }()
    fmt.Println("\nblockingChannels()")
    fmt.Println("This is a message from the 'blockingChannels' function thread")
    <- ch
    fmt.Println("Channel no longer blocks. 'blockingChannels()' Exiting")
}

func signalling(){
    ch := make(chan int)
    done := make(chan bool)

    go func(){
        var i int = 1
        for {
            select{
            case ch <- i:
                i++
            case <- done:
                return
            }
        }
    }()

    fmt.Println("\nsignalling()")
    for n:=0; n<20; n++ {
        fmt.Printf("%d\n", <-ch)
    }
    done <- true
    close(done)
    close(ch)
}

func main(){
    delays()
    blockingChannels()
    signalling()
}


