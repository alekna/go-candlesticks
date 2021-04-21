package candlestick

import (
	"testing"
	"time"
)

func TestCandlestick(t *testing.T) {
	var c = NewCandlestick(time.Minute)
	var start = time.Date(2009, time.November, 10, 23, 30, 5, 0, time.UTC)

	c.AddTrade(start, 5, 1)
	c.AddTrade(start.Add(5*time.Second), 25, 1)
	c.AddTrade(start.Add(25*time.Second), 3, 1)
	var c1 = c.Candles[0]

	c.AddTrade(start.Add(60*time.Second), 12, 5)
	c.AddTrade(start.Add(70*time.Second), 13, 2)
	var c2 = c.Candles[1]

	// Intentionally empty data series included here, to test flat candles
	c.AddTrade(start.Add(240*time.Second), 15, 6)
	var c3 = c.Candles[2]
	var c4 = c.Candles[3]

	if !(c1.Volume == 3 && c1.Open == 5 && c1.Close == 3 &&
		c1.High == 25 && c1.Low == 3) {
		t.Logf("Got wrong val: %v", c1)
		t.Fail()
	}

	if !(c2.Volume == 7 && c2.Open == 12 && c2.Close == 13 &&
		c2.High == 13 && c2.Low == 12) {
		t.Logf("Got wrong val: %v", c2)
		t.Fail()
	}

	if !(c3.Volume == 0 && c3.Open == 13 && c3.Close == 13 &&
		c3.High == 13 && c3.Low == 13) {
		t.Logf("Got wrong val: %v", c3)
		t.Fail()
	}

	if !(c4.Volume == 6 && c4.Open == 15 && c4.Close == 15 &&
		c4.High == 15 && c4.Low == 15) {
		t.Logf("Got wrong val: %v", c4)
		t.Fail()
	}

	if len(c.Candles) != 4 {
		t.Log("Got wrong len")
		t.Fail()
	}
}
