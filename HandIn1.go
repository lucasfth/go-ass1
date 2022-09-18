package main

import {
    "fmt"
    "time"
}

func main() {
    f0 := make(chan int)
    f1 := make(chan int)
    f2 := make(chan int)
    f3 := make(chan int)
    f4 := make(chan int)
    
    
}


func getUsed(index int, c chan int) {
    isUsed bool = false
    for true {
        m <- c
        if (!isUsed) {
            c <- m
        } else {
            c <- 10
        }
    }
}
