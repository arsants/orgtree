package orgtree

import "github.com/google/uuid"

// NodeType представляет тип узла в оргструктуре
type NodeType struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SysName string    `json:"sysname"`
}

// OrgNode представляет узел в оргструктуре
type OrgNode struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	SysName  string    `json:"sysname"`
	Position *Position `json:"position,omitempty"`
	Type     *NodeType `json:"type,omitempty"`
}

// EdgeType представляет тип связи между узлами
type EdgeType struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SysName string    `json:"sysname"`
}

// Edge представляет связь между узлами
type Edge struct {
	Type     *EdgeType `json:"type,omitempty"`
	FromNode uuid.UUID `json:"from_node"`
	ToNode   uuid.UUID `json:"to_node"`
}

// Position представляет должность
type Position struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SysName string    `json:"sysname"`
}
