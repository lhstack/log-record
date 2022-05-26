package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(Md5("123456"))
}

type Tree struct {
	Id       string
	Name     string
	ParentId string
	Children []*Tree
}

func TestTree(t *testing.T) {
	tree1 := &Tree{
		Id:       "1",
		ParentId: "0",
		Children: []*Tree{},
	}
	tree2 := &Tree{
		Id:       "2",
		ParentId: "1",
		Children: []*Tree{},
	}

	tree3 := &Tree{
		Id:       "3",
		ParentId: "1",
		Children: []*Tree{},
	}

	tree4 := &Tree{
		Id:       "4",
		ParentId: "2",
		Children: []*Tree{},
	}

	tree5 := &Tree{
		Id:       "5",
		ParentId: "3",
		Children: []*Tree{},
	}
	trees := []*Tree{
		tree1, tree2, tree3, tree4, tree5,
	}
	var result []*Tree
	for i := range trees {
		if trees[i].ParentId == "0" {
			result = append(result, trees[i])
		}
		for j := range trees {
			if trees[i].Id == trees[j].ParentId {
				trees[i].Children = append(trees[i].Children, trees[j])
			}
		}
	}

	marshal, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(marshal))
}
