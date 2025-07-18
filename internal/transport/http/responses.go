package http

import "encoding/json"

type ResponseQuotesItem struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline uint8   `json:"deadline"`
	Price    float64 `json:"price"`
}

type ResponseQuotes struct {
	Carrier []ResponseQuotesItem `json:"carrier"`
}

type ResponseQuoteErrorSet struct {
	Errors []string `json:"errors"`
}

func (r ResponseQuoteErrorSet) ToJSON() []byte {
	rBytes, err := json.Marshal(r)
	if err != nil {
		// TODO: handle error
		return nil
	}

	return rBytes
}

type ResponseGenericError struct {
	Message string `json:"message"`
}

func (r ResponseGenericError) ToJSON() []byte {
	rBytes, err := json.Marshal(r)
	if err != nil {
		// TODO: handle error
		return nil
	}

	return rBytes
}
