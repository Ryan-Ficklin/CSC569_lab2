package main

import (
	"fmt"
	"github.com/Ryan-Ficklin/CSC569_lab2/shared"
	"math/rand"
	"net/rpc"
	"os"
	"strconv"
	//"strings"
	"sync"
	"time"
)

const (
	MAX_NODES  = 8
	X_TIME     = 1
	Y_TIME     = 2
	Z_TIME_MAX = 100
	Z_TIME_MIN = 10
)
var (
  self_node shared.Node
  self_mutex sync.Mutex
)

// Send the current membership table to a neighboring node with the provided ID
func sendMessage(server *rpc.Client, id int, membership shared.Membership) {
	//TODO
  req := shared.Request{
    ID: id,
    Table: membership,
  }

  var reply bool
  
  // send message by adding to the requests
  err := server.Call("Requests.Add", req, &reply)
  if err != nil {
    fmt.Println("Message send failed", err)
  }
}

// Read incoming messages from other nodes
func readMessages(server *rpc.Client, id int, membership shared.Membership) *shared.Membership {
	//TODO
  incoming := shared.NewMembership()
  
  // receive any sent messages 
  err := server.Call("Requests.Listen", id, incoming)
  if err != nil {
    fmt.Println("Request read failed", err)
  }
  
  // received requests
  if len(incoming.Members) > 0 {
    membership = *shared.CombineTables(&membership, incoming)
  }

  return &membership
}

func calcTime() time.Time {
  return time.Now()
}

var wg = &sync.WaitGroup{}

func main() {
	rand.Seed(time.Now().UnixNano())
	Z_TIME := rand.Intn(Z_TIME_MAX - Z_TIME_MIN) + Z_TIME_MIN

	// Connect to RPC server
	server, _ := rpc.DialHTTP("tcp", "localhost:9005")

	args := os.Args[1:]

	// Get ID from command line argument
	if len(args) == 0 {
		fmt.Println("No args given")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Found Error", err)
	}

	fmt.Println("Node", id, "will fail after", Z_TIME, "seconds")

	currTime := calcTime()
	// Construct self
	self_node = shared.Node{ID: id, Hbcounter: 0, Time: currTime, Alive: true}
	var self_node_response shared.Node // Allocate space for a response to overwrite this

	// Add node with input ID
	if err := server.Call("Membership.Add", self_node, &self_node_response); err != nil {
		fmt.Println("Error:2 Membership.Add()", err)
	} else {
		fmt.Printf("Success: Node created with id= %d\n", id)
	}

	neighbors := self_node.InitializeNeighbors(id)
	fmt.Println("Neighbors:", neighbors)

	membership := shared.NewMembership()
	membership.Add(self_node, &self_node)

	sendMessage(server, neighbors[0], *membership)

	// crashTime := self_node.CrashTime()

	time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, &self_node, &membership, id) })
	time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, neighbors, &membership, id) })
	time.AfterFunc(time.Second*time.Duration(Z_TIME), func() { runAfterZ(server, id) })

	wg.Add(1)
	wg.Wait()
}

func runAfterX(server *rpc.Client, node *shared.Node, membership **shared.Membership, id int) {
	//TODO
  // stop if dead
  if !self_node.Alive {
    return 
  }

  // Incremement (still alive)
  node.Hbcounter++
  // update time after heartbeat
  node.Time = calcTime()

  // local update (requires locking)
  self_mutex.Lock()
  (**membership).Update(*node, node) 
  printMembership(**membership)
  self_mutex.Unlock()
  
  // update node in membership list
  var reply shared.Node
  err := server.Call("Membership.Update", *node, &reply)
  if err != nil {
    fmt.Println("Error updating hb", err);
  }

  // Schedule the next HB increment for the next X duration
  time.AfterFunc(time.Second*X_TIME, func() { runAfterX(server, node, membership, id) })
}

func runAfterY(server *rpc.Client, neighbors [2]int, membership **shared.Membership, id int) {
	//TODO
  // stop if dead
  if !self_node.Alive {
    return 
  }

  // look for messages sent from neighbors
  self_mutex.Lock()
  (*membership) = readMessages(server, id, **membership)
  
  // check for dead nodes
  for id, member := range (*membership).Members {
    // let time 5T be the threshold for detecting a dead node 
    if calcTime().Sub(member.Time) > 5*Y_TIME*time.Second {
      member.Alive = false
      (**membership).Members[id] = member
    }
  }
  self_mutex.Unlock()

  self_mutex.Lock()
  // send my table to my neighbors
  for _, neighbor := range neighbors {
    sendMessage(server, neighbor, **membership)
  }
  self_mutex.Unlock()

  time.AfterFunc(time.Second*Y_TIME, func() { runAfterY(server, neighbors, membership, id) })
}

func runAfterZ(server *rpc.Client, id int) {
	//TODO
  fmt.Printf("NODE %d IS NOW DEAD\n", id)
  self_node.Alive = false 
  wg.Done() 
  return
}



func printMembership(m shared.Membership){
	for _, val := range m.Members {
		status := "is Alive"
		if !val.Alive {
			status = "is Dead"
		}
		fmt.Printf("Node %d has hb %d, time %s and %s\n", 
      val.ID, 
      val.Hbcounter, 
      val.Time.Format("15:04:05"),
      status)
	}
	fmt.Println("")
}
