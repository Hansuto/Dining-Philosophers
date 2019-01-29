// Written By Christopher Taliaferro | Ch119541

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var numPhilosophers int

type Philosopher struct {
	number int
	hunger int
	wait bool
	chopstick chan bool
	neighbor *Philosopher
}

// Algorithm to determine if chopsticks are available to use.
func (current *Philosopher) pickupChopsticks() {
	fmt.Println(current.number, "is now hungry.")
	current.hunger++

	// If the current philosopher's wait flag is set to true, then it's
	// neighbor is starving. Therefore this philosopher will wait.
	if current.wait == true {
		time.Sleep(2 * time.Second)
		current.wait = false
	}

	// Waiting for the current Philosopher's chopstick to become available.
	<-current.chopstick

	// If the current Philosopher's hunger is greater than the number of
	// Philosophers at the table, the current Philospher's neighbor now
	// must wait for the current Philosopher to eat.
	if current.hunger >= 3 {
		current.neighbor.wait = true
	}

	// We use a selection because we must first check our neighbor's chopstick.
	// Since we are dealing with channels, if the chopstick is unavailable we
	// reach a blocking condition. To counter this, the select statement checks
	// all cases even if some are blocking to determine the execution path.
	// This allows the current Philosopher to set it's chopstick back down after
	// waiting for its neighbor's chopstick for 3 seconds. This prevents deadlock.
	select {
	case <-current.neighbor.chopstick:
		return
	case <-time.After(3 * time.Second):
		current.chopstick <- true
		current.pickupChopsticks()
	}
}

func (current *Philosopher) start() {
	fmt.Println(current.number,"is now thinking.")

	// Each Philosopher will run through these actions forever.
	for {
		time.Sleep(time.Duration(rand.Intn(1e5)) * time.Microsecond)
		current.pickupChopsticks()
		fmt.Println(current.number,"is now eating.")
		time.Sleep(time.Duration(rand.Intn(1e5)) * time.Microsecond)
		fmt.Println(current.number,"is now thinking.")
		current.chopstick <- true
		current.neighbor.chopstick <- true
		current.hunger = 0
	}
}

// Initializes Philosopher's number, assigns it a chopstick, and defines it's neighbor.
func makePhilosopher(number int, neighbor *Philosopher) *Philosopher {
	// Make note of "make(chan bool, 1)", this is our blocking condition in this program.
	// The channel can hold only 1 value. This means threads must wait when the channel
	// is empty. This also ensures only one thread can hold a chopstick at any given time.
	current := &Philosopher{number, 0, false, make(chan bool, 1), neighbor}

	// If the chopstick channel has a value then it's available to grab.
	current.chopstick <- true

	return current
}

func introduction() {
	fmt.Println("Be sure to press 'Enter' when you want to terminate execution")
	fmt.Println()
	time.Sleep(2 * time.Second)
	fmt.Printf("3 ")
	time.Sleep(1 * time.Second)
	fmt.Printf("2 ")
	time.Sleep(1 * time.Second)
	fmt.Printf("1 ")
	time.Sleep(1 * time.Second)
	fmt.Printf("GO! \n\n")
	fmt.Println()
}

func inputFail () {
	fmt.Println()
	fmt.Println("Please provide a valid command line argument. [1 < n < ∞]")
	fmt.Println()
	fmt.Println("EXAMPLE: go run problem2v4.go <n>")
	fmt.Println()
}

func main() {

	// Input Validation
	if len(os.Args) < 2 {
		inputFail()
		return
	}

	// Grabs number of philosophers from command line argument.
	arg := os.Args[1]
	numPhilosophers, _ := strconv.Atoi(arg)

	// Input Validation
	if numPhilosophers < 2 {
		inputFail()		
		return
	}

	fmt.Println()
	fmt.Println(numPhilosophers, "Philosophers are sitting at the table.")

	// Creates an array of size 'numPhilosophers' with data type 'Philosopher'.
	philosophers := make([]*Philosopher, numPhilosophers)

	// Creates variable of custom type "Philosopher"
	var current *Philosopher

	// Initializes each philosopher and adds to philosophers array.
	for i := 0; i < numPhilosophers; i++ {
		current = makePhilosopher(i+1, current)
		philosophers[i] = current
	}

	// Defines philosopher[0]'s neighbor to close the loop.
	philosophers[0].neighbor = current

	introduction()

	// Spawns a go routine for each philosopher.
	for i:= 0; i < numPhilosophers; i++ {
		go philosophers[i].start()
	}

	// Reads keyboard input to terminate program.
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
