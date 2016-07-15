package main

import (
    "fmt"
)

type Context struct {
    done chan bool
}

func (ctx *Context) Done() <-chan bool {
    return ctx.done
}

func (ctx *Context) Stop() {
    ctx.done <- true
}

func NewContext() *Context {
    ctx := new(Context)
    done := make(chan bool)
    ctx.done = done
    return ctx
}

type Counter struct {
    ctx Context
    c chan int
    count int
}

func NewCounter(ctx *Context, start int) *Counter {
    counter := new(Counter)
    counter.ctx = *ctx
    counter.c = make(chan int)
    counter.count = start

    go func(){
        for {
            select {
            case counter.c <- counter.count:
                counter.count += 1
            case <- counter.ctx.Done():
                return
            }
        }
    }()
    return counter
}

func (c *Counter) GetSource() <-chan int {
    return c.c
}

func (c *Counter) Stop(){
    c.ctx.Stop()
}

func main(){
    ctx := NewContext()
    c1 := NewCounter(ctx, 5)
    c2 := NewCounter(ctx, 10)

    count1 := c1.GetSource() // returns a read only channel
    count2 := c2.GetSource()

    for {
        n := <- count1
        m := <- count2
        if m > 100 {
            ctx.Stop()
            break
        } else {
            fmt.Printf("Counter values: [1]:%d\t[2]:%d\n",n, m)
        }
    }
}