package transaction

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var credit = func(value float64) float64 {
	return value / 100
}

var debit = func(value float64) float64 {
	return -(value / 100)
}

type calcValueI func(value float64) float64

var calcValue = map[int]calcValueI{
	1: credit,
	2: credit,
	3: debit,
	4: credit,
}

func parseSeller(row string, tran *Model, errs *[]error) {
	s := row[66:86]
	caser := cases.Title(language.BrazilianPortuguese)
	tran.Seller = caser.String(strings.ToLower(strings.Trim(s, " ")))
}

func parseValue(row string, tran *Model, errs *[]error) {
	vStr := row[56:66]
	v, err := strconv.ParseFloat(vStr, 32)
	if err != nil {
		err := errors.New("the value of sale is not parsed")
		*errs = append(*errs, err)
		return
	}

	if fn, ok := calcValue[tran.Type]; ok {
		tran.Value = fn(v)
	}
}

func parseProduct(row string, tran *Model, errs *[]error) {
	p := row[26:56]
	caser := cases.Title(language.BrazilianPortuguese)
	tran.Product = caser.String(strings.ToLower(strings.Trim(p, " ")))
}

func parseDate(row string, tran *Model, errs *[]error) {
	dStr := row[1:26]
	layout := "2006-01-02T15:04:05-07:00"
	d, err := time.Parse(layout, dStr)
	if err != nil {
		err := errors.New("the date of sale is not parsed")
		*errs = append(*errs, err)
		return
	}
	tran.Date = d
}

func parseType(row string, tran *Model, errs *[]error) {
	tStr := row[:1]
	t, err := strconv.Atoi(tStr)
	if err != nil {
		err := errors.New("the type of sale is not parsed")
		*errs = append(*errs, err)
	}

	types := []int{1, 2, 3, 4}
	if !slices.Contains(types, t) {
		err := errors.New("unsupported type")
		*errs = append(*errs, err)
		return
	}

	tran.Type = t
}
