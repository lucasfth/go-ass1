package main

import (
	"fmt"
	"time"
)

// REPRESENTING A FORK - THE FORK HAS A NAME, ONE CHANNEL "GOING OUT" TO THE PHILOSOPHERS WHO CAN USE IT
// AND TWO CHANNELS "COMING IN", ONE FROM EACH PHILOSOPHER WHO CAN USE IT
// THE CHANNELS ARE USED TO COMMUNICATE AVAILABILITY - THE FORK WILL USE THE CHANNEL GOING OUT WHEN IT IS AVAILABLE
// AND IT WILL RECEIVE A STATEMENT COMING IN WHEN IT IS AVAILABLE AGAIN AND BEEN "RETURNED" BY THE PHILOSOPHER
type Forks struct {
	name     string
	outChan  chan string
	leftCin  chan string
	rightCin chan string
}

// REPRESENTING A PHILOSOPHER
// THE PHILOSOPHER IS EITHER LEFT OR RIGHT HANDED; DECIDING HOW HE WILL PREFER CHOOSING THE FIRST FORK
// THEY HAVE TWO CHANNELS GOING OUT AND TWO CHANNELS COMING IN FROM THE FORKS THAT THEY CAN USE
// THESE CHANNELS ARE DESCRIBED ABOVE
type Phils struct {
	name       string
	leftHanded bool
	leftCin    chan string
	rightCin   chan string
	leftCout   chan string
	rightCout  chan string
}

// FUNCTION TO CREATE PHILOSOPHERS
func createPhils(name string, leftHanded bool, leftChannelIn chan string, rightChannelIn chan string, leftChannelOut chan string, rightChannelOut chan string) Phils {
	philosopher := Phils{name, leftHanded, rightChannelIn, leftChannelIn, rightChannelOut, leftChannelOut}
	return philosopher
}

// FUNCTION TO CREATE FORKS
func createForks(name string) Forks {
	outChan := make(chan string)
	leftChannelIn := make(chan string)
	rightChannelIn := make(chan string)
	fork := Forks{name, outChan, leftChannelIn, rightChannelIn}
	return fork
}

// START THE DINING PROCESS
func (phil Phils) startDining() {
	fmt.Println(phil.name, "has sat down at the table")
	time.Sleep(1 * time.Second)
	for i := 0; i < 3; i++ {
		phil.requestFork(i + 1)
	}
	fmt.Println("____________________________________", phil.name, "has finished eatings and left the table", "____________________________________________________")
}

// PLACE FORK AT TABLE
// SENDING OUT AVAILABLEBILITY
// CONTINOUSLY WAITING FOR A RESPONSE FROM PHILOSOPHERS RETURNING IT AND THEN IT WILL SEND OUT AGAIN THAT IT IS ABVAILABLE
// THE SELECT STATEMENT MAKES SURE THAT THERE IS NO DEADLOCK FOR THE FORKS JUST WAITING FOR ONE PHILOSOPHER WHILE THE OTHER ACTUALLY HAS IT
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

// PHILOSOPHERS REQUESTING FORKS
// MAKING TWO PHILOSOPHERS RIGHTHANDED - WHICH MEANS THE PHILOSOPHERS WILL TRY AND GRAB THE FORK TO THE RIGHT SIDE FIRST -  AND ADDING 3 SECONDS SLEEP TO THEM
// WHILE THE OTHER THREE ARE LEFTHANDED - WHICH MEANS THE PHILOSOPHERS WILL TRY AND GRAB THE FORK TO THE LEFT FIRST - WITH NO SLEEP TIME
// MAKES SURE THAT THERE IS NEVER A DEADLOCK. THE RIGHTHANDED PHILOSOPHERS WILL ALWAYS HAVE TO WAIT FOR THE LEFTHANDED TO FINISH UP FIRST THE FIRST TIME THE FUNCTION IS REACHED
// AFTER THIS THE SELECT STAEMENT WILL MAKE SURE THERE IS NEVER A DEADLOCK
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

// PRINT EAT
func (phil Phils) eat(lh string, rh string, i int) {
	time.Sleep(3 * time.Second)
	fmt.Println(phil.name, "is eating for the", i, "time")
	phil.releaseFork(lh, rh)
}

// THE PHILOSOPHERS MAKES THE FORKS AVAILABLE FOR USE AGIAN FOR OTHER PHILOSOPHERS
// BY SENDING A STATEMENT BACK TO THE FORK BY A CHANNEL, THAT IT IS AVAILABLE
func (phil Phils) releaseFork(lh string, rh string) {
	phil.think()
	phil.leftCout <- lh
	phil.rightCout <- rh
}

// PRINT THINKS
func (phil Phils) think() {
	fmt.Println(phil.name, "is thinking")
}

func main() {
	// NAMES FOR FORKS AND PHILOSOPHERS
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

	//GOROUTINES ARE STARTED HERE FOR PHILOSOPHERS
	for i := 0; i < 5; i++ {
		go Philosophers[i].startDining()
	}

	time.Sleep(3 * time.Second)

	//GOROUTINES FOR FORKS ARE STARTED HERE
	for i := 0; i < 5; i++ {
		go Forks[i].placeForks()
	}
	time.Sleep(45 * time.Second)
}
