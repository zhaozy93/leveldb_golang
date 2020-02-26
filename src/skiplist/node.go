package skiplist

type Node struct {
	key  interface{}   
	next []*Node
}

func newNode(key interface{}, height int) (x *Node) {
	x = &Node{}
	x.key = key
	x.next = make([]*Node, height+1)
	return
}
func (node *Node) getNext(level int) (n *Node) {
	if level < len(node.next){
		return node.next[level]
	}
	return 
}

func (node *Node) setNext(level int, x *Node) {
	if level < len(node.next){
		node.next[level] = x
	}
}
