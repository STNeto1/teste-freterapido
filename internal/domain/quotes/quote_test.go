package quotes_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stretchr/testify/assert"
)

func TestQuote_ValidateCategories(t *testing.T) {
	t.Parallel()

	for categoryID := range quotes.CategoryMap {
		quote := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{
				Address: quotes.RequestQuoteRecipientAddress{
					Zipcode: "123",
				},
			},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      categoryID,
					Amount:        1,
					UnitaryWeight: 1,
					Price:         decimal.NewFromFloat(1),
					Sku:           "SUT",
					Height:        1,
					Width:         1,
					Length:        1,
				},
			},
		}

		assert.Nil(t, quote.ErrorSet())
	}

	invalidQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "123",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      9999,
				Amount:        1,
				UnitaryWeight: 1,
				Price:         decimal.NewFromFloat(1),
				Sku:           "",
				Height:        1,
				Width:         1,
				Length:        1,
			},
		},
	}
	assert.Equal(t, []string{
		"Category on volume 1 is invalid",
	}, invalidQuote.ErrorSet())
}

func TestRequestQuote_ParseRecipientZipcode(t *testing.T) {
	t.Parallel()

	validQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "01311000",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{},
		},
	}
	validParsedZipcode, err := validQuote.ParseRecipientZipcode()
	assert.Equal(t, int64(1311000), validParsedZipcode)
	assert.NoError(t, err)

	invalidQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "abc",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{},
		},
	}
	invalidParsedZipcode, err := invalidQuote.ParseRecipientZipcode()
	assert.Equal(t, int64(-1), invalidParsedZipcode)
	assert.Error(t, err)
	assert.ErrorIs(t, err, quotes.QuoteInvalidZipcodeError{Message: "Invalid recipient zipcode"})
}

func TestRequestQuote_MustParseRecipientZipcode(t *testing.T) {
	t.Parallel()

	validQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "01311000",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{},
		},
	}
	validParsedZipcode := validQuote.MustParseRecipientZipcode()
	assert.Equal(t, int64(1311000), validParsedZipcode)

	invalidQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{
				Zipcode: "abc",
			},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{},
		},
	}
	invalidParsedZipcode := invalidQuote.MustParseRecipientZipcode()
	assert.Equal(t, int64(-1), invalidParsedZipcode)
}

func TestQuote_ErrorSetFullyInvalid(t *testing.T) {
	t.Parallel()

	invalidQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{
			Address: quotes.RequestQuoteRecipientAddress{},
		},
		Volumes: []quotes.RequestQuoteVolume{
			{},
		},
	}
	validationErrors := invalidQuote.ErrorSet()
	assert.Equal(t, []string{
		"Invalid recipient zipcode",
		"Category on volume 1 is invalid",
		"Dimensions on volume 1 are invalid",
		"Price on volume 1 is invalid",
		"Weight on volume 1 is invalid",
		"Amount on volume 1 is invalid",
	}, validationErrors)
}

func TestQuote_ErrorSetFullyWithoutErrors(t *testing.T) {
	t.Parallel()

	invalidQuote := quotes.RequestQuote{
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
	validationErrors := invalidQuote.ErrorSet()
	assert.Nil(t, validationErrors)
}
