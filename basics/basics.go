package main

import (
    "fmt"
    "time"
)

func delays() {
    message := "This comes from a go routine!\n"
    go func(){
        fmt.Printf(message)
    }()
    time.Sleep(1 * time.Second)
    fmt.Printf("This comes from the 'delays' function thread\n")
}

func blockingChannels() {
    ch := make(chan int)
    go func(){
       fmt.Println("This is a messsage from a go routine")
       ch <- 1
    }()

    fmt.Println("This is a message from the 'blockingChannels' function thread")
    <- ch
    fmt.Println("Channel no longer blocks. 'blockingChannels()' Exiting")
}

func main(){
    delays()
    blockingChannels()
}


