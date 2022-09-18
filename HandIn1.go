package main

import (
    "fmt"
    "time"
)

func main() {
    f0 := make(chan int)
    f1 := make(chan int)
    f2 := make(chan int)
    f3 := make(chan int)
    f4 := make(chan int)
    
    channels := [5] chan int {f0, f1, f2, f3, f4}
    
    for i := 0; i < 5; i++ {
        go getUsed (i, channels[i])
    }
    
    for i := 0; i < 5; i++ {
        go goEat(i, channels[i], channels[(i+1)%5])
    }
    
    time.Sleep(10 * time.Second)
}

func goEat(index int, rC chan int, lC chan int) {
    var i int = 0
    for i < 3 {
        rC <- index
        rM := <- rC
        if (rM == index) {
            lC <- index
            lM := <- lC
            if (lM == index) {
                fmt.Println("Phil", index, "eaten", i+1, "time(s)\t\t\t", i+1)
                i++
                time.Sleep(5 * time.Millisecond)
                lC <- 10
            }
            rC <- 10
        }
        //fmt.Println("Phil", index, "is thinking")
        time.Sleep(5 * time.Millisecond)
    }
}


func getUsed(index int, c chan int) {
    var isUsed bool = false
    for true {
        m := <- c
        if (m == 10) {
            isUsed = false
            //fmt.Println("Fork", index, "is free")
        } else if (!isUsed ) {
            isUsed = true
            //fmt.Println("Fork", index, "is used by", m)
            c <- m
        } else {
            c <- 10
        }
    }
}
