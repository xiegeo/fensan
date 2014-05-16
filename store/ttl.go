package store

import "time"

const ttl_base_year = 2000

//A TTL encodes the number of periods from a base time.
//For ease of local usage, it is always per month, from January, 2000.
//
//For more information, see fensan/docs
type TTL struct {
	month int16 //January, 2000 is 0.
}

var ttl_now = int16(0)

//Return a TTL for the current period (month), cached. This is a TTL of 0.
//
func TTLNow() TTL {
	if ttl_now == 0 {
		y, m, _ := time.Now().UTC().Date()
		tn := int16(y - ttl_base_year*12 + (int(m) - 1))
		ttl_now = tn
		go func() {
			time.Sleep(time.Second)
			ttl_now = 0
		}()
		return TTL{tn}
	}
	cached := TTL{ttl_now}
	if cached.month == 0 {
		//retry
		return TTLNow()
	}
	return cached
}

func TTLFromBytes(h, l byte) TTL {
	return TTL{int16(h)<<8 + int16(l)}
}

//Bytes returns the 2 bytes of the TTL for serialization
func (t TTL) Bytes() (h, l byte) {
	return byte(t.month >> 8), byte(t.month & 0xff)
}

func (t TTL) YearMonth() (year int, month time.Month) {
	y, m := t.month/12, t.month%12
	return int(y) + ttl_base_year, time.Month(m) + time.January
}

//End returns the time at the end of TTL, just before the 1st of next month.
func (t TTL) End() time.Time {
	y, m := t.YearMonth()
	return time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
}
