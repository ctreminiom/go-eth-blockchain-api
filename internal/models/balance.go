package models

type BalanceResponse struct {
	Address string `json:"address,omitempty"`
	Balance string `json:"balance,omitempty"`
	Symbol  string `json:"symbol,omitempty"`
	Units   string `json:"units,omitempty"`
}
