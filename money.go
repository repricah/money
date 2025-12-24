package money

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// all money operations use int to avoid floating point errors
const MONEY_INT_MULTIPLIER = 100

type DollarValue struct {
	Cents int `json:"cents"`
}

func NewMoneyValue(amount int) DollarValue {
	return DollarValue{
		Cents: amount,
	}
}

func NewZero() DollarValue {
	return DollarValue{
		Cents: 0,
	}
}

// NewOneCent returns a DollarValue representing one cent ($0.01).
func NewOneCent() DollarValue {
	return DollarValue{
		Cents: 1,
	}
}

// NewMoneyValueFromString creates a new DollarValue from a string representing dollars.
func NewMoneyValueFromString(amount string) DollarValue {
	// Remove any whitespace and split the string into whole and fractional parts
	parts := strings.Split(strings.TrimSpace(amount), ".")
	if len(parts) == 1 {
		// No decimal point, so treat the entire string as the whole part
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("Error parsing money value: %s, error: %s", amount, err))
		}
		return DollarValue{Cents: a * MONEY_INT_MULTIPLIER}
	} else if len(parts) == 2 {
		// Decimal point present, handle whole and fractional parts
		a := parts[0]
		b := parts[1]
		if len(b) == 1 {
			b += "0"
		} else if len(b) > 2 {
			b = b[:2]
		}
		ai, err := strconv.Atoi(a)
		if err != nil {
			panic(fmt.Sprintf("Error parsing whole part of money value: %s, error: %s", amount, err))
		}
		bi, err := strconv.Atoi(b)
		if err != nil {
			panic(fmt.Sprintf("Error parsing fractional part of money value: %s, error: %s", amount, err))
		}
		return NewMoneyValue(ai*MONEY_INT_MULTIPLIER + bi)
	}
	panic(fmt.Sprintf("Error parsing money value: %s", amount))
}

func NewMoneyValueFromDollars(amount float64) DollarValue {
	return DollarValue{
		Cents: int(math.Round(amount * MONEY_INT_MULTIPLIER)),
	}
}

func (m DollarValue) Multiply(multiplier int) DollarValue {
	return DollarValue{
		Cents: m.Cents * multiplier,
	}
}

func (m DollarValue) MultiplyFloat(multiplier float64) DollarValue {
	return DollarValue{
		Cents: int(math.Round(float64(m.Cents) * multiplier)),
	}
}

func (m DollarValue) MultiplyRatio(numerator, denominator int) DollarValue {
	if denominator == 0 {
		panic("Divide by zero")
	}

	// Use int64 for intermediate calculation to avoid overflow on 32-bit systems
	product := int64(m.Cents) * int64(numerator)
	denom := int64(denominator)

	var rounding int64
	if product < 0 {
		rounding = -denom / 2
	} else {
		rounding = denom / 2
	}

	// Standard rounding (round half up): (m.Cents * numerator + denominator/2) / denominator.
	// Note: Previous code used Banker's rounding (math.RoundToEven), but standard rounding is common for fees and simpler to implement.
	return DollarValue{
		Cents: int((product + rounding) / denom),
	}
}

func (m DollarValue) Divide(divisor int) DollarValue {
	if divisor == 0 {
		panic("Divide by zero")
	}
	return DollarValue{
		Cents: m.Cents / divisor,
	}
}

func (m DollarValue) Add(other DollarValue) DollarValue {
	return DollarValue{
		Cents: m.Cents + other.Cents,
	}
}

func (m DollarValue) Subtract(other DollarValue) DollarValue {
	return DollarValue{
		Cents: m.Cents - other.Cents,
	}
}

func (m DollarValue) IsZero() bool {
	return m.Cents == 0
}

// Compare returns -1 if m < other, 0 if m == other, and 1 if m > other.
func (m DollarValue) Compare(other DollarValue) int {
	switch {
	case m.Cents < other.Cents:
		return -1
	case m.Cents > other.Cents:
		return 1
	default:
		return 0
	}
}

func (m DollarValue) GreaterThan(other DollarValue) bool {
	return m.Compare(other) > 0
}

func (m DollarValue) LessThan(other DollarValue) bool {
	return m.Compare(other) < 0
}

func (m DollarValue) Equal(other DollarValue) bool {
	return m.Compare(other) == 0
}

func (m DollarValue) To2DecimalString() string {
	return fmt.Sprintf("%.2f", float64(m.Cents)/MONEY_INT_MULTIPLIER)
}

func (m DollarValue) PercentageDifference(other DollarValue) float64 {
	if other.Cents == 0 {
		return 1
	}
	difference := m.Cents - other.Cents
	return float64(difference) / float64(other.Cents)
}

func (m DollarValue) AbsDifference(other DollarValue) DollarValue {
	amount := m.Cents - other.Cents
	if amount < 0 {
		amount = -amount
	}
	return DollarValue{
		Cents: amount,
	}
}

func (r *DollarValue) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		*r = NewMoneyValue(0)
		return
	}
	*r = NewMoneyValueFromString(csv)
	return
}
func (m DollarValue) MarshalCSV() (string, error) {
	return m.To2DecimalString(), nil
}

func (m DollarValue) ToCents() int {
	return m.Cents
}

func (m DollarValue) GreaterThanOrEqualTo(value DollarValue) bool {
	return m.Compare(value) >= 0
}

func (m DollarValue) LessThanOrEqualTo(value DollarValue) bool {
	return m.Compare(value) <= 0
}

func (m DollarValue) BetweenInclusive(min DollarValue, max DollarValue) bool {
	return m.Compare(min) >= 0 && m.Compare(max) <= 0
}

func MoneyValueHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {

		// Ensure the source field is numeric (float64 or int)
		if f.Kind() != reflect.Float64 && f.Kind() != reflect.Int {
			return data, nil
		}

		// Ensure the target type is DollarValue
		if t != reflect.TypeOf(DollarValue{}) {
			return data, nil
		}

		// Safely convert the data to float64
		switch v := data.(type) {
		case float64:
			return NewMoneyValueFromDollars(v), nil
		case int:
			return NewMoneyValueFromDollars(float64(v)), nil
		default:
			return nil, fmt.Errorf("unsupported type: %T", data)
		}
	}
}

func MaxPrice(numbers ...DollarValue) DollarValue {
	if len(numbers) == 0 {
		panic("No values found")
	}

	maxValue := numbers[0]
	for _, num := range numbers {
		if num.GreaterThan(maxValue) {
			maxValue = num
		}
	}
	return maxValue
}
