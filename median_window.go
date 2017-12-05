package main

import "fmt"
import "math/rand"
import "time"

type MedianWindow struct {
	nodes  []Node
	head   *Node
	idx    int
	winLen int
	filled bool
}

type Node struct {
	value float64
	next  *Node
	prev  *Node
}

func New(winLen int) MedianWindow {
	// Initialise a new median window with max length
	mw := MedianWindow{}
	mw.winLen = winLen
	mw.nodes = make([]Node, winLen)

	return mw
}

func (mw *MedianWindow) skipInsert() {
	new := &mw.nodes[mw.idx]
	curr_p := &mw.head

	if *curr_p == new {
		// New node is already set as head
		return
	}

	var prev *Node

	for {
		current := *curr_p
		if new.value < current.value {
			new.next = current
			new.prev = prev
			current.prev = new
			*curr_p = new
			return
		}
		if current.next == nil {
			// If you are at the end of the list add to the end.
			current.next = new
			new.prev = current
			return
		}
		prev = current
		curr_p = &current.next
	}

}

func (mw *MedianWindow) getItems(target, items int) []float64 {
	vals := []float64{}

	// Parse the first half of the list
	current := mw.head
	for i := 0; i < target+items; i++ {
		if i >= target {
			fmt.Println("NODE", current.value)
			vals = append(vals, current.value)
		}
		current = current.next
	}
	return vals
}

func (mw *MedianWindow) Median() float64 {
	// Get the length of the window
	var len int
	if mw.filled {
		// If the window is filled then used defined window lenght.
		len = mw.winLen
	} else {
		len = mw.idx
	}

	// Target index
	var target int
	target = (len - 1) / 2
	// Number of items to get
	items_cnt := ((len + 1) % 2) + 1
	fmt.Println(items_cnt)
	vals := mw.getItems(target, items_cnt)

	var median float64
	for _, val := range vals {
		median += val
	}

	return median / float64(items_cnt)

}

func (mw *MedianWindow) skipRemove() {
	// Need to have a doubly linked list.
	rm := &mw.nodes[mw.idx] // Node to remove
	next := rm.next         // Next node
	prev := rm.prev         // Prev node

	if prev == nil {
		// If first item in the list then replace the head pointer and
		// reset the prev pointer.
		mw.head = next
		mw.head.prev = nil
	}

	if prev != nil {
		prev.next = next
	}

	if next != nil {
		next.prev = prev
	}
}

func (mw *MedianWindow) push(val float64) {
	// Cannot use append here because we don't want to reallocate memory
	new := Node{value: val}

	if !mw.filled {
		// Initialse the slice on the first loop through the window
		mw.nodes[mw.idx] = new
		if mw.head == nil {
			// Initialise the first node as the head
			mw.head = &mw.nodes[mw.idx]
		}
		// insert into the skip list here
		mw.skipInsert()

		// Increment the counter
		mw.idx = mw.idx + 1
		if mw.idx == mw.winLen {
			// The window is full, reset the counter and mark window as full
			mw.idx = 0
			mw.filled = true
		}
		return
	}
	// Insert into the skip list here
	mw.skipRemove()
	mw.nodes[mw.idx] = new
	mw.skipInsert()

	// Increment circular counter
	mw.idx = (mw.idx + 1) % mw.winLen

}

func (mw *MedianWindow) printList() {
	current := mw.head
	for {
		fmt.Println("NODE", current)
		if current.next == nil {
			break
		}

		current = current.next
	}
}

func main() {
	mw := New(1000000)

	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < 10000; i++ {
		mw.push(rand.Float64())
	}
	fmt.Println("results")
	fmt.Println(mw.Median())
}
