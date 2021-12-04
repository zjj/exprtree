package exprtree

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NodeOpType int

const (
	// OpType
	OpTypeNull           NodeOpType = 0
	OpTypeLogicOperation NodeOpType = 1 // and, or ...
	OpTypeIntCompare     NodeOpType = 2
	OpTypeStrCompare     NodeOpType = 3

	// NUM OP
	NUMGTE = ">=?"
	NUMGT  = ">?"
	NUMEQ  = "==?"
	NUMLT  = "<?"
	NUMLTE = "<=?"

	// STRING OP
	STREQ       = "==?"
	STRCONTAINS = "c?"

	// LOGIC OP
	AND = "AND"
	OR  = "OR"
)

func intCompare(op string, ai, bi int) (bool, error) {
	switch op {
	case NUMGTE:
		return ai >= bi, nil
	case NUMGT:
		return ai > bi, nil
	case NUMLTE:
		return ai <= bi, nil
	case NUMLT:
		return ai < bi, nil
	case NUMEQ:
		return ai == bi, nil
	default:
		return false, fmt.Errorf("unsupported op %s for intCompare", op)
	}
}

func strCompare(op string, s0, s1 string) (bool, error) {
	switch op {
	case STREQ:
		return s0 == s1, nil
	case STRCONTAINS:
		return strings.Contains(s0, s1), nil
	default:
		return false, fmt.Errorf("unsupported op %s for strCompare", op)
	}
}

type Node struct {
	OpType NodeOpType  `json:"op_type"`
	Value  interface{} `json:"value"`
	Child  []*Node     `json:"child"`
}

func New() *Node {
	return &Node{
		Child: make([]*Node, 0),
	}
}

func (node *Node) AddChild(c *Node) {
	node.Child = append(node.Child, c)
}

func (node *Node) SetOpType(t NodeOpType) {
	node.OpType = t
}

func (node *Node) SetValue(v interface{}) {
	node.Value = v
}

func (node Node) Serialize() ([]byte, error) {
	return json.Marshal(node)
}

func (node Node) doIntCompareWithArgs(args ...interface{}) (ret bool, err error) {
	if len(node.Child) != 2 {
		err = errors.New("node shall have two args")
		return
	}

	var ai, bi int
	var a, b interface{}
	a, err = node.Child[0].exprWithArgs(args...)
	if err != nil {
		return
	}

	switch v := a.(type) {
	case int:
		ai = v
	case float64:
		ai = int(v)
	default:
		err = fmt.Errorf("can't convert to num :%T", v)
		return
	}

	b, err = node.Child[1].exprWithArgs(args...)
	if err != nil {
		return
	}

	switch v := b.(type) {
	case int:
		bi = v
	case float64:
		bi = int(v)
	default:
		err = fmt.Errorf("can't convert to num :%T", v)
		return
	}

	op := node.Value.(string)
	return intCompare(op, int(ai), int(bi))
}

func (node Node) doStringCompareWithArgs(args ...interface{}) (ret bool, err error) {
	if len(node.Child) != 2 {
		err = errors.New("STRINGEQUAL? shall have two args")
		return
	}

	var a, b interface{}
	a, err = node.Child[0].exprWithArgs(args...)
	if err != nil {
		return
	}
	as, ok := a.(string)
	if !ok {
		err = fmt.Errorf("wrong type %s", a)
		return
	}
	b, err = node.Child[1].exprWithArgs(args...)
	if err != nil {
		return
	}
	bs, ok := b.(string)
	if !ok {
		err = fmt.Errorf("wrong type %s", b)
		return
	}

	op := node.Value.(string)
	return strCompare(op, as, bs)
}

func (node Node) doLogicOpWithArgs(args ...interface{}) (ret bool, err error) {
	op := node.Value.(string)
	switch op {
	case OR:
		for _, subNode := range node.Child {
			var r interface{}
			r, err = subNode.exprWithArgs(args...)
			if err != nil {
				return
			}
			if b, ok := r.(bool); ok {
				if b {
					return b, nil
				}
			}
		}
		return false, nil
	case AND:
		for _, subNode := range node.Child {
			var r interface{}
			r, err = subNode.exprWithArgs(args...)
			if err != nil {
				return
			}
			if b, ok := r.(bool); ok {
				if !b {
					return false, err
				}
			}
		}
		return true, nil
	default:
	}
	err = fmt.Errorf("unsupported op %s", op)
	return
}

func (node Node) isPlaceholder(s interface{}) (bool, int, error) {
	if arg, ok := s.(string); ok {
		if strings.HasPrefix(arg, placeholderSymbol) {
			pos, err := strconv.Atoi(arg[placeholderSymbolLen:])
			if err != nil {
				return true, pos, err
			}
			return true, pos, nil
		}
	}
	return false, 0, nil
}

// exprt with args
func (node Node) exprWithArgs(args ...interface{}) (interface{}, error) {
	if node.OpType == OpTypeNull {
		if len(args) > 0 {
			placehoder, pos, err := node.isPlaceholder(node.Value)
			if err == nil {
				if placehoder {
					if len(args) < pos {
						return nil, errors.New("input args too short")
					}
					return args[pos-1], nil
				} else {
					return node.Value, nil
				}
			}
		} else {
			return node.Value, nil
		}
	}

	switch node.OpType {
	case OpTypeIntCompare:
		return node.doIntCompareWithArgs(args...)
	case OpTypeStrCompare:
		return node.doStringCompareWithArgs(args...)
	case OpTypeLogicOperation:
		return node.doLogicOpWithArgs(args...)
	default:
	}
	return nil, fmt.Errorf("non supported op %d", node.OpType)
}

func (node Node) Expr(args ...interface{}) (interface{}, error) {
	return node.exprWithArgs(args...)
}

var (
	placeholderSymbol    = "$"
	placeholderSymbolLen = 1
)

func NewPlaceholder(pos int) string {
	return fmt.Sprintf("%s%d", placeholderSymbol, pos)
}

func SetDefaultPlaceholderSymbol(s string) {
	placeholderSymbol = s
	placeholderSymbolLen = len(s)
}
