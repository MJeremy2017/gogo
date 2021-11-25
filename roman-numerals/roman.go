package roman

import (
	"strings"
)

type RomanNumeral struct {
	value int
	symbol string
}

type RomanNumerals []RomanNumeral

var allRomanNumerals = RomanNumerals {
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func (r RomanNumerals) ValueOf(symbols ...byte) int {
	symbol := string(symbols)
	for _, v := range r {
		if symbol == v.symbol {
			return v.value
		}
	}
	return 0
}

func ConvertToRoman(value int) string {
	var result strings.Builder

	for _, numeral := range allRomanNumerals {
		for value >= numeral.value {
			result.WriteString(numeral.symbol)
			value -= numeral.value
		}
	}

	return result.String()
}

func ConvertToNumeral(roman string) int {
	total := 0
	for i := 0; i < len(roman); i++ {
		symbol := roman[i]
		if couldBeSubtractive(i, symbol, roman) {
			if value := allRomanNumerals.ValueOf(symbol, roman[i+1]); value != 0 {
				total += value
				i++
			} else {
				total += allRomanNumerals.ValueOf(symbol)
			}
		} else {
			total += allRomanNumerals.ValueOf(symbol)
		}
	}
	return total
}

func couldBeSubtractive(index int, currentSymbol uint8, roman string) bool {
	isSubtractiveSymbol := currentSymbol == 'I' || currentSymbol == 'X' || currentSymbol == 'C'
	return index+1 < len(roman) && isSubtractiveSymbol
}


