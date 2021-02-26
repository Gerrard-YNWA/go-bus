package bus

import (
	"testing"
)

func onMessage(n *Node, msg *Message) {}

func TestNewNode(t *testing.T) {
	n := NewNode("test", onMessage)
	if n == nil {
		t.Fail()
	}
}

func TestNodeName(t *testing.T) {
	cases := []map[string]string{
		{
			"name":   "a",
			"expect": "node-a",
		},
		{
			"name":   "b",
			"expect": "node-b",
		},
	}

	for _, v := range cases {
		c := NewNode(v["name"], onMessage)
		if c.Name() != v["expect"] {
			t.Fail()
		}
	}
}

func TestSubscribe(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"topics": "/a",
			"expect": 1,
		},
		{
			"topics": []string{"/a/b", "/a/c"},
			"expect": 2,
		},
	}

	for _, v := range cases {
		n := NewNode("test", onMessage)
		n.Subscribe(v["topics"])
		if len(n.subtopics) != v["expect"] {
			t.Fail()
		}
	}
}

func TestIsSubscribed(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"topic":  "/a",
			"expect": true,
		},
		{
			"topic":  "/ab",
			"expect": false,
		},
	}

	n := NewNode("test", onMessage)
	n.Subscribe("/a", "/b")
	for _, v := range cases {
		if n.IsSubscribed(v["topic"].(string)) != v["expect"] {
			t.Fail()
		}
	}
}

func TestUnSubscribe(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"topics": "/a",
			"expect": false,
		},
		{
			"topics": []string{"/a/b", "/a/c"},
			"expect": []bool{false, false},
		},
	}

	n := NewNode("test", onMessage)
	n.Subscribe("/a", "/a/b")
	for _, v := range cases {
		n.UnSubscribe(v["topics"])
		if s, ok := v["topics"].(string); ok {
			if n.IsSubscribed(s) != v["expect"] {
				t.Fail()
			}
		} else if s, ok := v["topics"].([]string); ok {
			for i, e := range s {
				if n.IsSubscribed(e) != v["expect"].([]bool)[i] {
					t.Fail()
				}
			}
		}
	}
}

func TestPublish(t *testing.T) {
	n := NewNode("test", onMessage)
	topic := "/a"
	data := "sss"
	n.Publish(topic, data)
	msg := <-n.OutQ
	if msg.From != n.Name() || msg.Topic != topic || data != msg.Data {
		t.Fail()
	}
}

func TestSpin(t *testing.T) {
	n := NewNode("test", onMessage)
	cases := []map[string]interface{}{
		{
			"from":  "a",
			"topic": "/a",
			"data":  "sss",
		},
		{
			"from":  "bus",
			"topic": "exit",
			"data":  "",
		},
	}

	go func() {
		for _, v := range cases {
			msg := buildMsg(v["from"].(string), v["topic"].(string), v["data"])
			n.InQ <- msg
		}
	}()
	n.Spin()
}
