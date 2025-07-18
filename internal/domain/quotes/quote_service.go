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

	var result FreteRapidoResponseQuote
	var err error
	for attempt := range s.config.TryQuotesRetries {
		result, err = s.freteRapidoRepository.TryQuote(ctx, s.CreateRequestPayload(req))
		if err == nil {
			break
		}

		// exponential backoff -> 100ms, 200ms, 400ms, .... LIMIT ms
		time.Sleep(s.config.TryQuotesTimeout * (1 << attempt))
	}
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

	go s.ProcessQuotes(context.Background(), quotes)

	return quotes, nil
}

// ProcessQuotes saves quotes to clickhouse, doesn't return anything because it's async
func (s *QuoteService) ProcessQuotes(ctx context.Context, quotes []Quote) {
	s.config.Logger.Info("processing quotes")

	if len(quotes) == 0 {
		s.config.Logger.Error("no quotes to save")
		return
	}

	var err error
	for attempt := range s.config.AddQuotesRetries {
		err = s.clickhouseRepository.AddQuotes(ctx, quotes)
		if err == nil {
			break
		}
		// exponential backoff -> 100ms, 200ms, 400ms, .... LIMIT ms
		time.Sleep(s.config.TryQuotesTimeout * (1 << attempt))
	}
	if err != nil {
		s.config.Logger.Error("failed to save quotes to clickhouse",
			slog.Any("error", err),
		)
	}
}
