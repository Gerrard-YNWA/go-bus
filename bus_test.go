package bus

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	b := New()

	if b == nil {
		t.Fail()
	}
}

func TestAddNode(t *testing.T) {
	b := New()
	b.AddNode(NewNode("test", onMessage))
	if len(b.nodes) != 1 {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	b := New()
	n := NewNode("test", onMessage)
	n.Subscribe("/a")
	b.AddNode(n)

	go b.Run()
	n.Publish("/a", "hello")
	time.Sleep(1 * time.Second)
	b.exit <- struct{}{}
}

func TestExit(t *testing.T) {
	b := New()
	n := NewNode("test", onMessage)
	n.Subscribe("/a")
	b.AddNode(n)

	go b.Run()
	b.Exit()
}
