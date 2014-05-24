package store

import "time"

const ttl_base_year = 2000

//A TTL encodes the number of periods from a base time.
//For ease of local usage, it is always per month, from January, 2000.
//
//For more information, see fensan/docs
type TTL int16 //January, 2000 is 0.

const (
	TTLLongAgo = TTL(-1)
)

var ttl_now = int16(0)

//Return a TTL for the current period (month), cached. This is a TTL of 0.
//
func TTLNow() TTL {
	cached := TTL(ttl_now)
	if cached == 0 {
		y, m, _ := time.Now().UTC().Date()
		tn := int16(y - ttl_base_year*12 + (int(m) - 1))
		ttl_now = tn
		go func() {
			time.Sleep(time.Second)
			ttl_now = 0
		}()
		return TTL(tn)
	}
	return cached
}

func TTLFromBytes(b []byte) TTL {
	return TTL(int16(b[0]) + int16(b[1])<<8)
}

//Bytes returns the 2 bytes of the TTL for serialization
func (t TTL) Bytes() []byte {
	return []byte{byte(t & 0xff), byte(t >> 8)}
}

func (t TTL) YearMonth() (year int, month time.Month) {
	y, m := t/12, t%12
	return int(y) + ttl_base_year, time.Month(m) + time.January
}

//End returns the time at the end of TTL, just before the 1st of next month.
func (t TTL) End() time.Time {
	y, m := t.YearMonth()
	return time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
}

//MonthUntil is the number of months until to, ie: 3.MonthUntil(5) = 2
func (t TTL) MonthUntil(to TTL) int16 {
	return int16(to - t)
}
