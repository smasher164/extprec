// MIT License

// Copyright (c) 2018 Akhil Indurti

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package extprec

import . "math/bits"

// Add returns the sum and carry-out bit (0 or 1) of
// x, y, and a carry-in bit (assumed to be 0 or 1).
func Add(x, y, carry uint) (sum, carryOut uint) {
	if UintSize == 32 {
		s32, c32 := Add32(uint32(x), uint32(y), uint32(carry))
		return uint(s32), uint(c32)
	}
	s64, c64 := Add64(uint64(x), uint64(y), uint64(carry))
	return uint(s64), uint(c64)
}

// Add32 returns the 32-bit sum and carry-out bit (0 or 1)
// of x, y, and a carry-in bit (assumed to be 0 or 1).
func Add32(x, y, carry uint32) (sum, carryOut uint32) {
	sum = uint32(uint64(x) + uint64(y) + uint64(carry))
	carryOut = ((x & y) | ((x | y) & ^sum)) >> 31
	return
}

// Add64 returns the 64-bit sum and carry-out bit (0 or 1)
// of x, y, and a carry-in bit (assumed to be 0 or 1).
func Add64(x, y, carry uint64) (sum, carryOut uint64) {
	// See "Hacker's Delight", Ch. 2-16: Double Length Add/Subtract
	//
	// Extract the higher and lower halves of x and y.
	// Take the sum of the lower halves and the given carry.
	// Compute the next carry (Ch. 2-13: Overflow Detection).
	// Take the sum of the higher halves and the new carry.
	// Compute carryOut in the same way.

	xlo := x & 0xFFFFFFFF
	xhi := x >> 32
	ylo := y & 0xFFFFFFFF
	yhi := y >> 32

	lo := (xlo + ylo + carry) & 0xFFFFFFFF
	c := ((xlo & ylo) | ((xlo | ylo) & ^lo)) >> 31
	hi := (xhi + yhi + c) & 0xFFFFFFFF
	carryOut = ((xhi & yhi) | ((xhi | yhi) & ^hi)) >> 31
	sum = (hi << 32) | lo
	return
}

// Sub returns the difference and borrow-out bit (0 or 1) of
// x, y, and a borrow-in bit (assumed to be 0 or 1).
func Sub(x, y, borrow uint) (difference, borrowOut uint) {
	if UintSize == 32 {
		d32, b32 := Sub32(uint32(x), uint32(y), uint32(borrow))
		return uint(d32), uint(b32)
	}
	d64, b64 := Sub64(uint64(x), uint64(y), uint64(borrow))
	return uint(d64), uint(b64)
}

// Sub32 returns the 32-bit difference and borrow-out bit (0 or 1) of
// x, y, and a borrow-in bit (assumed to be 0 or 1).
func Sub32(x, y, borrow uint32) (difference, borrowOut uint32) {
	difference = x - y - borrow
	borrowOut = ((^x & y) | (^(x ^ y) & difference)) >> 31
	return
}

// Sub64 returns the 64-bit difference and borrow-out bit (0 or 1) of
// x, y, and a borrow-in bit (assumed to be 0 or 1).
func Sub64(x, y, borrow uint64) (difference, borrowOut uint64) {
	// See "Hacker's Delight", Ch. 2-16: Double Length Add/Subtract
	//
	// Extract the higher and lower halves of x and y.
	// Take the difference of the lower halves and the given borrow.
	// Compute the next borrow (Ch. 2-13: Overflow Detection).
	// Take the difference of the higher halves and the new borrow.
	// Compute borrowOut in the same way.

	xlo := x & 0xFFFFFFFF
	xhi := x >> 32
	ylo := y & 0xFFFFFFFF
	yhi := y >> 32

	lo := (xlo - ylo - borrow) & 0xFFFFFFFF
	b := ((^xlo & ylo) | (^(xlo ^ ylo) & lo)) >> 31
	hi := (xhi - yhi - b) & 0xFFFFFFFF
	borrowOut = ((^xhi & yhi) | (^(xhi ^ yhi) & hi)) >> 31
	difference = (hi << 32) | lo
	return
}

// Mul returns the most significant and least significant
// bits from the product of x and y.
func Mul(x, y uint) (hi, lo uint) {
	if UintSize == 32 {
		hi32, lo32 := Mul32(uint32(x), uint32(y))
		return uint(hi32), uint(lo32)
	}
	hi64, lo64 := Mul64(uint64(x), uint64(y))
	return uint(hi64), uint(lo64)
}

// Mul32 returns the most significant and least significant
// 32 bits from the 64-bit product of x and y.
func Mul32(x, y uint32) (hi, lo uint32) {
	z := uint64(x) * uint64(y)
	return uint32(z >> 32), uint32(z)
}

// Mul64 returns the most significant and least significant
// 64 bits from the 128-bit product of x and y.
func Mul64(x, y uint64) (hi, lo uint64) {
	// See http://www.hackersdelight.org/MontgomeryMultiplication.pdf
	//
	// Montgomery Multiplication for fixed-width values:
	// Extract the higher and lower halves of x and y.
	// Product = x*y*r mod m, calculated more efficiently
	// with just three multiplications.

	xlo := x & 0xFFFFFFFF
	xhi := x >> 32
	ylo := y & 0xFFFFFFFF
	yhi := y >> 32

	t := xlo * ylo
	w0 := t & 0xFFFFFFFF
	k := t >> 32

	t = xhi*ylo + k
	w1 := t & 0xFFFFFFFF
	w2 := t >> 32

	t = xlo*yhi + w1
	k = t >> 32

	lo = (t << 32) + w0
	hi = xhi*yhi + w2 + k
	return
}

// Div returns the quotient and remainder of (hi || lo) and x,
// where hi and lo hold the most significant and least
// significant bits of the dividend.
func Div(hi, lo, x uint) (quo, rem uint) {
	if UintSize == 32 {
		q32, r32 := Div32(uint32(hi), uint32(lo), uint32(x))
		return uint(q32), uint(r32)
	}
	q64, r64 := Div64(uint64(hi), uint64(lo), uint64(x))
	return uint(q64), uint(r64)
}

// Div32 returns the 32-bit quotient and remainder of (hi || lo) and x,
// where hi and lo hold the most significant and least
// significant 32 bits of the dividend.
func Div32(hi, lo, x uint32) (quo, rem uint32) {
	y := (uint64(hi) << 32) | uint64(lo)
	quo = uint32(y / uint64(x))
	rem = uint32(y % uint64(x))
	return
}

// Div64 returns the 64-bit quotient and remainder of (hi || lo) and x,
// where hi and lo hold the most significant and least
// significant 64 bits of the dividend.
func Div64(hi, lo, x uint64) (quo, rem uint64) {
	// See "Hacker's Delight", Ch. 9-4: Unsigned Long Division
	const b = 1 << 32
	s := LeadingZeros64(x)
	us := uint64(s)
	// Normalize divisor.
	x = x << us
	// Break divisor up into two 32-bit digits.
	xhi := x >> 32
	xlo := x & 0xFFFFFFFF

	un64 := (hi << us) | ((lo >> (64 - us)) & uint64(-s>>31))
	// Shift dividend left.
	un10 := lo << us
	// Break right half of dividend into two digits.
	un1 := un10 >> 32
	un0 := un10 & 0xFFFFFFFF

	// Compute the first quotient digit, q1.
	q1 := un64 / xhi
	rhat := un64 - q1*xhi
	for q1 >= b || q1*xlo > b*rhat+un1 {
		q1 -= 1
		rhat += xhi
		if rhat >= b {
			break
		}
	}
	// Multiply and subtract
	un21 := un64*b + un1 - q1*x

	// Compute the second quotient digit, q0.
	q0 := un21 / xhi
	rhat = un21 - q0*xhi
	for q0 >= b || q0*xlo > b*rhat+un0 {
		q0 -= 1
		rhat += xhi
		if rhat >= b {
			break
		}
	}
	quo = q1*b + q0
	rem = (un21*b + un0 - q0*x) >> us
	return
}
