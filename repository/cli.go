package repository

import (
	"fmt"
)

type Node struct {
	Location   string
	APIKey     string
	Name       string
	Containers string
}

func (r *SQLite) InsertNode(node Node) {
	r.db.Exec("INSERT INTO nodes (location, api_key, name, containers) VALUES ($1, $2, $3, $4)", node.Location, node.APIKey, node.Name, node.Containers)
}

func (r *SQLite) DoesNodeExist(name string) bool {
	var node Node

	err := r.db.Get(&node, "SELECT name FROM nodes WHERE name = $1", name)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
