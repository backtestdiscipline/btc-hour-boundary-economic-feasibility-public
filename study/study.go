package study

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"path"
	"research-platform/seal"
	"sort"
	"strings"
	"time"
)

type Evaluation struct {
	Outcome seal.TerminalOutcome
	Payload []byte
}
type Verification struct {
	Status string `json:"status"`
	Notes  string `json:"notes"`
}
type Interval struct {
	Mean  float64 `json:"mean"`
	Lower float64 `json:"lower_95"`
	Upper float64 `json:"upper_95"`
}
type Distribution struct {
	Mean   float64   `json:"mean"`
	Median float64   `json:"median"`
	SD     float64   `json:"sample_sd"`
	P90    float64   `json:"p90"`
	P95    float64   `json:"p95"`
	P99    float64   `json:"p99"`
	Exceed []float64 `json:"fraction_above_5_10_15_20_30"`
}
type Scenario struct {
	Fee      float64  `json:"fee_per_side"`
	Extra    float64  `json:"additional_round_trip"`
	Net      Interval `json:"oracle_net"`
	Positive float64  `json:"fraction_positive"`
}
type Result struct {
	CloseCloseMean Interval     `json:"first_minute_close_close_mean_bps"`
	Version        string       `json:"version"`
	Rows           int          `json:"source_rows"`
	Defects        int          `json:"invalid_close_timestamp_rows"`
	Eligible       int          `json:"eligible_hours"`
	Excluded       int          `json:"excluded_hours"`
	Weeks          int          `json:"utc_weeks"`
	CloseClose     Distribution `json:"first_minute_close_close_bps"`
	OpenClose      Distribution `json:"first_minute_open_close_absolute_bps"`
	Other          Interval     `json:"other_minute_mean_close_close_bps"`
	Contrast       Interval     `json:"close_close_first_minus_other_bps"`
	Signed         Interval     `json:"signed_open_close_bps"`
	Short          Interval     `json:"short_open_close_bps"`
	Oracle         Interval     `json:"oracle_gross_bps"`
	Scenarios      []Scenario   `json:"scenarios"`
	Gate           string       `json:"gate"`
}
type hour struct {
	t                 time.Time
	cc, other, signed float64
}

func quant(x []float64, p float64) float64 {
	y := append([]float64(nil), x...)
	sort.Float64s(y)
	v := p * float64(len(y)-1)
	i := int(v)
	j := i + 1
	if j >= len(y) {
		return y[i]
	}
	return y[i] + (v-float64(i))*(y[j]-y[i])
}
func dist(x []float64) Distribution {
	d := Distribution{Exceed: make([]float64, 5)}
	for _, v := range x {
		d.Mean += v
		for j, t := range []float64{5, 10, 15, 20, 30} {
			if v > t {
				d.Exceed[j]++
			}
		}
	}
	d.Mean /= float64(len(x))
	for _, v := range x {
		d.SD += (v - d.Mean) * (v - d.Mean)
	}
	if len(x) > 1 {
		d.SD = math.Sqrt(d.SD / float64(len(x)-1))
	}
	for j := range d.Exceed {
		d.Exceed[j] /= float64(len(x))
	}
	d.Median = quant(x, .5)
	d.P90 = quant(x, .9)
	d.P95 = quant(x, .95)
	d.P99 = quant(x, .99)
	return d
}
func hours(c []Candle, start, end time.Time) ([]hour, int, error) {
	m := make(map[int64]Candle, len(c))
	for _, v := range c {
		k := v.Time.Unix()
		if _, ok := m[k]; ok {
			return nil, 0, fmt.Errorf("duplicate candle")
		}
		m[k] = v
	}
	var out []hour
	excluded := 0
	for t := start; t.Before(end); t = t.Add(time.Hour) {
		prev, ok := m[t.Unix()-60]
		if !ok {
			excluded++
			continue
		}
		var vals [60]Candle
		valid := true
		for j := range vals {
			vals[j], ok = m[t.Unix()+int64(j)*60]
			if !ok {
				valid = false
				break
			}
		}
		if !valid {
			excluded++
			continue
		}
		h := hour{t: t, cc: math.Abs(math.Log(vals[0].Close/prev.Close)) * 10000, signed: math.Log(vals[0].Close/vals[0].Open) * 10000}
		for j := 1; j < 60; j++ {
			h.other += math.Abs(math.Log(vals[j].Close/vals[j-1].Close)) * 10000 / 59
		}
		out = append(out, h)
	}
	return out, excluded, nil
}
func calculate(h []hour, excluded, rows, defects int) (Result, error) {
	if len(h) == 0 {
		return Result{}, fmt.Errorf("no eligible hours")
	}
	var cc, oc, signed, other, contrast []float64
	type cluster struct {
		n    int
		sums [5]float64
	}
	groups := map[int64]*cluster{}
	var keys []int64
	for _, v := range h {
		a := math.Abs(v.signed)
		cc = append(cc, v.cc)
		oc = append(oc, a)
		signed = append(signed, v.signed)
		other = append(other, v.other)
		contrast = append(contrast, v.cc-v.other)
		d := v.t.UTC()
		monday := d.AddDate(0, 0, -(int(d.Weekday())+6)%7)
		k := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC).Unix()
		g := groups[k]
		if g == nil {
			g = &cluster{}
			groups[k] = g
			keys = append(keys, k)
		}
		g.n++
		g.sums[0] += a
		g.sums[1] += v.signed
		g.sums[2] += v.other
		g.sums[3] += v.cc - v.other
		g.sums[4] += v.cc
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	if len(keys) < 2 {
		return Result{}, fmt.Errorf("fewer than two weeks")
	}
	var reps [5][]float64
	rng := rand.New(rand.NewSource(20261105))
	for r := 0; r < 10000; r++ {
		var sum [5]float64
		n := 0
		for range keys {
			g := groups[keys[rng.Intn(len(keys))]]
			n += g.n
			for j := range sum {
				sum[j] += g.sums[j]
			}
		}
		for j := range sum {
			reps[j] = append(reps[j], sum[j]/float64(n))
		}
	}
	mean := func(x []float64) float64 {
		s := 0.
		for _, v := range x {
			s += v
		}
		return s / float64(len(x))
	}
	ci := func(x []float64, j int) Interval {
		return Interval{mean(x), quant(reps[j], .025), quant(reps[j], .975)}
	}
	res := Result{Version: "economic-feasibility/v1", Rows: rows, Defects: defects, Eligible: len(h), Excluded: excluded, Weeks: len(keys), CloseClose: dist(cc), OpenClose: dist(oc), Other: ci(other, 2), Contrast: ci(contrast, 3), Signed: ci(signed, 1), Oracle: ci(oc, 0)}
	res.CloseCloseMean = ci(cc, 4)
	res.Short = Interval{-res.Signed.Mean, -res.Signed.Upper, -res.Signed.Lower}
	for _, fee := range []float64{0, 2, 5, 10} {
		for _, extra := range []float64{0, 2, 4, 8} {
			cost := 2*fee + extra
			n := 0
			for _, v := range oc {
				if v > cost {
					n++
				}
			}
			res.Scenarios = append(res.Scenarios, Scenario{fee, extra, Interval{res.Oracle.Mean - cost, res.Oracle.Lower - cost, res.Oracle.Upper - cost}, float64(n) / float64(len(oc))})
		}
	}
	primary := res.Scenarios[14].Net
	res.Gate = "INCONCLUSIVE"
	if primary.Lower > 0 {
		res.Gate = "WORTH_FOLLOWUP"
	} else if primary.Upper < 0 {
		res.Gate = "NOT_ECONOMICALLY_PROMISING"
	}
	return res, nil
}
func Evaluate(in Inputs) (Evaluation, error) {
	fail := func(e error) (Evaluation, error) {
		b, _ := json.Marshal(map[string]string{"error": e.Error()})
		return Evaluation{seal.TechnicalFailure, b}, nil
	}
	refs := append([]seal.ArtifactRef(nil), in.Sources...)
	sort.Slice(refs, func(i, j int) bool { return refs[i].Path < refs[j].Path })
	if len(refs) != 21 {
		return fail(fmt.Errorf("expected 21 archives"))
	}
	var all []Candle
	rows, defects := 0, 0
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start
	for _, ref := range refs {
		name := path.Base(ref.Path)
		expected := "BTCUSDT-1m-" + end.Format("2006-01") + ".zip"
		if name != expected {
			return fail(fmt.Errorf("unexpected archive %s", name))
		}
		c, d, e := ReadArchive(in, ref, true)
		if e != nil {
			return fail(e)
		}
		next := end.AddDate(0, 1, 0)
		if c[0].Time.Before(end) || !c[len(c)-1].Time.Before(next) {
			return fail(fmt.Errorf("out of month"))
		}
		all = append(all, c...)
		rows += d.Total
		defects += len(d.Excluded)
		end = next
	}
	h, excluded, e := hours(all, start, end)
	if e != nil {
		return fail(e)
	}
	res, e := calculate(h, excluded, rows, defects)
	if e != nil {
		return fail(e)
	}
	b, e := json.Marshal(res)
	return Evaluation{seal.Complete, b}, e
}
func Verify(in Inputs, result []byte) (Verification, error) {
	e, err := Evaluate(in)
	if err != nil {
		return Verification{}, err
	}
	if e.Outcome != seal.Complete || strings.TrimSpace(string(e.Payload)) != strings.TrimSpace(string(result)) {
		return Verification{"FAIL", "Source replay differs from result"}, nil
	}
	return Verification{"PASS", "Deterministic semantic replay of bound sources matches"}, nil
}
