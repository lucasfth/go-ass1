package main

import (
    "fmt"
    "time"
)

type fork struct {
    index int
    ch chan bool
}

type phil struct {
    index int
    
    leftCh chan bool
    rightCh chan bool
    
    forkR fork
    forkL fork
}

func main() {
    chf0 := make(chan bool)
    chf1 := make(chan bool)
    chf2 := make(chan bool)
    chf3 := make(chan bool)
    chf4 := make(chan bool)
    
    fmt.Println("some")
    
    chans := [5] chan bool {chf0, chf1, chf2, chf3, chf4}
    
    var forks = new([5]fork)
    for i:=0; i<5; i++ {
        var temp fork
        temp.index = i
        temp.ch = chans[i]
        forks[i] = temp
        go forkCom(temp)
    }
    
    for i:=0; i<5; i++ {
        var temp phil
        temp.index = i
        temp.leftCh = forks[(i+1)%5].ch
        temp.rightCh = forks[i].ch
        
        temp.forkL = forks[(i+1)%5]
        temp.forkR = forks[i]
        go philEat(temp)
    }
    time.Sleep(6 * time.Second)
}

func philEat(p phil) {
    for i:=0; i<3; i++ {
        for true {
            c0 := <- p.leftCh
            c1 := <- p.rightCh
            if (c0 == true && c1 == true) {
                p.leftCh <- true
                p.rightCh <- true
                fmt.Println("Philosipher", p.index, "has eaten", i+1, "times")

                p.leftCh <- false
                p.rightCh <- false
                break
            } else {
                fmt.Println("Philosipher", p.index, "is thinking")

            }
        }
    }
}

func forkCom(f fork) {
    for true {
        f.ch <- true
        r := <- f.ch
        if (r == true) {
            fmt.Println("\t", f.index, "is being used")
            for true {
                f.ch <- false
                r := <- f.ch
                if (r == false) {
                    break
                }
            }
        }
    }
}