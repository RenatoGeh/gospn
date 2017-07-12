package common

type Collection interface {
	// Give inserts an element into an arbitrary position of the data type, dependant on
	// implementation.
	Give(e interface{})
	// Take returns and removes an element from the data type. What element and in which position is
	// dependant on implementation.
	Take() interface{}
	// Peek returns an element from the data type. Which of the element is to be peeked at is
	// dependant on implementation.
	Peek() interface{}
	// Get returns the i-th element of the data type.
	Get(i int) interface{}
	// Size returns the size of the data type in number of elements.
	Size() int
	// Empty returns whether the data type is empty or not.
	Empty() bool
}
