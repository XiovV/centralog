package repository

type Node struct {
	Location   string
	APIKey     string `db:"api_key"`
	Name       string
	Containers string
}

func (r *SQLite) InsertNode(node Node) error {
	_, err := r.db.Exec("INSERT INTO nodes (location, api_key, name, containers) VALUES ($1, $2, $3, $4)", node.Location, node.APIKey, node.Name, node.Containers)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLite) GetNodes() ([]Node, error) {
	var nodes []Node

	if err := r.db.Select(&nodes, "SELECT location, api_key, name, containers FROM nodes"); err != nil {
		return nodes, err
	}

	return nodes, nil
}

func (r *SQLite) GetNode(name string) (Node, error) {
	var node Node

	if err := r.db.Get(&node, "SELECT location, api_key, name, containers FROM nodes WHERE name = $1", name); err != nil {
		return Node{}, err
	}

	return node, nil
}

func (r *SQLite) UpdateNodeName(oldName, newName string) error {
	_, err := r.db.Exec("UPDATE nodes SET name = $1 WHERE name = $2", newName, oldName)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLite) DoesNodeExist(name string) bool {
	var node Node

	if err := r.db.Get(&node, "SELECT name FROM nodes WHERE name = $1", name); err != nil {
		return false
	}

	return true
}
