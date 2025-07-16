package quotes

import (
	"time"

	"github.com/shopspring/decimal"
)

type RequestQuote struct {
	Shipper        RequestShipper       `json:"shipper"`
	Recipient      RequestRecipient     `json:"recipient"`
	Dispatchers    []RequestDispatchers `json:"dispatchers"`
	Channel        string               `json:"channel"`
	Filter         int                  `json:"filter"`
	Limit          int                  `json:"limit"`
	Identification string               `json:"identification"`
	Reverse        bool                 `json:"reverse"`
	SimulationType []int                `json:"simulation_type"`
	Returns        RequestReturns       `json:"returns"`
}

type RequestShipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type RequestRecipient struct {
	Type             int    `json:"type"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Country          string `json:"country"`
	Zipcode          int    `json:"zipcode"`
}

type RequestVolumes struct {
	Amount        int     `json:"amount"`
	AmountVolumes int     `json:"amount_volumes"`
	Category      string  `json:"category"`
	Sku           string  `json:"sku"`
	Tag           string  `json:"tag"`
	Description   string  `json:"description"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  int     `json:"unitary_price"`
	UnitaryWeight int     `json:"unitary_weight"`
	Consolidate   bool    `json:"consolidate"`
	Overlaid      bool    `json:"overlaid"`
	Rotate        bool    `json:"rotate"`
}

type RequestDispatchers struct {
	RegisteredNumber string           `json:"registered_number"`
	Zipcode          int              `json:"zipcode"`
	TotalPrice       int              `json:"total_price"`
	Volumes          []RequestVolumes `json:"volumes"`
}

type RequestReturns struct {
	Composition  bool `json:"composition"`
	Volumes      bool `json:"volumes"`
	AppliedRules bool `json:"applied_rules"`
}

type ResponseQuote struct {
	Dispatchers []ResponseDispatcher `json:"dispatchers"`
}

type ResponseDispatcher struct {
	ID                         string          `json:"id"`
	RequestID                  string          `json:"request_id"`
	RegisteredNumberShipper    string          `json:"registered_number_shipper"`
	RegisteredNumberDispatcher string          `json:"registered_number_dispatcher"`
	ZipcodeOrigin              int             `json:"zipcode_origin"`
	Offers                     []ResponseOffer `json:"offers"`
	TotalPrice                 decimal.Decimal `json:"total_price"`
}

type ResponseOffer struct {
	Offer                       int                  `json:"offer"`
	TableReference              string               `json:"table_reference,omitempty"`
	SimulationType              int                  `json:"simulation_type"`
	Carrier                     ResponseCarrier      `json:"carrier"`
	Service                     string               `json:"service"`
	ServiceCode                 string               `json:"service_code,omitempty"`
	ServiceDescription          string               `json:"service_description,omitempty"`
	DeliveryTime                ResponseDeliveryTime `json:"delivery_time"`
	Expiration                  time.Time            `json:"expiration"`
	CostPrice                   decimal.Decimal      `json:"cost_price"`
	FinalPrice                  decimal.Decimal      `json:"final_price"`
	Weights                     ResponseWeights      `json:"weights"`
	Correios                    *ResponseCorreios    `json:"correios,omitempty"`
	OriginalDeliveryTime        ResponseDeliveryTime `json:"original_delivery_time"`
	Identifier                  string               `json:"identifier,omitempty"`
	HomeDelivery                bool                 `json:"home_delivery"`
	CarrierOriginalDeliveryTime ResponseDeliveryTime `json:"carrier_original_delivery_time"`
	Modal                       string               `json:"modal"`
}

type ResponseCarrier struct {
	Name             string `json:"name"`
	RegisteredNumber string `json:"registered_number"`
	StateInscription string `json:"state_inscription"`
	Logo             string `json:"logo"`
	Reference        int    `json:"reference"`
	CompanyName      string `json:"company_name"`
}

type ResponseDeliveryTime struct {
	Days          int    `json:"days,omitempty"`
	Hours         int    `json:"hours,omitempty"`
	Minutes       int    `json:"minutes,omitempty"`
	EstimatedDate string `json:"estimated_date"`
}

type ResponseWeights struct {
	Real  float64 `json:"real"`
	Cubed float64 `json:"cubed,omitempty"`
	Used  float64 `json:"used,omitempty"`
}

type ResponseCorreios struct {
	DeclaredValue bool `json:"declared_value"`
}
