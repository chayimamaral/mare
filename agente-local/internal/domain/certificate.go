package domain

type Certificate struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Subject    string `json:"subject"`
	SerialHex  string `json:"serial_hex"`
	SlotID     uint   `json:"slot_id"`
	TokenLabel string `json:"token_label"`
}
