package quotes_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stneto1/teste-freterapido/internal/domain/quotes"
	"github.com/stneto1/teste-freterapido/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestQuote_ValidateCategories(t *testing.T) {
	t.Parallel()

	for categoryID := range quotes.CategoryMap {
		quote := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      categoryID,
					Amount:        0,
					UnitaryWeight: 0,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        0,
					Width:         0,
					Length:        0,
				},
			},
		}

		assert.Equal(t, -1, quote.ValidateCategories())
	}

	invalidQuote := quotes.RequestQuote{
		Recipient: quotes.RequestQuoteRecipient{},
		Volumes: []quotes.RequestQuoteVolume{
			{
				Category:      9999,
				Amount:        0,
				UnitaryWeight: 0,
				Price:         decimal.NewFromFloat(0),
				Sku:           "",
				Height:        0,
				Width:         0,
				Length:        0,
			},
		},
	}
	assert.Equal(t, 0, invalidQuote.ValidateCategories())
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

func TestRequestQuote_ValidateDimensions(t *testing.T) {
	t.Parallel()

	dimensions := utils.Map(utils.RangeWithStep(-1, 1, 1), func(val int) float64 {
		return float64(val)
	})

	for _, heightDimension := range dimensions {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        0,
					UnitaryWeight: 0,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        heightDimension,
					Width:         10,
					Length:        10,
				},
			},
		}

		if heightDimension > 0 {
			assert.Equal(t, -1, sample.ValidateDimensions())
		} else {
			assert.Equal(t, 0, sample.ValidateDimensions())
		}
	}

	for _, widthDimension := range dimensions {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        0,
					UnitaryWeight: 0,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        10,
					Width:         widthDimension,
					Length:        10,
				},
			},
		}

		if widthDimension > 0 {
			assert.Equal(t, -1, sample.ValidateDimensions())
		} else {
			assert.Equal(t, 0, sample.ValidateDimensions())
		}
	}

	for _, lengthDimension := range dimensions {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        0,
					UnitaryWeight: 0,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        10,
					Width:         10,
					Length:        lengthDimension,
				},
			},
		}

		if lengthDimension > 0 {
			assert.Equal(t, -1, sample.ValidateDimensions())
		} else {
			assert.Equal(t, 0, sample.ValidateDimensions())
		}
	}
}

func TestRequestQuote_ValidateAmount(t *testing.T) {
	t.Parallel()

	amounts := utils.RangeWithStep(-1, 1, 1)

	for _, amount := range amounts {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        amount,
					UnitaryWeight: 0,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        0,
					Width:         0,
					Length:        0,
				},
			},
		}

		if amount > 0 {
			assert.Equal(t, -1, sample.ValidateAmount())
		} else {
			assert.Equal(t, 0, sample.ValidateAmount())
		}
	}
}

func TestRequestQuote_ValidatePrice(t *testing.T) {
	t.Parallel()

	prices := utils.Map(utils.RangeWithStep(-1, 1, 1), func(val int) decimal.Decimal {
		return decimal.NewFromInt(int64(val))
	})

	for _, price := range prices {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        0,
					UnitaryWeight: 0,
					Price:         price,
					Sku:           "",
					Height:        0,
					Width:         0,
					Length:        0,
				},
			},
		}

		if price.IsPositive() {
			assert.Equal(t, -1, sample.ValidatePrice())
		} else {
			assert.Equal(t, 0, sample.ValidatePrice())
		}
	}
}

func TestRequestQuote_ValidateWeight(t *testing.T) {
	t.Parallel()

	weights := utils.Map(utils.RangeWithStep(-1, 1, 1), func(val int) float64 {
		return float64(val)
	})

	for _, weight := range weights {
		sample := quotes.RequestQuote{
			Recipient: quotes.RequestQuoteRecipient{},
			Volumes: []quotes.RequestQuoteVolume{
				{
					Category:      0,
					Amount:        0,
					UnitaryWeight: weight,
					Price:         decimal.NewFromFloat(0),
					Sku:           "",
					Height:        0,
					Width:         0,
					Length:        0,
				},
			},
		}

		if weight > 0 {
			assert.Equal(t, -1, sample.ValidateWeight())
		} else {
			assert.Equal(t, 0, sample.ValidateWeight())
		}
	}
}
