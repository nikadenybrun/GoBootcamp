package main

import (
	"container/heap"
	"fmt"
	t "src/task00"
	p "src/task02"
)

func returnTree() *t.TreeNode {
	t := &t.TreeNode{
		HasToy: true,
		Left: &t.TreeNode{
			HasToy: true,
			Left: &t.TreeNode{
				HasToy: true,
			},
		},

		Right: &t.TreeNode{
			HasToy: false,
			Left: &t.TreeNode{
				HasToy: true,
			},
			Right: &t.TreeNode{
				HasToy: true,
			},
		},
	}
	return t
}

func test00() {
	tree := returnTree()
	fmt.Println(t.AreToysBalanced(tree))
}

func test01() {
	tree := returnTree()
	fmt.Println(t.UnrollGarland(tree))
}
func test02() {
	h := &p.Presents{{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2}}

	heap.Init(h)
	fmt.Println(h.GetNCoolestPresents(4))
}
func test03() {
	h := &[]p.Present{{Value: 5, Size: 1},
		{Value: 4, Size: 5},
		{Value: 3, Size: 1},
		{Value: 5, Size: 2}}
	capacity := 8

	selected := p.Knapsack(*h, capacity)
	fmt.Println(selected)
}

func main() {
	test00()
	test01()
	test02()
	test03()
}
