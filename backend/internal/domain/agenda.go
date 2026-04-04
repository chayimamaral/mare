package domain

type AgendaEvent struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	RotinaID        string `json:"rotina_id,omitempty"`
	PassoID         string `json:"passo_id,omitempty"`
	AgendaID        string `json:"agenda_id,omitempty"`
	Start           string `json:"start"`
	End             string `json:"end"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
	BorderColor     string `json:"borderColor"`
}

type ConcluirPassoResult struct {
	AgendaID              string `json:"agenda_id"`
	AgendaItemID          string `json:"agenda_item_id"`
	TodosPassosConcluidos bool   `json:"todos_passos_concluidos"`
}
