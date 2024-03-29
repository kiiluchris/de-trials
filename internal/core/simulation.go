package core

type SimulationState[T any] struct {
	TxID  string `json:"tx"`
	Query string `json:"query"`
	Rows  []T    `json:"rows"`
}
