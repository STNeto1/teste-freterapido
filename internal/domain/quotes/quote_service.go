package quotes

import (
	"context"
	"fmt"

	"github.com/stneto1/teste-freterapido/internal/utils"
)

type QuoteService struct {
	freteRapidoRepository FreteRapidoQuotesRepository
}

func NewQuoteService(freteRapidoRepository FreteRapidoQuotesRepository) *QuoteService {
	return &QuoteService{
		freteRapidoRepository: freteRapidoRepository,
	}
}

func (s *QuoteService) CreateRequestPayload(requestQuote *RequestQuote) FreteRapidoRequestQuote {
	return FreteRapidoRequestQuote{
		Shipper: FreteRapidoRequestShipper{
			RegisteredNumber: "SUT",
			Token:            "SUT",
			PlatformCode:     "SUT",
		},
		Recipient: FreteRapidoRequestRecipient{
			Type:             0,
			RegisteredNumber: "SUT",
			StateInscription: "SUT",
			Country:          "BRA",
			// Zipcode:          requestQuote.Recipient.Address.Zipcode,
		},
		Dispatchers: []FreteRapidoRequestDispatchers{
			{
				RegisteredNumber: "SUT",
				Zipcode:          0, // TODO: use option from config
				// TotalPrice:       0, // TODO: reduce from volumes?
				Volumes: utils.Map(requestQuote.Volumes, func(vol RequestQuoteVolume) FreteRapidoRequestVolumes {
					return FreteRapidoRequestVolumes{
						Amount: vol.Amount,
						// AmountVolumes: 0,
						Category:      fmt.Sprintf("%d", vol.Category),
						Sku:           vol.Sku,
						Tag:           "",
						Description:   "",
						Height:        vol.Height,
						Width:         vol.Width,
						Length:        vol.Length,
						UnitaryPrice:  vol.Price.InexactFloat64(),
						UnitaryWeight: vol.UnitaryWeight,
						Consolidate:   false,
						Overlaid:      false,
						Rotate:        false,
					}
				}),
			},
		},
		Channel:        "",
		Filter:         0,
		Limit:          0,
		Identification: "",
		Reverse:        false,
		SimulationType: []int{0},
		Returns:        FreteRapidoRequestReturns{},
	}
}

func (s *QuoteService) GetFreteRapidoQuotes(ctx context.Context, req *RequestQuote) ([]Quote, error) {
	if _, err := req.ParseRecipientZipcode(); err != nil {
		return nil, QuoteInvalidZipcodeError{
			Message: "Invalid recipient zipcode",
		}
	}

	if invalidIdx := req.ValidateCategories(); invalidIdx != -1 {
		return nil, QuoteInvalidCategoryError{
			Message: fmt.Sprintf("Category on volume %d is invalid", invalidIdx+1),
		}
	}

	if invalidIdx := req.ValidateAmount(); invalidIdx != -1 {
		return nil, QuoteInvalidAmountError{
			Message: fmt.Sprintf("Amount on volume %d is invalid", invalidIdx+1),
		}
	}

	if invalidIdx := req.ValidateDimensions(); invalidIdx != -1 {
		return nil, QuoteInvalidDimensionError{
			Message: fmt.Sprintf("Dimensions on volume %d are invalid", invalidIdx+1),
		}
	}

	if invalidIdx := req.ValidatePrice(); invalidIdx != -1 {
		return nil, QuoteInvalidPriceError{
			Message: fmt.Sprintf("Price on volume %d is invalid", invalidIdx+1),
		}
	}

	if invalidIdx := req.ValidateWeight(); invalidIdx != -1 {
		return nil, QuoteInvalidWeightError{
			Message: fmt.Sprintf("Weight on volume %d is invalid", invalidIdx+1),
		}
	}

	return nil, nil
}
