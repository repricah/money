package money

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMoneyValue(t *testing.T) {
	mv := NewMoneyValue(500)
	assert.Equal(t, 500, mv.Cents)
}

func TestNewMoneyValueFromDollars(t *testing.T) {
	mv := NewMoneyValueFromDollars(5.00)
	assert.Equal(t, 500, mv.Cents)
}

func TestMultiply(t *testing.T) {
	mv := NewMoneyValue(100)
	result := mv.Multiply(2)
	assert.Equal(t, 200, result.Cents)
}

func TestMultiplyFloat(t *testing.T) {
	mv := NewMoneyValue(100)
	result := mv.MultiplyFloat(2.5)
	assert.Equal(t, 250, result.Cents)
}

func TestMultiplyRatio(t *testing.T) {
	// 100 cents * 895 / 10000 = 8.95 -> 9 cents (rounded)
	mv := NewMoneyValue(100)
	result := mv.MultiplyRatio(895, 10000)
	assert.Equal(t, 9, result.Cents)

	// 50 cents * 5000 / 10000 = 25 cents
	mv = NewMoneyValue(50)
	result = mv.MultiplyRatio(5000, 10000)
	assert.Equal(t, 25, result.Cents)

	// Rounding check:
	// 10 * 1 / 3 = 3.33 -> 3
	mv = NewMoneyValue(10)
	result = mv.MultiplyRatio(1, 3)
	assert.Equal(t, 3, result.Cents)

	// 10 * 2 / 3 = 6.66 -> 7
	mv = NewMoneyValue(10)
	result = mv.MultiplyRatio(2, 3)
	assert.Equal(t, 7, result.Cents)

	// 10 * 1 / 2 = 5
	mv = NewMoneyValue(10)
	result = mv.MultiplyRatio(1, 2)
	assert.Equal(t, 5, result.Cents)
}

func TestMultiplyRatio_PanicOnZeroDenominator(t *testing.T) {
	assert.Panics(t, func() {
		mv := NewMoneyValue(100)
		mv.MultiplyRatio(1, 0)
	})
}

func TestMultiplyRatio_OverflowCheck(t *testing.T) {
	// Test potential overflow with large numbers
	// Max int32 is approx 2 billion
	// 10,000,000 cents ($100,000) * 895 = 8,950,000,000 (overflows int32)

	largeVal := NewMoneyValue(10_000_000) // $100,000.00
	numerator := 895
	denominator := 10000

	// Expected: 10,000,000 * 895 / 10000 = 895,000
	expected := 895_000

	result := largeVal.MultiplyRatio(numerator, denominator)
	assert.Equal(t, expected, result.Cents, "Should handle large numbers without overflow")
}

func TestMultiplyRatio_NegativeNumbers(t *testing.T) {
	// Test with negative value
	negVal := NewMoneyValue(-100) // -$1.00
	numerator := 5000             // 50%
	denominator := 10000

	// Expected: -100 * 5000 / 10000 = -50
	expected := -50

	result := negVal.MultiplyRatio(numerator, denominator)
	assert.Equal(t, expected, result.Cents, "Should handle negative numbers correctly")

	// Test rounding for negative numbers
	// -10 * 1 / 3 = -3.33 -> -3
	negVal = NewMoneyValue(-10)
	result = negVal.MultiplyRatio(1, 3)
	assert.Equal(t, -3, result.Cents)

	// -10 * 2 / 3 = -6.66 -> -7
	negVal = NewMoneyValue(-10)
	result = negVal.MultiplyRatio(2, 3)
	assert.Equal(t, -7, result.Cents)
}

func TestAdd(t *testing.T) {
	mv1 := NewMoneyValue(100)
	mv2 := NewMoneyValue(200)
	result := mv1.Add(mv2)
	assert.Equal(t, 300, result.Cents)
}

func TestSubtract(t *testing.T) {
	mv1 := NewMoneyValue(200)
	mv2 := NewMoneyValue(100)
	result := mv1.Subtract(mv2)
	assert.Equal(t, 100, result.Cents)
}

func TestIsZero(t *testing.T) {
	mv := NewMoneyValue(0)
	assert.True(t, mv.IsZero())

	mv = NewMoneyValue(100)
	assert.False(t, mv.IsZero())
}

func TestGreaterThan(t *testing.T) {
	mv1 := NewMoneyValue(200)
	mv2 := NewMoneyValue(100)
	assert.True(t, mv1.GreaterThan(mv2))
	assert.False(t, mv2.GreaterThan(mv1))
}

func TestLessThan(t *testing.T) {
	mv1 := NewMoneyValue(100)
	mv2 := NewMoneyValue(200)
	assert.True(t, mv1.LessThan(mv2))
	assert.False(t, mv2.LessThan(mv1))
}

func TestEqual(t *testing.T) {
	mv1 := NewMoneyValue(100)
	mv2 := NewMoneyValue(100)
	mv3 := NewMoneyValue(200)
	assert.True(t, mv1.Equal(mv2))
	assert.False(t, mv1.Equal(mv3))
}

func TestCompare(t *testing.T) {
	t.Parallel()

	mv1 := NewMoneyValue(100)
	mv2 := NewMoneyValue(200)

	assert.Equal(t, -1, mv1.Compare(mv2))
	assert.Equal(t, 1, mv2.Compare(mv1))
	assert.Equal(t, 0, mv1.Compare(NewMoneyValue(100)))
}

func TestTo2DecimalString(t *testing.T) {
	mv := NewMoneyValue(125)
	assert.Equal(t, "1.25", mv.To2DecimalString())
}

func TestPercentageDifference(t *testing.T) {
	mv1 := NewMoneyValue(200)
	mv2 := NewMoneyValue(100)
	assert.Equal(t, 1.0, mv1.PercentageDifference(mv2))

	mv3 := NewMoneyValue(50)
	assert.Equal(t, -0.5, mv3.PercentageDifference(mv2))

	// test the return 1 case
	mv4 := NewMoneyValue(0)
	assert.Equal(t, 1.0, mv4.PercentageDifference(NewMoneyValue(0)))

}

func TestAbsDifference(t *testing.T) {
	mv1 := NewMoneyValue(200)
	mv2 := NewMoneyValue(100)
	result := mv1.AbsDifference(mv2)
	assert.Equal(t, 100, result.Cents)

	mv3 := NewMoneyValue(50)
	result = mv3.AbsDifference(mv2)
	assert.Equal(t, 50, result.Cents)
}

func TestUnmarshalCSV(t *testing.T) {
	var mv DollarValue
	err := mv.UnmarshalCSV("1.25")
	assert.NoError(t, err)
	assert.Equal(t, 125, mv.Cents)

	err = mv.UnmarshalCSV("")
	assert.NoError(t, err)
	assert.Equal(t, 0, mv.Cents)
}

func TestMarshalCSV(t *testing.T) {
	mv := NewMoneyValue(125)
	csv, err := mv.MarshalCSV()
	assert.NoError(t, err)
	assert.Equal(t, "1.25", csv)
}

// test string conversion
func TestNewMoneyValueFromString(t *testing.T) {
	// multiple cases
	mv := NewMoneyValueFromString("1.25")
	assert.Equal(t, 125, mv.Cents)

	mv = NewMoneyValueFromString("1.00")
	assert.Equal(t, 100, mv.Cents)

	mv = NewMoneyValueFromString("0.00")
	assert.Equal(t, 0, mv.Cents)

	mv = NewMoneyValueFromString("0.29")
	assert.Equal(t, 29, mv.Cents)

	mv = NewMoneyValueFromString("6.4400")
	assert.Equal(t, 644, mv.Cents)

	mv = NewMoneyValueFromString("17")
	assert.Equal(t, 1700, mv.Cents)

	mv = NewMoneyValueFromString("3.4")
	assert.Equal(t, 340, mv.Cents)

	mv = NewMoneyValueFromString("2.345")
	assert.Equal(t, 234, mv.Cents)

	assert.Panics(t, func() {
		NewMoneyValueFromString("abc")
	})

	assert.Panics(t, func() {
		NewMoneyValueFromString("1.2.3")
	})
	assert.Panics(t, func() {
		NewMoneyValueFromString("abc.12")
	})

	assert.Panics(t, func() {
		NewMoneyValueFromString("12.ab")
	})
}

func TestMaxPrice_PanicWhenNoPricesPassed(t *testing.T) {
	assert.Panics(t, func() {
		MaxPrice()
	})
}

func TestMaxPriceGetsMaxPrice(t *testing.T) {
	t.Parallel()
	cases := []struct {
		prices []DollarValue
		max    DollarValue
	}{
		{
			prices: []DollarValue{
				NewMoneyValue(100),
				NewMoneyValueFromDollars(1000),
			},
			max: NewMoneyValueFromDollars(1000),
		},
	}
	for _, c := range cases {
		t.Run("MaxPrice", func(t *testing.T) {
			maxPrice := MaxPrice(c.prices...)
			assert.Equal(t, c.max, maxPrice)
		})
	}
}

func TestNewZero(t *testing.T) {
	mv := NewZero()
	assert.Equal(t, 0, mv.Cents)
	assert.True(t, mv.IsZero())
}

func TestNewOneCent(t *testing.T) {
	mv := NewOneCent()
	assert.Equal(t, 1, mv.Cents)
	assert.False(t, mv.IsZero())
}

func TestDivide(t *testing.T) {
	mv := NewMoneyValue(100)
	result := mv.Divide(2)
	assert.Equal(t, 50, result.Cents)

	mv = NewMoneyValue(300)
	result = mv.Divide(3)
	assert.Equal(t, 100, result.Cents)
}

func TestDivide_PanicOnZero(t *testing.T) {
	assert.Panics(t, func() {
		mv := NewMoneyValue(100)
		mv.Divide(0)
	})
}

func TestToCents(t *testing.T) {
	mv := NewMoneyValue(125)
	assert.Equal(t, 125, mv.ToCents())

	mv = NewMoneyValueFromDollars(5.50)
	assert.Equal(t, 550, mv.ToCents())
}

func TestGreaterThanOrEqualTo(t *testing.T) {
	mv1 := NewMoneyValue(200)
	mv2 := NewMoneyValue(100)
	mv3 := NewMoneyValue(200)

	assert.True(t, mv1.GreaterThanOrEqualTo(mv2))
	assert.True(t, mv1.GreaterThanOrEqualTo(mv3))
	assert.False(t, mv2.GreaterThanOrEqualTo(mv1))
}

func TestLessThanOrEqualTo(t *testing.T) {
	mv1 := NewMoneyValue(100)
	mv2 := NewMoneyValue(200)
	mv3 := NewMoneyValue(100)

	assert.True(t, mv1.LessThanOrEqualTo(mv2))
	assert.True(t, mv1.LessThanOrEqualTo(mv3))
	assert.False(t, mv2.LessThanOrEqualTo(mv1))
}

func TestBetweenInclusive(t *testing.T) {
	min := NewMoneyValue(100)
	max := NewMoneyValue(300)

	// Within range
	mv := NewMoneyValue(200)
	assert.True(t, mv.BetweenInclusive(min, max))

	// At min boundary
	mv = NewMoneyValue(100)
	assert.True(t, mv.BetweenInclusive(min, max))

	// At max boundary
	mv = NewMoneyValue(300)
	assert.True(t, mv.BetweenInclusive(min, max))

	// Below range
	mv = NewMoneyValue(50)
	assert.False(t, mv.BetweenInclusive(min, max))

	// Above range
	mv = NewMoneyValue(400)
	assert.False(t, mv.BetweenInclusive(min, max))
}
func TestMoneyValueHookFunc(t *testing.T) {
	hook := MoneyValueHookFunc()
	type customInt int

	// Test float64 conversion
	result, err := hook(
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(DollarValue{}),
		5.50,
	)
	assert.NoError(t, err)
	mv := result.(DollarValue)
	assert.Equal(t, 550, mv.Cents)

	// Test int conversion
	result, err = hook(
		reflect.TypeOf(int(0)),
		reflect.TypeOf(DollarValue{}),
		10,
	)
	assert.NoError(t, err)
	mv = result.(DollarValue)
	assert.Equal(t, 1000, mv.Cents)

	// Test non-numeric type (should pass through)
	result, err = hook(
		reflect.TypeOf(""),
		reflect.TypeOf(DollarValue{}),
		"test",
	)
	assert.NoError(t, err)
	assert.Equal(t, "test", result)

	// Test wrong target type (should pass through)
	result, err = hook(
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(""),
		5.50,
	)
	assert.NoError(t, err)
	assert.Equal(t, 5.50, result)

	// Test unsupported numeric type
	_, err = hook(
		reflect.TypeOf(customInt(0)),
		reflect.TypeOf(DollarValue{}),
		customInt(10),
	)
	assert.EqualError(t, err, "unsupported type: money.customInt")
}
