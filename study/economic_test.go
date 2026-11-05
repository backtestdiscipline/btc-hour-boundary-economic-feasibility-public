package study

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSyntheticHoursAndCosts(t *testing.T) {
	start := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	var candles []Candle
	for i := -1; i < 60; i++ {
		p := 100.
		if i >= 0 {
			p = 100 * math.Exp(.001)
		}
		candles = append(candles, Candle{Time: start.Add(time.Duration(i) * time.Minute), Open: 100, Close: p})
	}
	h, n, e := hours(candles, start, start.Add(time.Hour))
	if e != nil || n != 0 || len(h) != 1 {
		t.Fatal(h, n, e)
	}
	if math.Abs(h[0].cc-10) > 1e-9 || math.Abs(h[0].signed-10) > 1e-9 || h[0].other != 0 {
		t.Fatal(h)
	}
	h2 := append([]hour(nil), h...)
	h2 = append(h2, hour{t: start.AddDate(0, 0, 7), cc: 5, signed: -5})
	a, e := calculate(h2, 0, 0, 0)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := calculate(h2, 0, 0, 0)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("nondeterministic")
	}
	if math.Abs(a.Oracle.Mean-7.5) > 1e-9 || math.Abs(a.Signed.Mean-2.5) > 1e-9 || math.Abs(a.Scenarios[14].Net.Mean-(-16.5)) > 1e-9 {
		t.Fatal(a)
	}
	if a.Scenarios[14].Positive != 0 || a.Gate != "NOT_ECONOMICALLY_PROMISING" {
		t.Fatal(a)
	}
	candles = append(candles, candles[0])
	if _, _, e = hours(candles, start, start.Add(time.Hour)); e == nil {
		t.Fatal("duplicate accepted")
	}
	candles = candles[:len(candles)-2]
	if _, n, e = hours(candles, start, start.Add(time.Hour)); e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
func TestDecodeTimestampAndDefinitions(t *testing.T) {
	row := "1736121600000000,100,102,99,101,2,1736121659999999,200,1,1,100,0\n"
	c, _, e := DecodeDetailed(strings.NewReader(row), true)
	if e != nil || len(c) != 1 || c[0].Time.Minute() != 0 || c[0].Open != 100 || c[0].Close != 101 {
		t.Fatal(c, e)
	}
	_, _, e = DecodeDetailed(strings.NewReader(row+row), true)
	if e == nil {
		t.Fatal("duplicate accepted")
	}
	d := dist([]float64{0, 10, 20})
	if d.Median != 10 || d.Mean != 10 || d.Exceed[0] != 2./3 || d.Exceed[1] != 1./3 {
		t.Fatal(d)
	}
}
