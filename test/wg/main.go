package main

import "sync"

func main(){
    var wg sync.WaitGroup
    wg.Add(1)
    go func(){
        defer func(){
            if err := recover(); err != nil {
                println("recover in wg")
            }
            // 挂了也要Done，否则会阻塞
            wg.Done()
        }()
        p()
        wg.Done()
        println("hello")
    }()
    wg.Wait()
    println("main end")
}

func p(){
    panic("panic in p")
}