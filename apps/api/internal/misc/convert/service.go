package convert

import (
	"errors"
	"fmt"
	"math"
)

// unit represents a single unit with its conversion factor relative to the base unit.
// For temperature, the factor field is unused; conversions are handled separately.
type unit struct {
	category string
	// toBase converts a value in this unit to the base unit.
	toBase func(v float64) float64
	// fromBase converts a value from the base unit to this unit.
	fromBase func(v float64) float64
	// factor is toBase(1) – the number of base units in one of this unit.
	// Used only for building the human-readable formula for linear units.
	factor float64
}

var units = map[string]unit{
	// Length (base: meters)
	"mm":    {category: "length", factor: 0.001, toBase: mul(0.001), fromBase: mul(1 / 0.001)},
	"cm":    {category: "length", factor: 0.01, toBase: mul(0.01), fromBase: mul(1 / 0.01)},
	"m":     {category: "length", factor: 1, toBase: mul(1), fromBase: mul(1)},
	"km":    {category: "length", factor: 1000, toBase: mul(1000), fromBase: mul(1 / 1000.0)},
	"in":    {category: "length", factor: 0.0254, toBase: mul(0.0254), fromBase: mul(1 / 0.0254)},
	"ft":    {category: "length", factor: 0.3048, toBase: mul(0.3048), fromBase: mul(1 / 0.3048)},
	"yd":    {category: "length", factor: 0.9144, toBase: mul(0.9144), fromBase: mul(1 / 0.9144)},
	"miles": {category: "length", factor: 1609.344, toBase: mul(1609.344), fromBase: mul(1 / 1609.344)},
	"nmi":   {category: "length", factor: 1852, toBase: mul(1852), fromBase: mul(1 / 1852.0)},

	// Weight (base: grams)
	"mg":    {category: "weight", factor: 0.001, toBase: mul(0.001), fromBase: mul(1 / 0.001)},
	"g":     {category: "weight", factor: 1, toBase: mul(1), fromBase: mul(1)},
	"kg":    {category: "weight", factor: 1000, toBase: mul(1000), fromBase: mul(1 / 1000.0)},
	"t":     {category: "weight", factor: 1e6, toBase: mul(1e6), fromBase: mul(1 / 1e6)},
	"oz":    {category: "weight", factor: 28.3495, toBase: mul(28.3495), fromBase: mul(1 / 28.3495)},
	"lb":    {category: "weight", factor: 453.592, toBase: mul(453.592), fromBase: mul(1 / 453.592)},
	"stone": {category: "weight", factor: 6350.29, toBase: mul(6350.29), fromBase: mul(1 / 6350.29)},

	// Volume (base: milliliters)
	"ml":    {category: "volume", factor: 1, toBase: mul(1), fromBase: mul(1)},
	"l":     {category: "volume", factor: 1000, toBase: mul(1000), fromBase: mul(1 / 1000.0)},
	"tsp":   {category: "volume", factor: 4.92892, toBase: mul(4.92892), fromBase: mul(1 / 4.92892)},
	"tbsp":  {category: "volume", factor: 14.7868, toBase: mul(14.7868), fromBase: mul(1 / 14.7868)},
	"fl_oz": {category: "volume", factor: 29.5735, toBase: mul(29.5735), fromBase: mul(1 / 29.5735)},
	"cup":   {category: "volume", factor: 236.588, toBase: mul(236.588), fromBase: mul(1 / 236.588)},
	"pt":    {category: "volume", factor: 473.176, toBase: mul(473.176), fromBase: mul(1 / 473.176)},
	"qt":    {category: "volume", factor: 946.353, toBase: mul(946.353), fromBase: mul(1 / 946.353)},
	"gal":   {category: "volume", factor: 3785.41, toBase: mul(3785.41), fromBase: mul(1 / 3785.41)},

	// Temperature (base: celsius; special handling)
	"c": {category: "temperature"},
	"f": {category: "temperature"},
	"k": {category: "temperature"},

	// Area (base: square meters)
	"mm2":  {category: "area", factor: 1e-6, toBase: mul(1e-6), fromBase: mul(1 / 1e-6)},
	"cm2":  {category: "area", factor: 1e-4, toBase: mul(1e-4), fromBase: mul(1 / 1e-4)},
	"m2":   {category: "area", factor: 1, toBase: mul(1), fromBase: mul(1)},
	"km2":  {category: "area", factor: 1e6, toBase: mul(1e6), fromBase: mul(1 / 1e6)},
	"in2":  {category: "area", factor: 6.4516e-4, toBase: mul(6.4516e-4), fromBase: mul(1 / 6.4516e-4)},
	"ft2":  {category: "area", factor: 0.092903, toBase: mul(0.092903), fromBase: mul(1 / 0.092903)},
	"yd2":  {category: "area", factor: 0.836127, toBase: mul(0.836127), fromBase: mul(1 / 0.836127)},
	"acre": {category: "area", factor: 4046.86, toBase: mul(4046.86), fromBase: mul(1 / 4046.86)},
	"ha":   {category: "area", factor: 10000, toBase: mul(10000), fromBase: mul(1 / 10000.0)},

	// Speed (base: km/h)
	"m_s":   {category: "speed", factor: 3.6, toBase: mul(3.6), fromBase: mul(1 / 3.6)},
	"km_h":  {category: "speed", factor: 1, toBase: mul(1), fromBase: mul(1)},
	"mph":   {category: "speed", factor: 1.60934, toBase: mul(1.60934), fromBase: mul(1 / 1.60934)},
	"knots": {category: "speed", factor: 1.852, toBase: mul(1.852), fromBase: mul(1 / 1.852)},
}

func mul(f float64) func(float64) float64 {
	return func(v float64) float64 { return v * f }
}

// ErrUnknownUnit is returned when a unit is not recognised.
var ErrUnknownUnit = errors.New("unknown unit")

// ErrIncompatibleUnits is returned when from and to belong to different categories.
var ErrIncompatibleUnits = errors.New("incompatible units: cannot convert between different measurement types")

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Convert(from, to string, value float64) (Result, error) {
	fromUnit, ok := units[from]
	if !ok {
		return Result{}, fmt.Errorf("%w: %q", ErrUnknownUnit, from)
	}

	toUnit, ok := units[to]
	if !ok {
		return Result{}, fmt.Errorf("%w: %q", ErrUnknownUnit, to)
	}

	if fromUnit.category != toUnit.category {
		return Result{}, ErrIncompatibleUnits
	}

	var result float64
	var formula string

	if fromUnit.category == "temperature" {
		result, formula = convertTemperature(from, to, value)
	} else {
		factor := fromUnit.factor / toUnit.factor
		result = roundTo(value*factor, 6)
		formula = fmt.Sprintf("%s × %s", from, formatFactor(factor))
	}

	return Result{
		From:    from,
		To:      to,
		Input:   value,
		Result:  result,
		Formula: formula,
	}, nil
}

func convertTemperature(from, to string, value float64) (float64, string) {
	if from == to {
		return roundTo(value, 6), fmt.Sprintf("%s (no conversion needed)", from)
	}

	var celsius float64

	switch from {
	case "c":
		celsius = value
	case "f":
		celsius = (value - 32) * 5 / 9
	case "k":
		celsius = value - 273.15
	}

	var result float64
	var formula string

	switch to {
	case "c":
		result = celsius
	case "f":
		result = celsius*9/5 + 32
	case "k":
		result = celsius + 273.15
	}

	result = roundTo(result, 6)

	switch {
	case from == "c" && to == "f":
		formula = "°C × 9/5 + 32"
	case from == "f" && to == "c":
		formula = "(°F − 32) × 5/9"
	case from == "c" && to == "k":
		formula = "°C + 273.15"
	case from == "k" && to == "c":
		formula = "K − 273.15"
	case from == "f" && to == "k":
		formula = "(°F − 32) × 5/9 + 273.15"
	case from == "k" && to == "f":
		formula = "(K − 273.15) × 9/5 + 32"
	}

	return result, formula
}

func roundTo(v float64, decimals int) float64 {
	p := math.Pow(10, float64(decimals))
	return math.Round(v*p) / p
}

func formatFactor(f float64) string {
	// Round to 6 decimal places first to avoid floating-point artefacts.
	r := roundTo(f, 6)

	if r == math.Trunc(r) {
		return fmt.Sprintf("%.0f", r)
	}

	s := fmt.Sprintf("%.6f", r)

	// Trim trailing zeros but keep at least one decimal place.
	i := len(s) - 1
	for i > 0 && s[i] == '0' {
		i--
	}

	if s[i] == '.' {
		i++
	}

	return s[:i+1]
}
