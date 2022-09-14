package main

import (
	"fmt"
	"time"
)

type fork struct {
	index  int
	isUsed bool

	inChannel     chan int
	rightPhilChan chan bool
	leftPhilChan  chan bool
}

type phil struct {
	index int

	leftForkOut  chan int
	rightForkOut chan int
	leftForkIn   chan bool
	rightForkIn  chan bool
}

func main() {
	fork0in := make(chan int)
	fork1in := make(chan int)
	fork2in := make(chan int)
	fork3in := make(chan int)
	fork4in := make(chan int)

	fork0leftOut := make(chan bool)
	fork1leftOut := make(chan bool)
	fork2leftOut := make(chan bool)
	fork3leftOut := make(chan bool)
	fork4leftOut := make(chan bool)
	fork0RightOut := make(chan bool)
	fork1RightOut := make(chan bool)
	fork2RightOut := make(chan bool)
	fork3RightOut := make(chan bool)
	fork4RightOut := make(chan bool)

	inChannels := [5]chan int{fork0in, fork1in, fork2in, fork3in, fork4in}
	leftOutChannels := [5]chan bool{fork0leftOut, fork1leftOut, fork2leftOut, fork3leftOut, fork4leftOut}
	rightOutChannels := [5]chan bool{fork0RightOut, fork1RightOut, fork2RightOut, fork3RightOut, fork4RightOut}

	for i := 0; i < 5; i++ {
		var temp fork
		temp.index = i
		temp.inChannel = inChannels[i]
		temp.rightPhilChan = rightOutChannels[i]
		temp.leftPhilChan = leftOutChannels[i]

		go forkCom(temp)
	}

	for i := 0; i < 5; i++ {
		var temp phil
		temp.index = i

		temp.rightForkOut = inChannels[i]
		temp.leftForkOut = inChannels[(i+1)%5]

		temp.rightForkIn = leftOutChannels[i]
		temp.leftForkIn = rightOutChannels[(i+1)%5]

		go philEat(temp)
	}

	time.Sleep(12 * time.Second)
}

func philEat(p phil) {
	var timesEaten int = 1

	for timesEaten < 4 {

		p.rightForkOut <- p.index
		m0 := <-p.rightForkIn
		if m0 == true {
			p.leftForkOut <- p.index
			m1 := <-p.leftForkIn
			if m1 == true {
				fmt.Println("Philosopher", p.index, "has eaten", timesEaten, "times-----------------------------")
				timesEaten = timesEaten + 1
				p.leftForkOut <- 10
			}
			p.rightForkOut <- 10
		}
		fmt.Println("Philosopher", p.index, "is thinking")
	}
}

func forkCom(f fork) {
	for true {
		m0 := <-f.inChannel
		if f.index == (m0) {
			if f.isUsed {
				f.leftPhilChan <- false
			} else {
				f.isUsed = true
				f.leftPhilChan <- true

			}
		}
		if f.index == (m0+1)%5 {
			if f.isUsed {
				f.rightPhilChan <- false
			} else {
				f.isUsed = true
				f.rightPhilChan <- true

			}
		}
		if 10 == m0 {
			f.isUsed = false
			fmt.Println("\tFork:", f.index, f.isUsed)
		}
	}
}
