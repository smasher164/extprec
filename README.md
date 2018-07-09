This package provides a portable implementation of extended precision integer arithmetic operations in Go. Specifically, it implements the functions outlined for the math/bits package in [#24813](golang.org/issue/24813) as follows:
```
// Add with carry
// The carry inputs are assumed to be 0 or 1; otherwise behavior undefined.
// The carryOut outputs are guaranteed to be 0 or 1.
func Add(x, y, carry uint) (sum, carryOut uint)
func Add32(x, y, carry uint32) (sum, carryOut uint32)
func Add64(x, y, carry uint64) (sum, carryOut uint64)

// Subtract with borrow
// The borrow inputs are assumed to be 0 or 1; otherwise behavior undefined.
// The borrowOut outputs are guaranteed to be 0 or 1.
func Sub(x, y, borrow uint) (difference, borrowOut uint)
func Sub32(x, y, borrow uint32) (difference, borrowOut uint32)
func Sub64(x, y, borrow uint64) (difference, borrowOut uint64)

// Full-width multiply: 32x32->64, 64x64->128
func Mul(x, y uint) (hi, lo uint)
func Mul32(x, y uint32) (hi, lo uint32)
func Mul64(x, y uint64) (hi, lo uint64)

// Full-width divide: 64/32 -> 32,32, 128/64 -> 64,64
// Behavior undefined if hi >= x (because quotient will not fit).
func Div(hi, lo, x uint) (quo, rem uint)
func Div32(hi, lo, x uint32) (quo, rem uint32)
func Div64(hi, lo, x uint64) (quo, rem uint64)
```

I still need to add unit tests and hammer out bugs in the library, so feel free to file an issue!