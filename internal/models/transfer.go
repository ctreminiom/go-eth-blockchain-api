package models

type TransferEthRequest struct {
	PrivKey string `json:"privKey,omitempty"`
	To      string `json:"to,omitempty"`
	Amount  int64  `json:"amount,omitempty"`
}

type HashResponse struct {
	Hash string `json:"hash,omitempty"`
}
