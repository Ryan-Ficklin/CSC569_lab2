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
  
  req.Pending[payload.ID] = payload.Table

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
  
  // lock these tables s.t. we can read from them without updates occuring
  // mid-read
  table1.mutex.Lock()
  defer table1.mutex.Unlock()
  table2.mutex.Lock()
  defer table2.mutex.Unlock()
  
  // add all of table 1 to the combined list
  for id, member := range table1.Members {
    combined.Members[id] = member
  }

  // iterate through table2 and add or update members only when table2 is 
  // more recent
  timeout := 3*time.Second
  for id, member := range table2.Members {
    maybeMember, exists := combined.Members[id]
    // add/update with table2's member if there isnt a corresponding node in
    // table1 or if table2 has a larger heartbeat
    if !exists || member.Hbcounter > maybeMember.Hbcounter {
      member.Time = time.Now()
      combined.Members[id] = member
    // finding dead nodes:
    // 1) member exists in both
    // 2) Heartbeat has not been updated
    // 3) time since the member has been updated in both tables is longer than 
    //    some tuned arbitrary timeout (I chose 3 seconds after some testing)
    } else if exists && 
              member.Hbcounter == maybeMember.Hbcounter && 
              (time.Since(maybeMember.Time) > timeout && time.Since(member.Time) > timeout) {
              //(member.Time.Sub(maybeMember.Time).Abs() > 3*time.Second) { 
      member.Alive = false
      combined.Members[id] = member
    }
  }
  
  return combined
}

