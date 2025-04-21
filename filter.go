package orgtree

import (
	"regexp"
)

// FilterSubtree возвращает новое дерево, содержащее только узлы, значения которых соответствуют условию.
// Если узел не соответствует, но его потомки соответствуют — они "поднимаются".
func (n *Node) FilterSubtree(predicate func(interface{}) bool) *Node {
	var cloneIfMatched func(node *Node) *Node

	cloneIfMatched = func(node *Node) *Node {
		matchingChildren := []*Node{}
		for _, child := range node.Children {
			if filtered := cloneIfMatched(child); filtered != nil {
				matchingChildren = append(matchingChildren, filtered)
			}
		}

		if predicate(node.Value) || len(matchingChildren) > 0 {
			newNode := NewNode(node.Value)
			for _, child := range matchingChildren {
				newNode.AddChild(child)
			}
			return newNode
		}
		return nil
	}

	return cloneIfMatched(n)
}

// FilterSubtreeByRegex фильтрует дерево по регулярному выражению, применяемому к значениям типа string
func (n *Node) FilterSubtreeByRegex(pattern string) (*Node, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	predicate := func(val interface{}) bool {
		if s, ok := val.(string); ok {
			return re.MatchString(s)
		}
		return false
	}

	return n.FilterSubtree(predicate), nil
}
