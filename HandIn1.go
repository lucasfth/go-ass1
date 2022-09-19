package main

import (
	"fmt"
	"time"
)

type Forks struct {
	name     string
	inUse    bool
	outChan  chan string
	leftCin  chan string
	rightCin chan string
}

type Phils struct {
	name       string
	leftHanded bool
	leftCin    chan string
	rightCin   chan string
	leftCout   chan string
	rightCout  chan string
}

// FUNC TO CREATE PHILS
func createPhils(name string, lh bool, l chan string, r chan string, li chan string, ri chan string) Phils {
	phil := Phils{name, lh, r, l, ri, li}
	return phil
}

// FUNC TO CREATE FORKS
func createForks(name string) Forks {
	outerChan := make(chan string)
	leftCin := make(chan string)
	rightCin := make(chan string)
	fork := Forks{name, false, outerChan, leftCin, rightCin}
	return fork
}

// START THE DINING
func (phil Phils) startDining() {
	fmt.Println(phil.name, "has sat down at the table")
	time.Sleep(1 * time.Second)
	for i := 0; i < 3; i++ {
		phil.requestFork(i + 1)
	}
}

// PLACE FORK AT TABLE
func (fork Forks) placeForks() {
	fmt.Println(fork.name, "has been placed on the table")
	fork.outChan <- fork.name
	for true {
		select {
		case one := <-fork.leftCin:
			fork.outChan <- one
		case two := <-fork.rightCin:
			fork.outChan <- two
		}
	}
}

/*
PHILOSOPHERS REQUESTING FORKS

MAKING TWO PHILOSOPHERS RIGHTHANDED - WHICH MEANS THE PHILOSOPHERS WILL TRY AND GRAB WITH THEIR RIGHT HAND FIRST -  AND ADDING 3 SECONDS SLEEP TO THEM
WHILE THE OTHER THREE ARE LEFTHANDED - WHICH MEANS THE PHILOSOPHERS WILL TRY AND GRAB WITH THEIR LEFT HAND FIRST - WITH NO SLEEP
MAKES SURE THAT THERE IS NEVER A DEADLOCK. THE RIGHTHANDED PHILOSOPHERS WILL ALWAYS HAVE TO WAIT FOR THE LEFTHANDED TO FINISH UP FIRST IN RANDOM ORDER
SO THEY WILL SIT BY AND WAIT FOR THE LEFTHANDED TO DROP THEIR FORK
*/
func (phil Phils) requestFork(i int) {
	if phil.leftHanded {
		lh := <-phil.leftCin
		rh := <-phil.rightCin
		phil.eat(lh, rh, i)
	} else {
		time.Sleep(3 * time.Second)
		rh := <-phil.rightCin
		lh := <-phil.leftCin
		phil.eat(lh, rh, i)
	}
}

func (phil Phils) eat(lh string, rh string, i int) {
	time.Sleep(3 * time.Second)
	fmt.Println(phil.name, "is eating for the", i, "time")
	phil.releaseFork(lh, rh)
}

// THE PHILOSOPHERS MAKES THE FORKS AVAILABLE FOR USE AGIAN FOR OTHER PHILOSOPHERS
func (phil Phils) releaseFork(lh string, rh string) {
	phil.think()
	phil.leftCout <- lh
	phil.rightCout <- rh
}

// PHILOSOPHERS THINKS
func (phil Phils) think() {
	fmt.Println(phil.name, "is thinking")
}

func main() {
	namesP := []string{"phil0", "phil1", "phil2", "phil3", "phil4"}
	namesF := []string{"Fork0", "Fork1", "Fork2", "Fork3", "Fork4"}
	Philosophers := make([]Phils, 5)
	Forks := make([]Forks, 5)

	//CREATE FORKS
	for i := 0; i < 5; i++ {
		Forks[i] = createForks(namesF[i])
	}

	//CREATE PHILS
	for i := 0; i < 5; i++ {
		if i == 0 {
			Philosophers[i] = createPhils(namesP[i], false, Forks[4].outChan, Forks[0].outChan, Forks[4].rightCin, Forks[0].leftCin)
		} else if i == 3 {
			Philosophers[i] = createPhils(namesP[i], false, Forks[i-1].outChan, Forks[i].outChan, Forks[i-1].rightCin, Forks[i].leftCin)
		} else {
			Philosophers[i] = createPhils(namesP[i], true, Forks[i-1].outChan, Forks[i].outChan, Forks[i-1].rightCin, Forks[i].leftCin)
		}
	}

	for i := 0; i < 5; i++ {
		go Philosophers[i].startDining()
	}

	time.Sleep(3 * time.Second)

	for i := 0; i < 5; i++ {
		go Forks[i].placeForks()
	}

	time.Sleep(45 * time.Second)

}
