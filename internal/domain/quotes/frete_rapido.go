package quotes

import (
	"time"

	"github.com/shopspring/decimal"
)

type FreteRapidoRequestQuote struct {
	Shipper        FreteRapidoRequestShipper       `json:"shipper"`
	Recipient      FreteRapidoRequestRecipient     `json:"recipient"`
	Dispatchers    []FreteRapidoRequestDispatchers `json:"dispatchers"`
	Channel        string                          `json:"channel"`
	Filter         int                             `json:"filter"`
	Limit          int                             `json:"limit"`
	Identification string                          `json:"identification"`
	Reverse        bool                            `json:"reverse"`
	SimulationType []int                           `json:"simulation_type"`
	Returns        FreteRapidoRequestReturns       `json:"returns"`
}

type FreteRapidoRequestShipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type FreteRapidoRequestRecipient struct {
	Type             int    `json:"type"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Country          string `json:"country"`
	Zipcode          int64  `json:"zipcode"`
}

type FreteRapidoRequestVolumes struct {
	Amount int `json:"amount"`
	// AmountVolumes int     `json:"amount_volumes"`
	Category      string  `json:"category"`
	Sku           string  `json:"sku"`
	Tag           string  `json:"tag"`
	Description   string  `json:"description"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Consolidate   bool    `json:"consolidate"`
	Overlaid      bool    `json:"overlaid"`
	Rotate        bool    `json:"rotate"`
}

type FreteRapidoRequestDispatchers struct {
	RegisteredNumber string `json:"registered_number"`
	Zipcode          int64  `json:"zipcode"`
	// TotalPrice       int                         `json:"total_price"`
	Volumes []FreteRapidoRequestVolumes `json:"volumes"`
}

type FreteRapidoRequestReturns struct {
	Composition  bool `json:"composition"`
	Volumes      bool `json:"volumes"`
	AppliedRules bool `json:"applied_rules"`
}

type FreteRapidoResponseQuote struct {
	Dispatchers []FreteRapidoResponseDispatcher `json:"dispatchers"`
}

type FreteRapidoResponseDispatcher struct {
	ID                         string                     `json:"id"`
	FreteRapidoRequestID       string                     `json:"request_id"`
	RegisteredNumberShipper    string                     `json:"registered_number_shipper"`
	RegisteredNumberDispatcher string                     `json:"registered_number_dispatcher"`
	ZipcodeOrigin              int                        `json:"zipcode_origin"`
	Offers                     []FreteRapidoResponseOffer `json:"offers"`
	TotalPrice                 decimal.Decimal            `json:"total_price"`
}

type FreteRapidoResponseOffer struct {
	Offer                       int                             `json:"offer"`
	TableReference              string                          `json:"table_reference,omitempty"`
	SimulationType              int                             `json:"simulation_type"`
	Carrier                     FreteRapidoResponseCarrier      `json:"carrier"`
	Service                     string                          `json:"service"`
	ServiceCode                 string                          `json:"service_code,omitempty"`
	ServiceDescription          string                          `json:"service_description,omitempty"`
	DeliveryTime                FreteRapidoResponseDeliveryTime `json:"delivery_time"`
	Expiration                  time.Time                       `json:"expiration"`
	CostPrice                   decimal.Decimal                 `json:"cost_price"`
	FinalPrice                  decimal.Decimal                 `json:"final_price"`
	Weights                     FreteRapidoResponseWeights      `json:"weights"`
	Correios                    *FreteRapidoResponseCorreios    `json:"correios,omitempty"`
	OriginalDeliveryTime        FreteRapidoResponseDeliveryTime `json:"original_delivery_time"`
	Identifier                  string                          `json:"identifier,omitempty"`
	HomeDelivery                bool                            `json:"home_delivery"`
	CarrierOriginalDeliveryTime FreteRapidoResponseDeliveryTime `json:"carrier_original_delivery_time"`
	Modal                       string                          `json:"modal"`
}

type FreteRapidoResponseCarrier struct {
	Name             string `json:"name"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Logo             string `json:"logo"`
	Reference        int    `json:"reference"`
	CompanyName      string `json:"company_name"`
}

type FreteRapidoResponseDeliveryTime struct {
	Days          int    `json:"days,omitempty"`
	Hours         int    `json:"hours,omitempty"`
	Minutes       int    `json:"minutes,omitempty"`
	EstimatedDate string `json:"estimated_date"`
}

type FreteRapidoResponseWeights struct {
	Real  float64 `json:"real"`
	Cubed float64 `json:"cubed,omitempty"`
	Used  float64 `json:"used,omitempty"`
}

type FreteRapidoResponseCorreios struct {
	DeclaredValue bool `json:"declared_value"`
}
