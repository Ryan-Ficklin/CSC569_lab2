package shared

import (
	//"fmt"
	"math/rand"
	"time"
  "sync"
)

const (
	MAX_NODES = 8
)

// Node struct represents a computing node.
type Node struct {
	ID        int
	Hbcounter int
	Time      time.Time
	Alive     bool
}

// Generate random crash time from 10-60 seconds
func (n Node) CrashTime() int {
	rand.Seed(time.Now().UnixNano())
	max := 60
	min := 10
	return rand.Intn(max-min) + min
}

func (n Node) InitializeNeighbors(id int) [2]int {
	neighbor1 := RandInt()
	for neighbor1 == id {
		neighbor1 = RandInt()
	}
	neighbor2 := RandInt()
	for neighbor1 == neighbor2 || neighbor2 == id {
		neighbor2 = RandInt()
	}
	return [2]int{neighbor1, neighbor2}
}

func RandInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(MAX_NODES-1+1) + 1
}

/*---------------*/

// Membership struct represents participanting nodes
type Membership struct {
	mutex sync.RWMutex
  Members map[int]Node
}

// Returns a new instance of a Membership (pointer).
func NewMembership() *Membership {
	return &Membership{
		Members: make(map[int]Node),
	}
}

// Adds a node to the membership list.
func (m *Membership) Add(payload Node, reply *Node) error {
	//TODO
  m.mutex.Lock()
  defer m.mutex.Unlock()
  
  m.Members[payload.ID] = payload
  return nil
}

// Updates a node in the membership list.
func (m *Membership) Update(payload Node, reply *Node) error {
	//TODO
  m.mutex.Lock()
  defer m.mutex.Unlock()
 
 m.Members[payload.ID] = payload
 return nil
}

// Returns a node with specific ID.
func (m *Membership) Get(payload int, reply *Node) error {
	//TODO
  m.mutex.Lock()
  defer m.mutex.Unlock()
  
  *reply = m.Members[payload]
  return nil
}

/*---------------*/

// Request struct represents a new message request to a client
type Request struct {
	ID    int
	Table Membership
}

// Requests struct represents pending message requests
type Requests struct {
	mutex sync.RWMutex
  Pending map[int]Membership
}

// Returns a new instance of a Membership (pointer).
func NewRequests() *Requests {
	//TODO
  return &Requests{
    Pending: make(map[int]Membership),
  }
}

// Adds a new message request to the pending list
func (req *Requests) Add(payload Request, reply *bool) error {
	//TODO
  req.mutex.Lock()
  defer req.mutex.Unlock()
  
  existing, exists := req.Pending[payload.ID]
  if exists {
    req.Pending[payload.ID] = *CombineTables(&existing, &payload.Table)
  } else {
    // no request exists already
    req.Pending[payload.ID] = payload.Table
  }

  *reply = true
  return nil
}

// Listens to communication from neighboring nodes.
func (req *Requests) Listen(ID int, reply *Membership) error {
	//TODO
  req.mutex.Lock()
  defer req.mutex.Unlock()

  membership, ok := req.Pending[ID]

  if ok{
    *reply = membership
    // now that the message has been received, remove it from the pending list
    delete(req.Pending, ID)
  }
  return nil
}

func CombineTables(table1 *Membership, table2 *Membership) *Membership {
	//TODO

  combined := NewMembership()
  
  table1.mutex.Lock()
  defer table1.mutex.Unlock()
  // add all of table 1 to the combined list
  for id, member := range table1.Members {
    combined.Members[id] = member
  }

  table2.mutex.Lock()
  defer table2.mutex.Unlock()
  // iterate through table2 and add or update members only when table2 is 
  // more recent
  for id, member := range table2.Members {
    maybeMember, exists := combined.Members[id]

    if !exists || member.Hbcounter > maybeMember.Hbcounter {
      member.Time = time.Now()
      combined.Members[id] = member 
    }
  }
  
  return combined
}

