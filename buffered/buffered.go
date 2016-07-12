/*
    Illustrating a basic use of buffered channels as semaphores
 */

package main

import (
    "strings"
    "fmt"
)

func process(word string, out chan int) {
    fmt.Printf("Input: %s\n", word)
    word = strings.ToUpper(word)
    fmt.Printf("Output: %s\n", word)
    out <- 1
}

func main() {
    names := []string{
        "abigail",
        "bjorn",
        "chimanda",
        "dunun",
        "efere",
        "francois",
        "greta",
        "hanan",
        "iqbal",
        "janislaw",
        "kelena",
        "louis",
        "mohammad",
        "nali",
        "opobo",
    }

    queue := make(chan int, 10) // Set a limit of 10 parallel processes
    var count int
    ch := make(chan int)

    for _, name := range names {
        queue <- 1 // queue will block once the buffer is full
        go func(n string){
            process(n, ch)
        }(name)
        <- queue // queue discards the value freeing up space for new processes to spawn
    }

    // Blocking until all processes complete
    for count < len(names) - 1 {
        count += <- ch
    }
}