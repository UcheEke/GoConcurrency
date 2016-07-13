/*
    leakybuffer.go

    The client can either reuse an existing buffer from a list or create a new one for processing on the server
    Program uses buffered channels to make the workflow clearer.

 */

package main

import (
    "strings"
    "fmt"
    "os"
    "bufio"
)

type Buffer struct {
    data string
}

func (b *Buffer) Load(data string){
    b.data = data
}

func (b *Buffer) Process(){
    fmt.Println("Processed data:", strings.ToUpper(b.data))
}

var available = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)
var quit = make(chan bool)

func dataGenerator() chan string {
    filename := "/home/ekeu/Programming/Projects/GoConcurrency/leakybuffer/words.txt"
    fd, err := os.Open(filename)
    if err != nil {
        panic(err)
    }

    ch := make(chan string)
    sc := bufio.NewScanner(fd)

    go func(){
        for sc.Scan() {
            ch <- sc.Text()
        }
        quit <- true
        close(ch)
    }()

    return ch
}

func client() {
    ch := dataGenerator()
    clientLoop:
    for {
        var b *Buffer
        // Grab a buffer if available; allocate if not.
        select {
        case b = <-available:
        // Got one; nothing more to do.
        default:
        // None free, so allocate a new one.
            b = new(Buffer)
        }

        select {
        case data := <-ch:
            b.Load(data)         // Read next buffer
            serverChan <- b      // Send to server.
        case <- quit:            // Receive signal from dataGenerator? Clean up and close
            break clientLoop
        }
    }
}

func server() {
    serverLoop:
    for {
        select{
        case b := <- serverChan:
            b.Process()
                // Reuse buffer if there's room.
                select {
                case available <- b:
                // Sent buffer back for reuse.
                default:
                // Available channel is full, drop buffer for GC to clean up.
                }
        case <- quit:
            break serverLoop
        }
    }
    close(available)
    close(serverChan)
}

func main(){
    go client()
    server()
}