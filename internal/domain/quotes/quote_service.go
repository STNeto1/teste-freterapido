package quotes

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
	"github.com/stneto1/teste-freterapido/internal/utils"
)

type QuoteService struct {
	config                *system.QuotesServiceConfig
	freteRapidoRepository FreteRapidoQuotesRepository
	clickhouseRepository  ClickhouseQuotesRepository
}

func NewQuoteService(
	config *system.QuotesServiceConfig,
	freteRapidoRepository FreteRapidoQuotesRepository,
	clickhouseRepository ClickhouseQuotesRepository,
) *QuoteService {
	return &QuoteService{
		config:                config,
		freteRapidoRepository: freteRapidoRepository,
		clickhouseRepository:  clickhouseRepository,
	}
}

func (s *QuoteService) CreateRequestPayload(requestQuote *RequestQuote) FreteRapidoRequestQuote {
	return FreteRapidoRequestQuote{
		Shipper: FreteRapidoRequestShipper{
			RegisteredNumber: s.config.RegisteredNumber,
			Token:            s.config.Token,
			PlatformCode:     s.config.PlatformCode,
		},
		Recipient: FreteRapidoRequestRecipient{
			Type:    0,
			Country: "BRA",
			Zipcode: requestQuote.MustParseRecipientZipcode(),
		},
		Dispatchers: []FreteRapidoRequestDispatchers{
			{
				RegisteredNumber: s.config.RegisteredNumber,
				Zipcode:          s.config.DispatcherZipCode,
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
	errorSet := req.ErrorSet()
	if errorSet != nil {
		return nil, QuoteRequestErrorSetError{
			Errors: errorSet,
		}
	}

	// TODO: make retries
	result, err := s.freteRapidoRepository.TryQuote(ctx, s.CreateRequestPayload(req))
	if err != nil {
		return nil, QuoteRequestError{
			Message: err.Error(),
		}
	}

	var quotes []Quote
	for _, dispatcher := range result.Dispatchers {
		for _, offer := range dispatcher.Offers {
			quotes = append(quotes, Quote{
				ID:        uuid.Must(uuid.NewV7()),
				Name:      offer.Carrier.Name,
				Service:   offer.Service,
				Deadline:  uint8(offer.DeliveryTime.Days),
				Price:     offer.FinalPrice,
				CreatedAt: time.Now(),
			})
		}
	}

	// TODO: maybe handle no quote on response?
	go s.ProcessQuotes(ctx, quotes)

	return quotes, nil
}

// ProcessQuotes saves quotes to clickhouse, doesn't return anything because it's async
func (s *QuoteService) ProcessQuotes(ctx context.Context, quotes []Quote) {
	slog.Info("processing quotes")

	if len(quotes) == 0 {
		slog.Error("no quotes to save")
		return
	}

	// TODO: handle retries
	if err := s.clickhouseRepository.AddQuotes(ctx, quotes); err != nil {
		slog.Error("failed to save quotes to clickhouse",
			slog.Any("error", err),
		)
		return
	}
}
