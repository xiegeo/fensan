// Adopted from Go crypto/sha256
// To use sha2 as a compression function (512 to 256)
// It is used as the inner hash of the hashtree, using
// sha2-224 init vectors to be a different hash function
// from the leaf hash, which uses sha2-256 unmodified.
//
// Origin carried the following notices:
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SHA256 block step.
// In its own file so that a faster assembly or C version
// can be substituted easily.

package hashtree

const (
	init0_224 = 0xC1059ED8
	init1_224 = 0x367CD507
	init2_224 = 0x3070DD17
	init3_224 = 0xF70E5939
	init4_224 = 0xFFC00B31
	init5_224 = 0x68581511
	init6_224 = 0x64F98FA7
	init7_224 = 0xBEFA4FA4
)

var _K = []uint32{
	0x428a2f98,
	0x71374491,
	0xb5c0fbcf,
	0xe9b5dba5,
	0x3956c25b,
	0x59f111f1,
	0x923f82a4,
	0xab1c5ed5,
	0xd807aa98,
	0x12835b01,
	0x243185be,
	0x550c7dc3,
	0x72be5d74,
	0x80deb1fe,
	0x9bdc06a7,
	0xc19bf174,
	0xe49b69c1,
	0xefbe4786,
	0x0fc19dc6,
	0x240ca1cc,
	0x2de92c6f,
	0x4a7484aa,
	0x5cb0a9dc,
	0x76f988da,
	0x983e5152,
	0xa831c66d,
	0xb00327c8,
	0xbf597fc7,
	0xc6e00bf3,
	0xd5a79147,
	0x06ca6351,
	0x14292967,
	0x27b70a85,
	0x2e1b2138,
	0x4d2c6dfc,
	0x53380d13,
	0x650a7354,
	0x766a0abb,
	0x81c2c92e,
	0x92722c85,
	0xa2bfe8a1,
	0xa81a664b,
	0xc24b8b70,
	0xc76c51a3,
	0xd192e819,
	0xd6990624,
	0xf40e3585,
	0x106aa070,
	0x19a4c116,
	0x1e376c08,
	0x2748774c,
	0x34b0bcb5,
	0x391c0cb3,
	0x4ed8aa4a,
	0x5b9cca4f,
	0x682e6ff3,
	0x748f82ee,
	0x78a5636f,
	0x84c87814,
	0x8cc70208,
	0x90befffa,
	0xa4506ceb,
	0xbef9a3f7,
	0xc67178f2,
}

func ht_sha256block(left, right *H256) *H256 {
	var w [64]uint32
	var h0, h1, h2, h3, h4, h5, h6, h7 uint32
	h0, h1, h2, h3, h4, h5, h6, h7 = init0_224, init1_224, init2_224, init3_224, init4_224, init5_224, init6_224, init7_224
	for i := 0; i < 8; i++ {
		w[i] = left[i]
	}
	for i := 0; i < 8; i++ {
		w[i+8] = right[i]
	}
	for i := 16; i < 64; i++ {
		t1 := (w[i-2]>>17 | w[i-2]<<(32-17)) ^ (w[i-2]>>19 | w[i-2]<<(32-19)) ^ (w[i-2] >> 10)

		t2 := (w[i-15]>>7 | w[i-15]<<(32-7)) ^ (w[i-15]>>18 | w[i-15]<<(32-18)) ^ (w[i-15] >> 3)

		w[i] = t1 + w[i-7] + t2 + w[i-16]
	}

	a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7

	for i := 0; i < 64; i++ {
		t1 := h + ((e>>6 | e<<(32-6)) ^ (e>>11 | e<<(32-11)) ^ (e>>25 | e<<(32-25))) + ((e & f) ^ (^e & g)) + _K[i] + w[i]

		t2 := ((a>>2 | a<<(32-2)) ^ (a>>13 | a<<(32-13)) ^ (a>>22 | a<<(32-22))) + ((a & b) ^ (a & c) ^ (b & c))

		h = g
		g = f
		f = e
		e = d + t1
		d = c
		c = b
		b = a
		a = t1 + t2
	}

	h0 += a
	h1 += b
	h2 += c
	h3 += d
	h4 += e
	h5 += f
	h6 += g
	h7 += h

	return &H256{h0, h1, h2, h3, h4, h5, h6, h7}
}
