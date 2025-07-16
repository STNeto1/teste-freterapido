package quotes_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/mocks/quotesmocks"
	"github.com/stretchr/testify/assert"
)

func TestQuoteService_CreateRequestPayload(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	result := svc.CreateRequestPayload(
		&quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
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
			Type:             0,
			RegisteredNumber: "SUT",
			StateInscription: "SUT",
			Country:          "BRA",
			Zipcode:          0,
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

func TestQuoteService_GetFreteRapidoQuotes_InvalidZipcode(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "abc",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidZipcodeError{
		Message: "Invalid recipient zipcode",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidCategory(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category: 0,
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidCategoryError{
		Message: "Category on volume 1 is invalid",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidAmount(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category: 1,
				Amount:   0,
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidAmountError{
		Message: "Amount on volume 1 is invalid",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidDimensions(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category: 1,
				Amount:   1,
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidDimensionError{
		Message: "Dimensions on volume 1 are invalid",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidPrices(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category: 1,
				Amount:   1,
				Width:    1,
				Height:   1,
				Length:   1,
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidPriceError{
		Message: "Price on volume 1 is invalid",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_InvalidWeights(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

	quote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category: 1,
				Amount:   1,
				Width:    1,
				Height:   1,
				Length:   1,
				Price:    decimal.NewFromFloat(1),
			},
		},
	}

	result, err := svc.GetFreteRapidoQuotes(t.Context(), &quote)
	assert.Nil(t, result)
	assert.Equal(t, err, quotes.QuoteInvalidWeightError{
		Message: "Weight on volume 1 is invalid",
	})
}

func TestQuoteService_GetFreteRapidoQuotes_NoValidationError(t *testing.T) {
	t.Parallel()

	freteRapidoMock := quotesmocks.NewMockFreteRapidoQuotesRepository(t)

	svc := quotes.NewQuoteService(freteRapidoMock)

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
