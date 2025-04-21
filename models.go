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
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SysName string    `json:"sysname"`
	TypeID  uuid.UUID `json:"type_id"`
	Type    *NodeType `json:"type,omitempty"`
}

// EdgeType представляет тип связи между узлами
type EdgeType struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	SysName string    `json:"sysname"`
}

// Edge представляет связь между узлами
type Edge struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	SysName  string    `json:"sysname"`
	TypeID   uuid.UUID `json:"type_id"`
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

// PositionNodeRelation представляет связь между должностью и узлом
type PositionNodeRelation struct {
	ID         uuid.UUID `json:"id"`
	NodeID     uuid.UUID `json:"node_id"`
	PositionID uuid.UUID `json:"position_id"`
	Position   *Position `json:"position,omitempty"`
}
