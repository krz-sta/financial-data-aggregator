package models

type AssetInfo struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	ApiID  string `json:"-"` //coingecko
}

type FrankfurterAsset struct {
	Base  string  `json:"base"`
	Quote string  `json:"quote"`
	Rate  float64 `json:"rate"`
}
