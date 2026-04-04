package domain

// NodePasso is a flat row returned from DB queries.
type NodePasso struct {
	ID        string  `json:"id"`
	Descricao string  `json:"descricao"`
	ParentID  *string `json:"parent_id"`
}
