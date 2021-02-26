# go-bus
message bus for pub/sub

```
func onMessage(n *bus.Node, msg *bus.Message) {
    fmt.Printf("[%s] recv message from [%s] topic [%s] data [%#v] \n", n.Name(), msg.From, msg.Topic, msg.Data)
}

n1 := bus.NewNode("a", onMessage)
n2 := bus.NewNode("b", onMessage)

n1.Subscribe("/a", "/a/b")
n2.Subscribe("/b", "/a/b")

b := bus.New()
b.AddNode(n1)
b.AddNode(n2)

go b.Run()

n1.Publish("/a/b", []string{"bbbb", "sss"})
n1.Publish("/b", "hello")
```
