package models

type AssetInfo struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	ApiID  string `json:"-"` //coingecko
}
