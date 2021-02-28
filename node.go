package bus

type Node struct {
	name      string
	subtopics map[string]bool
	inQ       chan *Message
	outQ      chan *Message
	exit      chan struct{}
	onMessage func(c *Node, msg *Message)
}

func NewNode(name string, onMessage func(c *Node, msg *Message)) *Node {
	return &Node{
		name:      name,
		inQ:       make(chan *Message, 10),
		outQ:      make(chan *Message, 10),
		exit:      make(chan struct{}, 1),
		subtopics: make(map[string]bool),
		onMessage: onMessage,
	}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Subscribe(topics ...interface{}) {
	for _, v := range topics {
		if s, ok := v.(string); ok {
			n.subtopics[s] = true
		} else if s, ok := v.([]string); ok {
			for _, e := range s {
				n.subtopics[e] = true
			}
		}
	}
}

func (n *Node) UnSubscribe(topics ...interface{}) {
	for _, v := range topics {
		if s, ok := v.(string); ok {
			delete(n.subtopics, s)
		} else if s, ok := v.([]string); ok {
			for _, v := range s {
				delete(n.subtopics, v)
			}
		}
	}
}

func (n *Node) IsSubscribed(topic string) bool {
	if b, ok := n.subtopics[topic]; ok {
		return b
	}
	return false
}

func buildMsg(name, topic string, data interface{}) *Message {
	return &Message{
		From:  name,
		Topic: topic,
		Data:  data,
	}
}

func (n *Node) Publish(topic string, data interface{}) {
	msg := buildMsg(n.name, topic, data)
	n.outQ <- msg
}

func (n *Node) Spin() {
	for msg := range n.inQ {
		if msg.From == "_bus" && msg.Topic == "exit" {
			break
		}
		n.onMessage(n, msg)
	}
	n.exit <- struct{}{}
}

func (n *Node) Exit() {
	<-n.exit
}
