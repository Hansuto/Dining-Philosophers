# Dining-Philosophers
Dijkstra's classic concurrency exercise written in Go. 

## Compilation Steps

My program usees GoLang version 1.11.4. To install the latest version of Go, run the 
following command:

	sudo snap install go --classic
		
To verify installation run the following command:

	go version
		
You should see this output if installation comleted successfully:

	go version go1.11.4 linux/amd64
    
In my implementation you're allowed to pick the number of philosophers 'n' where 1 ≤ n < ∞. 
To run the program, enter the following command:

	go run philosophers.go n
		

## Explanation
  To avoid two philosophers holding the same chopstick at the same time, the chopsticks
are implemented as channels. When using buffered channels in Golang, when 1 value is in the
channel, it can be read. Once the value is read, it is popped from the channel, leaving it empty.
Channels are blocking by nature so once a channel is empty, any other threads trying to read
from that channel must wait. Each philosopher has a chopstick and has access to only one
neighbor’s chopstick. If the chopstick channel has a value, the Philosopher reads from the
channel therefore making that chopstick unusable to all other Philosophers until another value
is appended to the chopstick channel. A value is added back to the channel once a Philosopher
is done eating.

  To prevent deadlocking in my solution, I use a selection statement in my
pickupChopsticks() function. After a Philosopher waits for their chopstick, they must wait for
their neighbor’s chopstick before eating. If all Philosophers pick up their chopstick at the same
time, then they all will be waiting for their neighbor, causing deadlock. Because I’m using
channels, the only way to break out of waiting is to have another execution path available at
the same time. The select statement grants this by allowing the Philosophers to select the
“exit” case if they must wait for their neighbor’s chopstick. The second select case will always
be a viable exit condition if the neighbor’s chopstick is unavailable after 3 seconds. If a
philosopher reaches the “exit” case, then they will set their chopstick back down (by adding a
value back into the channel) and try to pick up chopsticks once again.

  To prevent starvation in my solution I implemented a “hunger” attribute to each
philosopher which starts at 0. Every time a philosopher gets hungry, their hunger increments by
one. Once they finish eating, their hunger goes back to 0. If a Philosopher’s hunger becomes
greater than the number of Philosophers at the table, the Philosopher’s neighbor must now
wait before picking up their chopstick, so the original philosopher has time to pick it up and eat.
The Philosophers are set up vary similarly to a linked list, so scalability was very simple.
Instead of hard-coding how many philosophers to spawn, I captured the command line
argument and just made that many Philosophers sit at the table. 
