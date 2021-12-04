package exprtree

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestNode_Expr(t *testing.T) {
	n0 := Node{
		OpType: OpTypeLogicOperation,
		Value:  AND,
		Child: []*Node{
			{
				OpType: OpTypeLogicOperation,
				Value:  OR,
				Child: []*Node{
					{
						OpType: OpTypeStrCompare,
						Value:  STREQ,
						Child: []*Node{
							{
								Value: "36",
							},
							{
								Value: "360",
							},
						},
					},

					{
						OpType: OpTypeIntCompare,
						Value:  NUMGT,
						Child: []*Node{
							{
								Value: 3600,
							},
							{
								Value: 100,
							},
						},
					},
				},
			},
			{
				OpType: OpTypeIntCompare,
				Value:  NUMGT,
				Child: []*Node{
					{
						Value: 360,
					},
					{
						Value: 0,
					},
				},
			},
		},
	}
	fmt.Println(n0.Expr())
	s, err := json.Marshal(n0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(s))
}

func TestNodeUnS(t *testing.T) {
	s := `{"op_type":1,"value":"AND","child":[{"op_type":1,"value":"OR","child":[{"op_type":3,"value":"==?","child":[{"op_type":0,"value":"36","child":null},{"op_type":0,"value":"360","child":null}]},{"op_type":2,"value":"\u003e?","child":[{"op_type":0,"value":360,"child":null},{"op_type":0,"value":100,"child":null}]}]},{"op_type":2,"value":"\u003e?","child":[{"op_type":0,"value":360,"child":null},{"op_type":0,"value":0,"child":null}]}]}`
	node := &Node{}
	json.Unmarshal([]byte(s), node)
	fmt.Println(node.Expr())
}

func TestNode_ExprWithArgs(t *testing.T) {
	n0 := Node{
		OpType: OpTypeLogicOperation,
		Value:  AND,
		Child: []*Node{
			{
				OpType: OpTypeLogicOperation,
				Value:  AND,
				Child: []*Node{
					{
						OpType: OpTypeStrCompare,
						Value:  STRCONTAINS,
						Child: []*Node{
							{
								Value: "$1",
							},
							{
								Value: "360",
							},
						},
					},
					{
						OpType: OpTypeLogicOperation,
						Value:  OR,
						Child: []*Node{
							{
								OpType: OpTypeStrCompare,
								Value:  STREQ,
								Child: []*Node{
									{
										Value: "$2",
									},
									{
										Value: "1.1.2",
									},
								},
							},
							{
								OpType: OpTypeStrCompare,
								Value:  STREQ,
								Child: []*Node{
									{
										Value: "$2",
									},
									{
										Value: "1.1.2",
									},
								},
							},
						},
					},
				},
			},
			{
				OpType: OpTypeIntCompare,
				Value:  NUMGT,
				Child: []*Node{
					{
						Value: int(360),
					},
					{
						Value: int(0),
					},
				},
			},
		},
	}
	fmt.Println(n0.Expr("360", "1.1.2"))
	s, err := json.Marshal(n0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(s))
}

func TestNewPlaceholder(t *testing.T) {
	SetDefaultPlaceholderSymbol("#$#")
	s := NewPlaceholder(10)
	fmt.Println(s)
}
