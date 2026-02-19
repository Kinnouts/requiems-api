package convert

import (
	"fmt"
	"math"
	"strings"
)

// category groups units that can be converted between each other.
type category int

const (
	categoryLength category = iota
	categoryWeight
	categoryVolume
	categoryTemperature
	categoryArea
	categorySpeed
	categoryData
	categoryTime
)

// unitInfo holds the SI base-unit multiplier and the category for a unit name.
// Temperature units are handled separately via explicit conversion functions.
type unitInfo struct {
	factor   float64
	category category
}

// units maps lowercase unit names/aliases to their info.
var units = map[string]unitInfo{
	// Length — base: metre
	"mm":             {0.001, categoryLength},
	"millimeter":     {0.001, categoryLength},
	"millimeters":    {0.001, categoryLength},
	"millimetre":     {0.001, categoryLength},
	"millimetres":    {0.001, categoryLength},
	"cm":             {0.01, categoryLength},
	"centimeter":     {0.01, categoryLength},
	"centimeters":    {0.01, categoryLength},
	"centimetre":     {0.01, categoryLength},
	"centimetres":    {0.01, categoryLength},
	"m":              {1, categoryLength},
	"meter":          {1, categoryLength},
	"meters":         {1, categoryLength},
	"metre":          {1, categoryLength},
	"metres":         {1, categoryLength},
	"km":             {1000, categoryLength},
	"kilometer":      {1000, categoryLength},
	"kilometers":     {1000, categoryLength},
	"kilometre":      {1000, categoryLength},
	"kilometres":     {1000, categoryLength},
	"in":             {0.0254, categoryLength},
	"inch":           {0.0254, categoryLength},
	"inches":         {0.0254, categoryLength},
	"ft":             {0.3048, categoryLength},
	"foot":           {0.3048, categoryLength},
	"feet":           {0.3048, categoryLength},
	"yd":             {0.9144, categoryLength},
	"yard":           {0.9144, categoryLength},
	"yards":          {0.9144, categoryLength},
	"mile":           {1609.344, categoryLength},
	"miles":          {1609.344, categoryLength},
	"mi":             {1609.344, categoryLength},
	"nautical mile":  {1852, categoryLength},
	"nautical miles": {1852, categoryLength},
	"nmi":            {1852, categoryLength},

	// Weight — base: gram
	"mg":          {0.001, categoryWeight},
	"milligram":   {0.001, categoryWeight},
	"milligrams":  {0.001, categoryWeight},
	"g":           {1, categoryWeight},
	"gram":        {1, categoryWeight},
	"grams":       {1, categoryWeight},
	"kg":          {1000, categoryWeight},
	"kilogram":    {1000, categoryWeight},
	"kilograms":   {1000, categoryWeight},
	"tonne":       {1e6, categoryWeight},
	"tonnes":      {1e6, categoryWeight},
	"metric ton":  {1e6, categoryWeight},
	"metric tons": {1e6, categoryWeight},
	"lb":          {453.59237, categoryWeight},
	"lbs":         {453.59237, categoryWeight},
	"pound":       {453.59237, categoryWeight},
	"pounds":      {453.59237, categoryWeight},
	"oz":          {28.349523, categoryWeight},
	"ounce":       {28.349523, categoryWeight},
	"ounces":      {28.349523, categoryWeight},
	"stone":       {6350.2932, categoryWeight},
	"stones":      {6350.2932, categoryWeight},
	"st":          {6350.2932, categoryWeight},

	// Volume — base: litre
	"ml":           {0.001, categoryVolume},
	"milliliter":   {0.001, categoryVolume},
	"milliliters":  {0.001, categoryVolume},
	"millilitre":   {0.001, categoryVolume},
	"millilitres":  {0.001, categoryVolume},
	"l":            {1, categoryVolume},
	"liter":        {1, categoryVolume},
	"liters":       {1, categoryVolume},
	"litre":        {1, categoryVolume},
	"litres":       {1, categoryVolume},
	"fl oz":        {0.0295735, categoryVolume},
	"fl_oz":        {0.0295735, categoryVolume},
	"fluid ounce":  {0.0295735, categoryVolume},
	"fluid ounces": {0.0295735, categoryVolume},
	"cup":          {0.2365882, categoryVolume},
	"cups":         {0.2365882, categoryVolume},
	"pint":         {0.4731765, categoryVolume},
	"pints":        {0.4731765, categoryVolume},
	"pt":           {0.4731765, categoryVolume},
	"quart":        {0.9463529, categoryVolume},
	"quarts":       {0.9463529, categoryVolume},
	"qt":           {0.9463529, categoryVolume},
	"gallon":       {3.7854118, categoryVolume},
	"gallons":      {3.7854118, categoryVolume},
	"gal":          {3.7854118, categoryVolume},

	// Temperature — factors unused; handled by explicit functions
	"c":          {0, categoryTemperature},
	"celsius":    {0, categoryTemperature},
	"°c":         {0, categoryTemperature},
	"f":          {0, categoryTemperature},
	"fahrenheit": {0, categoryTemperature},
	"°f":         {0, categoryTemperature},
	"k":          {0, categoryTemperature},
	"kelvin":     {0, categoryTemperature},

	// Area — base: square metre
	"mm2":          {1e-6, categoryArea},
	"cm2":          {1e-4, categoryArea},
	"m2":           {1, categoryArea},
	"km2":          {1e6, categoryArea},
	"sqft":         {0.09290304, categoryArea},
	"sq ft":        {0.09290304, categoryArea},
	"square foot":  {0.09290304, categoryArea},
	"square feet":  {0.09290304, categoryArea},
	"sqmi":         {2589988.11, categoryArea},
	"sq mi":        {2589988.11, categoryArea},
	"square mile":  {2589988.11, categoryArea},
	"square miles": {2589988.11, categoryArea},
	"acre":         {4046.8564, categoryArea},
	"acres":        {4046.8564, categoryArea},
	"hectare":      {10000, categoryArea},
	"hectares":     {10000, categoryArea},
	"ha":           {10000, categoryArea},

	// Speed — base: metres per second
	"m/s":   {1, categorySpeed},
	"mps":   {1, categorySpeed},
	"km/h":  {1.0 / 3.6, categorySpeed},
	"kph":   {1.0 / 3.6, categorySpeed},
	"mph":   {0.44704, categorySpeed},
	"knot":  {0.514444, categorySpeed},
	"knots": {0.514444, categorySpeed},
	"ft/s":  {0.3048, categorySpeed},

	// Data — base: bit
	"bit":       {1, categoryData},
	"bits":      {1, categoryData},
	"byte":      {8, categoryData},
	"bytes":     {8, categoryData},
	"kb":        {8e3, categoryData},
	"kilobyte":  {8e3, categoryData},
	"kilobytes": {8e3, categoryData},
	"mb":        {8e6, categoryData},
	"megabyte":  {8e6, categoryData},
	"megabytes": {8e6, categoryData},
	"gb":        {8e9, categoryData},
	"gigabyte":  {8e9, categoryData},
	"gigabytes": {8e9, categoryData},
	"tb":        {8e12, categoryData},
	"terabyte":  {8e12, categoryData},
	"terabytes": {8e12, categoryData},

	// Time — base: second
	"ms":           {0.001, categoryTime},
	"millisecond":  {0.001, categoryTime},
	"milliseconds": {0.001, categoryTime},
	"s":            {1, categoryTime},
	"second":       {1, categoryTime},
	"seconds":      {1, categoryTime},
	"min":          {60, categoryTime},
	"minute":       {60, categoryTime},
	"minutes":      {60, categoryTime},
	"hr":           {3600, categoryTime},
	"hour":         {3600, categoryTime},
	"hours":        {3600, categoryTime},
	"day":          {86400, categoryTime},
	"days":         {86400, categoryTime},
	"week":         {604800, categoryTime},
	"weeks":        {604800, categoryTime},
}

// Service holds unit-conversion business logic.
type Service struct{}

// NewService creates a new conversion Service.
func NewService() *Service {
	return &Service{}
}

// Convert converts value from one unit to another and returns a Response.
func (s *Service) Convert(from, to string, value float64) (Response, error) {
	fromKey := strings.ToLower(strings.TrimSpace(from))
	toKey := strings.ToLower(strings.TrimSpace(to))

	fromInfo, ok := units[fromKey]
	if !ok {
		return Response{}, fmt.Errorf("unknown unit: %s", from)
	}

	toInfo, ok := units[toKey]
	if !ok {
		return Response{}, fmt.Errorf("unknown unit: %s", to)
	}

	if fromInfo.category != toInfo.category {
		return Response{}, fmt.Errorf("cannot convert %s to %s: incompatible unit types", from, to)
	}

	var result float64
	var formula string

	if fromInfo.category == categoryTemperature {
		result, formula = convertTemperature(fromKey, toKey, value)
	} else {
		factor := fromInfo.factor / toInfo.factor
		result = roundTo(value * factor)
		formula = fmt.Sprintf("%s × %s", from, formatFactor(factor))
	}

	return Response{
		From:    from,
		To:      to,
		Input:   value,
		Result:  result,
		Formula: formula,
	}, nil
}

// convertTemperature handles temperature conversions and returns the result and
// a human-readable formula string.
func convertTemperature(from, to string, value float64) (result float64, formula string) {
	// Normalise aliases to canonical names
	from = tempCanonical(from)
	to = tempCanonical(to)

	if from == to {
		return value, fmt.Sprintf("%s (no conversion needed)", from)
	}

	switch from + "→" + to {
	case "celsius→fahrenheit":
		return roundTo(value*9.0/5.0 + 32), "°C × 9/5 + 32"
	case "fahrenheit→celsius":
		return roundTo((value - 32) * 5.0 / 9.0), "(°F − 32) × 5/9"
	case "celsius→kelvin":
		return roundTo(value + 273.15), "°C + 273.15"
	case "kelvin→celsius":
		return roundTo(value - 273.15), "°K − 273.15"
	case "fahrenheit→kelvin":
		return roundTo((value-32)*5.0/9.0 + 273.15), "(°F − 32) × 5/9 + 273.15"
	case "kelvin→fahrenheit":
		return roundTo(value*9.0/5.0 - 459.67), "°K × 9/5 − 459.67"
	}

	return value, ""
}

func tempCanonical(s string) string {
	switch s {
	case "c", "celsius", "°c":
		return "celsius"
	case "f", "fahrenheit", "°f":
		return "fahrenheit"
	default:
		return "kelvin"
	}
}

// formatFactor formats a multiplier for the formula field, showing up to 6
// significant figures and stripping unnecessary trailing zeros.
func formatFactor(f float64) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.6f", f), "0"), ".")
}

// roundTo rounds v to at most 10 decimal places.
func roundTo(v float64) float64 {
	const pow = 1e10
	return math.Round(v*pow) / pow
}
