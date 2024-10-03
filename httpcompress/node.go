package httpcompress

import "net/http"

type Chain struct {
	head  INode
	local map[string]any
}

func (chain *Chain) Init(nodes ...INode) {
	if len(nodes) == 0 {
		return
	}
	chain.head = nodes[0]

	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}
	chain.local = make(map[string]any)
}
func (chain *Chain) GetLocal(key string) any {
	return chain.local[key]
}
func (chain *Chain) PutLocal(key string, value any) {
	chain.local[key] = value
}

type INode interface {
	Process(w http.ResponseWriter, r *http.Request)
	Next(w http.ResponseWriter, r *http.Request)
	SetNext(next INode)
}

type Node struct {
	INode
	nextNode INode
	chain    *Chain
}

func (node *Node) Next(w http.ResponseWriter, r *http.Request) {
	node.nextNode.Process(w, r)
}

func (node *Node) SetNext(next INode) {
	node.nextNode = next
}

func (node *Node) BindChain(chain *Chain) {
	node.chain = chain
}

func (node *Node) GetChain() *Chain {
	return node.chain
}

type FinalNode func(w http.ResponseWriter, r *http.Request)

func (node FinalNode) Next(w http.ResponseWriter, r *http.Request) {
	panic("FinalNode do not have Next")
}

func (node FinalNode) SetNext(next INode) {
	panic("FinalNode do not have SetNext")
}

func (node FinalNode) Process(w http.ResponseWriter, r *http.Request) {
	node(w, r)
}
