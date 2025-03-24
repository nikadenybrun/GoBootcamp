package task00

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func AreToysBalanced(root *TreeNode) bool {
	count := CountNode(root.Left) - CountNode(root.Right)
	if count == 0 {
		return true
	}
	return false
}
func CountNode(node *TreeNode) int {
	count := 0
	if node != nil {
		if node.HasToy == true {
			count += 1
		}
		count += CountNode(node.Left)
		count += CountNode(node.Right)

	}

	return count
}

func UnrollGarland(node *TreeNode) []bool {
	slice := make([]bool, 0)
	queue := []*TreeNode{node}
	level := 1
	degree := 1
	for len(queue) > 0 {
		count := 0
		for i := 0; i < degree && len(queue) > 0; i++ {
			node := queue[0]
			queue = queue[1:]

			slice = append(slice, node.HasToy)

			if node.Left != nil {
				count++
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				count++
				queue = append(queue, node.Right)

			}

		}
		size := len(queue)
		if level%2 == 0 && size > 1 {
			for i := 0; i < count/2; i++ {
				queue[i], queue[count-1-i] = queue[count-1-i], queue[i]

			}

		}

		for i := 0; i < level; i++ {
			degree *= 2
		}
		level++
	}
	return slice

}
