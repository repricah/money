# Money

A Go library for safe, integer-based monetary calculations. Avoids floating-point errors by storing all values in cents.

[![Go Reference](https://pkg.go.dev/badge/github.com/repricah/money.svg)](https://pkg.go.dev/github.com/repricah/money)
[![Go Report Card](https://goreportcard.com/badge/github.com/repricah/money)](https://goreportcard.com/report/github.com/repricah/money)
[![CI](https://github.com/repricah/money/actions/workflows/ci.yml/badge.svg)](https://github.com/repricah/money/actions/workflows/ci.yml)

## Installation

```bash
go get github.com/repricah/money
```

## Features

- **Integer-based arithmetic** - All operations use cents to avoid floating-point errors
- **Comprehensive operations** - Add, subtract, multiply, divide, compare
- **String/float parsing** - Create values from strings like "19.99" or floats
- **CSV marshaling** - Built-in CSV marshal/unmarshal support
- **Mapstructure hook** - Decode hook for configuration parsing

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/repricah/money"
)

func main() {
    // Create from various sources
    price := money.NewMoneyValueFromString("19.99")
    tax := money.NewMoneyValueFromDollars(1.50)
    discount := money.NewMoneyValue(200) // 200 cents = $2.00

    // Arithmetic
    subtotal := price.Add(tax)
    total := subtotal.Subtract(discount)

    fmt.Printf("Total: $%s\n", total.To2DecimalString()) // Total: $19.49

    // Comparisons
    if total.GreaterThan(money.NewZero()) {
        fmt.Println("You owe money!")
    }

    // Percentage calculations (e.g., 8.95% fee)
    fee := total.MultiplyRatio(895, 10000)
    fmt.Printf("Fee: $%s\n", fee.To2DecimalString())
}
```

## API Reference

### Creating Values

```go
// From cents (integer)
v := money.NewMoneyValue(1999)           // $19.99

// From dollars (float)
v := money.NewMoneyValueFromDollars(19.99) // $19.99

// From string
v := money.NewMoneyValueFromString("19.99") // $19.99

// Zero value
v := money.NewZero()                      // $0.00

// One cent
v := money.NewOneCent()                   // $0.01
```

### Arithmetic Operations

```go
a := money.NewMoneyValue(1000) // $10.00
b := money.NewMoneyValue(300)  // $3.00

a.Add(b)              // $13.00
a.Subtract(b)         // $7.00
a.Multiply(2)         // $20.00
a.MultiplyFloat(1.5)  // $15.00
a.Divide(2)           // $5.00
a.MultiplyRatio(1, 3) // $3.33 (rounded)
```

### Comparison Operations

```go
a.GreaterThan(b)         // true
a.LessThan(b)            // false
a.Equal(b)               // false
a.GreaterThanOrEqualTo(b) // true
a.LessThanOrEqualTo(b)   // false
a.Compare(b)             // 1 (a > b), 0 (equal), -1 (a < b)
a.BetweenInclusive(min, max) // true if min <= a <= max
```

### Utility Methods

```go
v.IsZero()              // true if value is $0.00
v.ToCents()             // int value in cents
v.To2DecimalString()    // "19.99"
v.AbsDifference(other)  // absolute difference
v.PercentageDifference(other) // percentage difference as float
```

## Mapstructure Integration

Use the decode hook for parsing configuration files:

```go
import "github.com/mitchellh/mapstructure"

type Config struct {
    MinPrice money.DollarValue
    MaxPrice money.DollarValue
}

decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
    DecodeHook: money.MoneyValueHookFunc(),
    Result:     &config,
})
decoder.Decode(rawConfig)
```

## License

MIT License - see LICENSE file for details.
