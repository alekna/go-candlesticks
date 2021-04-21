package candlestick

import "time"

type Candle struct {
	Time   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}

type Candlestick struct {
	Candles          []*Candle
	Resolution       time.Duration
	TimeSeries       map[time.Time]*Candle
	LastCandle       *Candle
	CurrentCandle    *Candle
	CurrentCandleNew bool
	StartTime        time.Time
	EndTime          time.Time
}

func NewCandlestick(res time.Duration) *Candlestick {
	return &Candlestick{
		Resolution: res,
		Candles:    make([]*Candle, 0),
		TimeSeries: map[time.Time]*Candle{},
	}
}

func NewCandle(ti time.Time, value float64, volume float64) *Candle {
	return &Candle{
		Time:   ti,
		High:   value,
		Low:    value,
		Open:   value,
		Close:  value,
		Volume: volume,
	}
}

func (cs *Candlestick) AddCandle(candle *Candle) {
	cs.CurrentCandle = candle
	cs.Candles = append(cs.Candles, candle)
	cs.TimeSeries[candle.Time] = candle

	if candle.Time.Before(cs.StartTime) {
		cs.StartTime = candle.Time
	} else if candle.Time.After(cs.EndTime) {
		cs.EndTime = candle.Time
	}
}

func (cs *Candlestick) AddTrade(ti time.Time, value float64, volume float64) {
	var x = ti.Truncate(cs.Resolution)
	var candle = cs.TimeSeries[x]

	if candle != nil {
		candle.Add(value, volume)
		cs.CurrentCandleNew = false
	} else {
		candle = NewCandle(x, value, volume)
		cs.CurrentCandleNew = true
		cs.setLastCandle(candle)

		if cs.LastCandle != nil && x.After(cs.LastCandle.Time.Add(cs.Resolution)) {
			cs.backfill(candle.Time, cs.LastCandle.Close)
		}
		cs.AddCandle(candle)
	}
}

func (cs *Candlestick) backfill(x time.Time, value float64) {
	var flatCandle *Candle

	for ti := x; !ti.Equal(cs.LastCandle.Time); ti = ti.Add(-cs.Resolution) {
		if cs.TimeSeries[x] == nil {
			flatCandle = NewCandle(x, value, 0)
			cs.Candles = append(cs.Candles, flatCandle)
			cs.TimeSeries[x] = flatCandle
		}
	}
}

func (cs *Candlestick) setLastCandle(candle *Candle) {
	if cs.CurrentCandle == nil {
		cs.LastCandle = candle
	} else {
		cs.LastCandle = cs.CurrentCandle
	}
}

func (c *Candle) Add(value float64, volume float64) {
	if value > c.High {
		c.High = value
	} else if value < c.Low {
		c.Low = value
	}

	c.Volume += volume
	c.Close = value
}
