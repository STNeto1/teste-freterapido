package quotes_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/internal/domain/system"
	"github.com/stneto1/teste-freterapido/mocks/quotesmocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQuoteService_CreateRequestPayload(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)
	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	result := svc.CreateRequestPayload(
		&quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{
				Address: quotes.RequestQuoteRecipientAddress{
					Zipcode: "123",
				},
			},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      1,
					Amount:        1,
					UnitaryWeight: 10.1,
					Price:         decimal.NewFromFloat(123),
					Sku:           "SKU-SUT",
					Height:        10,
					Width:         10,
					Length:        10,
				},
			},
		},
	)

	expected := quotes.FreteRapidoRequestQuote{
		Shipper: quotes.FreteRapidoRequestShipper{
			RegisteredNumber: "SUT",
			Token:            "SUT",
			PlatformCode:     "SUT",
		},
		Recipient: quotes.FreteRapidoRequestRecipient{
			Type: 0,
			// RegisteredNumber: "SUT",
			// StateInscription: "SUT",
			Country: "BRA",
			Zipcode: 123,
		},
		Dispatchers: []quotes.FreteRapidoRequestDispatchers{
			{
				RegisteredNumber: "SUT",
				Zipcode:          0,
				Volumes: []quotes.FreteRapidoRequestVolumes{
					{
						Amount: 1,
						// AmountVolumes: 0,
						Category:      "1",
						Sku:           "SKU-SUT",
						Tag:           "",
						Description:   "",
						Height:        10,
						Width:         10,
						Length:        10,
						UnitaryPrice:  123,
						UnitaryWeight: 10.1,
						Consolidate:   false,
						Overlaid:      false,
						Rotate:        false,
					},
				},
			},
		},
		Channel:        "",
		Filter:         0,
		Limit:          0,
		Identification: "",
		Reverse:        false,
		SimulationType: []int{0},
		Returns:        quotes.FreteRapidoRequestReturns{},
	}

	assert.Equal(t, expected, result)
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidPayload(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)
	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      0,
				Amount:        0,
				Width:         0,
				Height:        0,
				Length:        0,
				Price:         decimal.NewFromFloat(0),
				UnitaryWeight: 0,
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteRequestErrorSetError{
		Errors: []string{
			"Invalid recipient zipcode",
			"Category on volume 1 is invalid",
			"Dimensions on volume 1 are invalid",
			"Price on volume 1 is invalid",
			"Weight on volume 1 is invalid",
			"Amount on volume 1 is invalid",
		},
	})
}

func TestQuoteService_GetFreteRapidoQuotes_SuccessfulValidation(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)
	freteRapidoMock.
		EXPECT().
		TryQuote(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ quotes.FreteRapidoRequestQuote) (quotes.FreteRapidoResponseQuote, error) {
			return quotes.FreteRapidoResponseQuote{}, nil
		}).AnyTimes()

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	clickhouseMock.
		EXPECT().
		AddQuotes(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []quotes.Quote) error {
			return nil
		}).AnyTimes()

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      1,
				Amount:        1,
				Width:         1,
				Height:        1,
				Length:        1,
				Price:         decimal.NewFromFloat(1),
				UnitaryWeight: 1,
			},
		},
	}

	_, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	// assert.Nil(t, result)
	assert.NoError(t, err)
	// assert.Equal(t, err, quotes.QuoteInvalidWeightError{
	// 	Message: "Weight on volume 1 is invalid",
	// })
}

func TestQuoteService_GetFreteRapidoQuotes_TryQuotesFailure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)
	freteRapidoMock.
		EXPECT().
		TryQuote(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ quotes.FreteRapidoRequestQuote) (quotes.FreteRapidoResponseQuote, error) {
			return quotes.FreteRapidoResponseQuote{}, fmt.Errorf("error")
		}).AnyTimes()

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
		TryQuotesRetries:  1,
		TryQuotesTimeout:  time.Millisecond,
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      1,
				Amount:        1,
				Width:         1,
				Height:        1,
				Length:        1,
				Price:         decimal.NewFromFloat(1),
				UnitaryWeight: 1,
			},
		},
	}

	resultQuotes, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, resultQuotes)
	assert.ErrorIs(t, err, quotes.QuoteRequestError{
		Message: "error",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_TryQuotesValidReturn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)
	freteRapidoMock.
		EXPECT().
		TryQuote(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ quotes.FreteRapidoRequestQuote) (quotes.FreteRapidoResponseQuote, error) {
			return quotes.FreteRapidoResponseQuote{
				Dispatchers: []quotes.FreteRapidoResponseDispatcher{
					{
						ID:                         "SUT",
						FreteRapidoRequestID:       "SUT",
						RegisteredNumberShipper:    "SUT",
						RegisteredNumberDispatcher: "SUT",
						ZipcodeOrigin:              0,
						Offers: []quotes.FreteRapidoResponseOffer{
							{
								Offer:          1,
								TableReference: "SUT",
								SimulationType: 0,
								Carrier: quotes.FreteRapidoResponseCarrier{
									Name:             "SUT",
									RegisteredNumber: "SUT",
									StateInscription: "SUT",
									Logo:             "SUT",
									Reference:        0,
									CompanyName:      "SUT",
								},
								Service:            "SUT",
								ServiceCode:        "SUT",
								ServiceDescription: "SUT",
								DeliveryTime: quotes.FreteRapidoResponseDeliveryTime{
									Days:          1,
									Hours:         1,
									Minutes:       1,
									EstimatedDate: "SUT",
								},
								Expiration:                  time.Now(),
								CostPrice:                   decimal.NewFromFloat(10),
								FinalPrice:                  decimal.NewFromFloat(10),
								Weights:                     quotes.FreteRapidoResponseWeights{},
								Correios:                    &quotes.FreteRapidoResponseCorreios{},
								OriginalDeliveryTime:        quotes.FreteRapidoResponseDeliveryTime{},
								Identifier:                  "SUT",
								HomeDelivery:                false,
								CarrierOriginalDeliveryTime: quotes.FreteRapidoResponseDeliveryTime{},
								Modal:                       "SUT",
							},
						},
						TotalPrice: decimal.NewFromFloat(10),
					},
				},
			}, nil
		}).AnyTimes()

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	clickhouseMock.
		EXPECT().
		AddQuotes(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []quotes.Quote) error {
			return nil
		}).AnyTimes()

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
		TryQuotesRetries:  1,
		TryQuotesTimeout:  time.Millisecond,
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      1,
				Amount:        1,
				Width:         1,
				Height:        1,
				Length:        1,
				Price:         decimal.NewFromFloat(1),
				UnitaryWeight: 1,
			},
		},
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	resultQuotes, err := svc.GetFreteRapidoQuotes(ctx, &quote)
	assert.NotNil(t, resultQuotes)
	assert.Len(t, resultQuotes, 1)
	assert.NoError(t, err)

}

func TestQuoteService_ProcessQuotes_BadInput(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	clickhouseMock.
		EXPECT().
		AddQuotes(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []quotes.Quote) error {
			return nil
		}).Times(0)

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quotes := []quotes.Quote{}

	svc.ProcessQuotes(t.Context(), quotes)
}

func TestQuoteService_ProcessQuotes_Successful(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	clickhouseMock.
		EXPECT().
		AddQuotes(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []quotes.Quote) error {
			return nil
		}).Times(1)

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
		AddQuotesRetries:  1,
		AddQuoesTimeout:   time.Millisecond,
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quotes := []quotes.Quote{
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "SUT",
			Service:   "SUT",
			Deadline:  1,
			Price:     decimal.NewFromFloat(1),
			CreatedAt: time.Now(),
		},
	}

	svc.ProcessQuotes(t.Context(), quotes)
}

func TestQuoteService_ProcessQuotes_Failure(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	clickhouseMock.
		EXPECT().
		AddQuotes(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ []quotes.Quote) error {
			return fmt.Errorf("SUT")
		}).Times(1)

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
		AddQuotesRetries:  1,
		AddQuoesTimeout:   time.Millisecond,
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quotes := []quotes.Quote{
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "SUT",
			Service:   "SUT",
			Deadline:  1,
			Price:     decimal.NewFromFloat(1),
			CreatedAt: time.Now(),
		},
	}

	svc.ProcessQuotes(t.Context(), quotes)
}

func TestQuoteService_ProcessQuotes_FailureThenWorks(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(ctrl)

	clickhouseMock := quotesmocks.NewMockClickhouseQuotesRepository(ctrl)
	gomock.InOrder(
		clickhouseMock.
			EXPECT().
			AddQuotes(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("temporary failure")),
		clickhouseMock.
			EXPECT().
			AddQuotes(gomock.Any(), gomock.Any()).
			Return(nil),
	)

	quoteCfg := &system.QuotesServiceConfig{
		RegisteredNumber:  "SUT",
		Token:             "SUT",
		PlatformCode:      "SUT",
		DispatcherZipCode: 0,
		Logger:            slog.New(slog.NewJSONHandler(io.Discard, nil)),
		AddQuotesRetries:  2,
		AddQuoesTimeout:   time.Millisecond,
	}
	svc := quotes.NewQuoteService(quoteCfg, freteRapidoMock, clickhouseMock)

	quotes := []quotes.Quote{
		{
			ID:        uuid.Must(uuid.NewV7()),
			Name:      "SUT",
			Service:   "SUT",
			Deadline:  1,
			Price:     decimal.NewFromFloat(1),
			CreatedAt: time.Now(),
		},
	}

	svc.ProcessQuotes(t.Context(), quotes)
}
