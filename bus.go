package bus

import (
	"reflect"
)

type Message struct {
	From  string
	Topic string
	Data  interface{}
}

type Bus struct {
	nodes []*Node
	exit  chan struct{}
}

func New() *Bus {
	return &Bus{
		nodes: make([]*Node, 0),
		exit:  make(chan struct{}, 1),
	}
}

func (b *Bus) AddNode(n *Node) {
	b.nodes = append(b.nodes, n)
}

func (b *Bus) Run() {
	n := len(b.nodes)
	cases := make([]reflect.SelectCase, n+1)
	for i, n := range b.nodes {
		//start client
		go n.Spin()
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(n.OutQ),
		}
	}

	//exit chan
	cases[n] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(b.exit),
	}

	for {
		chosen, value, ok := reflect.Select(cases)
		if chosen < n && ok {
			msg := value.Interface().(*Message)
			for _, n := range b.nodes {
				if n.IsSubscribed(msg.Topic) {
					n.InQ <- msg
				}
			}
		} else if chosen == n {
			msg := &Message{From: "bus", Topic: "exit"}
			for _, n := range b.nodes {
				n.InQ <- msg
			}
		}
	}
}

func (b *Bus) Exit() {
	b.exit <- struct{}{}
	for _, n := range b.nodes {
		n.Exit()
	}
}
