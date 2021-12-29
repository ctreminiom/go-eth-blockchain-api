package models

type Transaction struct {
	Hash     string `json:"hash"`
	Value    string `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
	Nonce    uint64 `json:"nonce"`
	To       string `json:"to"`
	Pending  bool   `json:"pending"`
}
