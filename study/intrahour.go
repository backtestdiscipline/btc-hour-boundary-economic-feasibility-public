package study

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
)

// Candle holds a UTC minute-open time and its close, after validation.
type Candle struct {
	Time  time.Time
	Close float64
	Open  float64
}

// DecodeCandles validates a headerless Binance spot kline CSV in a declared timestamp unit.
type RowDefect struct {
	Row      int       `json:"row"`
	Open     string    `json:"raw_open"`
	Close    string    `json:"raw_close"`
	Expected int64     `json:"expected_close"`
	Time     time.Time `json:"utc_open"`
	Reason   string    `json:"reason"`
}
type DecodeReport struct {
	Total    int         `json:"total_rows"`
	Accepted int         `json:"accepted_rows"`
	Excluded []RowDefect `json:"excluded_rows"`
}

func DecodeCandles(r io.Reader, micros bool) ([]Candle, error) {
	c, d, e := DecodeDetailed(r, micros)
	if e == nil && len(d.Excluded) > 0 {
		return nil, fmt.Errorf("excluded %d rows; use DecodeDetailed", len(d.Excluded))
	}
	return c, e
}
func DecodeDetailed(r io.Reader, micros bool) ([]Candle, DecodeReport, error) {
	c := csv.NewReader(r)
	c.FieldsPerRecord = 12
	unit := int64(1000)
	if micros {
		unit = 1000000
	}
	var out []Candle
	var report DecodeReport
	var previous time.Time
	for n := 1; ; n++ {
		row, err := c.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, report, fmt.Errorf("row %d: %w", n, err)
		}
		report.Total++
		if report.Total > 50000 {
			return nil, report, fmt.Errorf("row limit exceeded")
		}
		op, e := strconv.ParseInt(row[0], 10, 64)
		if e != nil || op < 0 || op%(60*unit) != 0 {
			return nil, report, fmt.Errorf("row %d: invalid open time", n)
		}
		cl, closeErr := strconv.ParseInt(row[6], 10, 64)
		vals := make([]float64, 8)
		for j, k := range []int{1, 2, 3, 4, 5, 7, 9, 10} {
			v, e := strconv.ParseFloat(row[k], 64)
			if e != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
				return nil, report, fmt.Errorf("row %d: invalid number", n)
			}
			vals[j] = v
		}
		o, h, l, close := vals[0], vals[1], vals[2], vals[3]
		if l <= 0 || o < l || close < l || h < o || h < close || h < l || vals[6] > vals[4] || vals[7] > vals[5] {
			return nil, report, fmt.Errorf("row %d: invalid OHLCV", n)
		}
		if _, e = strconv.ParseUint(row[8], 10, 64); e != nil {
			return nil, report, fmt.Errorf("row %d: invalid trades", n)
		}
		if row[11] != "0" {
			return nil, report, fmt.Errorf("row %d: unexpected final field", n)
		}
		t := time.Unix(op/unit, 0).UTC()
		if !previous.IsZero() && !t.After(previous) {
			return nil, report, fmt.Errorf("row %d: duplicate or unordered timestamp", n)
		}
		previous = t
		if closeErr != nil || cl != op+60*unit-1 {
			if len(report.Excluded) >= 1000 {
				return nil, report, fmt.Errorf("too many defective rows")
			}
			report.Excluded = append(report.Excluded, RowDefect{n, row[0], row[6], op + 60*unit - 1, t, "invalid_close_timestamp"})
			continue
		}
		out = append(out, Candle{Time: t, Close: close, Open: o})
		report.Accepted++
	}
	if len(out) == 0 {
		return nil, report, fmt.Errorf("empty CSV")
	}
	return out, report, nil
}

type HourProfile struct {
	Start  time.Time
	Values [60]float64
}

// Profiles requires the preceding candle and omits incomplete hours and zero denominators.
func Profiles(c []Candle) ([]HourProfile, int, error) {
	var out []HourProfile
	zeros := 0
	for i, v := range c {
		if v.Close <= 0 || math.IsNaN(v.Close) || math.IsInf(v.Close, 0) || v.Time.Second() != 0 || v.Time.Nanosecond() != 0 {
			return nil, 0, fmt.Errorf("invalid candle")
		}
		if i > 0 && !v.Time.After(c[i-1].Time) {
			return nil, 0, fmt.Errorf("duplicate or unordered candle")
		}
	}
	for i := 1; i+59 < len(c); i++ {
		if c[i].Time.Minute() != 0 || !c[i].Time.Equal(c[i-1].Time.Add(time.Minute)) {
			continue
		}
		complete := true
		for j := 1; j < 60; j++ {
			if !c[i+j].Time.Equal(c[i].Time.Add(time.Duration(j) * time.Minute)) {
				complete = false
				break
			}
		}
		if !complete {
			continue
		}
		var a [60]float64
		sum := 0.0
		for j := 0; j < 60; j++ {
			a[j] = math.Abs(math.Log(c[i+j].Close) - math.Log(c[i+j-1].Close))
			sum += a[j]
		}
		if sum == 0 {
			zeros++
			continue
		}
		h := HourProfile{Start: c[i].Time}
		for j := range a {
			h.Values[j] = 60 * a[j] / sum
		}
		out = append(out, h)
	}
	return out, zeros, nil
}
