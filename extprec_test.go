package extprec

import (
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("edge", edgeAdd)
	t.Run("carryOut==0", carryOutZero)
	t.Run("carryOut==1,32-bit", carryOutOne32)
	t.Run("carryOut==1,64-bit", carryOutOne64)
}

func edgeAdd(t *testing.T) {
	// Edge cases with extremely large and small values
	edge := []struct {
		x, y, carry, sum, carryOut uint64
	}{
		{0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0},
		{0, 1<<64 - 1, 0, 1<<64 - 1, 0},
		{1<<64 - 1, 0, 1, 0, 1},
		{0xAAAAAAAAAAAAAAAA, 0x5555555555555555, 0, 1<<64 - 1, 0},
		{0xAAAAAAAAAAAAAAAA, 0x5555555555555555, 1, 0, 1},
	}
	for _, c := range edge {
		s64, c64 := Add64(c.x, c.y, c.carry)
		if s64 != c.sum || c64 != c.carryOut {
			t.Errorf("Add64(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
				c.x, c.y, c.carry,
				s64, c64,
				c.sum, c.carryOut)
		}
		s32, c32 := Add32(uint32(c.x), uint32(c.y), uint32(c.carry))
		if s32 != uint32(c.sum) || c32 != uint32(c.carryOut) {
			t.Errorf("Add32(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
				uint32(c.x), uint32(c.y), uint32(c.carry),
				s32, c32,
				uint32(c.sum), uint32(c.carryOut))
		}
		sum, carryOut := Add(uint(c.x), uint(c.y), uint(c.carry))
		if sum != uint(c.sum) || carryOut != uint(c.carryOut) {
			t.Errorf("Add(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
				uint(c.x), uint(c.y), uint(c.carry),
				sum, carryOut,
				uint(c.sum), uint(c.carryOut))
		}
	}
}

func carryOutZero(t *testing.T) {
	// Sums where carryOut == 0
	const (
		valid64  = (1<<64 - 1) / 2
		valid32  = (1<<32 - 1) / 2
		interval = valid64 / 256
	)
	for x := uint64(0); x < valid64; x += interval {
		for y := uint64(0); y < valid64; y += interval {
			if x < valid32 && y < valid32 {
				if sum, carryOut := Add32(uint32(x), uint32(y), 0); sum != uint32(x+y+0) || carryOut != 0 {
					t.Errorf("Add32(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
						uint32(x), uint32(y), 0,
						sum, carryOut,
						uint32(x+y+0), 0)
				}
				if sum, carryOut := Add32(uint32(x), uint32(y), 1); sum != uint32(x+y+1) || carryOut != 0 {
					t.Errorf("Add32(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
						uint32(x), uint32(y), 1,
						sum, carryOut,
						uint32(x+y+1), 0)
				}
			}
			if sum, carryOut := Add64(x, y, 0); sum != x+y+0 || carryOut != 0 {
				t.Errorf("Add64(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 0,
					sum, carryOut,
					x+y+0, 0)
			}
			if sum, carryOut := Add64(x, y, 1); sum != x+y+1 || carryOut != 0 {
				t.Errorf("Add64(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 1,
					sum, carryOut,
					x+y+1, 0)
			}
		}
	}
}

func carryOutOne32(t *testing.T) {
	// Sums where carryOut == 1
	// 32-bit
	const (
		halfRange = (1 << 32) / 2
		valid32   = 1<<32 - 1
		interval  = halfRange / 256
	)
	for x := uint32(valid32 - interval + 1); x >= halfRange; x -= interval {
		for y := uint32(valid32 - interval + 1); y >= halfRange; y -= interval {
			if sum, carryOut := Add32(x, y, 0); sum != x+y+0 || carryOut != 1 {
				t.Errorf("Add32(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 0,
					sum, carryOut,
					x+y+0, 1)
			}
			if sum, carryOut := Add32(x, y, 1); sum != x+y+1 || carryOut != 1 {
				t.Errorf("Add32(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 1,
					sum, carryOut,
					x+y+1, 1)
			}
		}
	}
}

func carryOutOne64(t *testing.T) {
	// Sums where carryOut == 1
	// 64-bit
	const (
		halfRange = (1 << 64) / 2
		valid64   = 1<<64 - 1
		interval  = halfRange / 256
	)
	for x := uint64(valid64 - interval + 1); x >= halfRange; x -= interval {
		for y := uint64(valid64 - interval + 1); y >= halfRange; y -= interval {
			if sum, carryOut := Add64(x, y, 0); sum != x+y+0 || carryOut != 1 {
				t.Errorf("Add64(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 0,
					sum, carryOut,
					x+y+0, 0)
			}
			if sum, carryOut := Add64(x, y, 1); sum != x+y+1 || carryOut != 1 {
				t.Errorf("Add64(0x%X, 0x%X, 0x%X) == (0x%X, 0x%X); want (0x%X, 0x%X)",
					x, y, 1,
					sum, carryOut,
					x+y+1, 0)
			}
		}
	}
}
