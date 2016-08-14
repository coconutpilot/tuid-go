package tuid

import (
	"fmt"
	"time"
)

type tuidCtx struct {
	nsec_offset        uint64 /* nsec offset from unixtime epoch */
	nsec               uint64 /* unixtime in nsec, adjusted by offset */
	nsec_mask          uint64 /* the bits of nsec to be used in the TUID */
	nsec_shift         uint8  /* how far to left shift the nsec bits */
	nsec_min_increment uint64 /* minimum increment required if nsec mask was too coarse */
	id                 uint64 /* static id, ie machine id */
	counter            uint64 /* a counter, current value */
	counter_max        uint64 /* when to reset the counter to 0 */
	counter_shift      uint8  /* how far to left shift the counter bits */
	random             uint64 /* a random #, also used as seed */
	random_mask        uint64 /* the bits of the random # to be used */
	random_shift       uint8  /* XXX: turns out this is not needed */
}

func New(spec string) (*tuidCtx, error) {
	// XXX: temp logging
	fmt.Printf("Decoding TUID spec: %s\n", spec)

	tuidGen := &tuidCtx{random: 0x5248c8561600f46d}

	bitpos := uint8(64)

	max := len(spec)
	i := 0

	// this string parser is a direct translation from C which means
	// it isn't idiomatic Go.  Undecided if I should rewrite it or
	// keep the libraries as similar as possible ...
	for i < max {
		identifier := spec[i]
		i++

		var value uint64
		for i < max {
			v := spec[i]
			if v < '0' || v > '9' {
				break
			}
			i++
			value = value*10 + uint64(v-'0')
		}
		// XXX: temp logging
		fmt.Printf("Key: %v Val: %v\n", string(identifier), value)

		switch identifier {
		case 'E':
			tuidGen.nsec_offset = value

		case 'N':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			tuidGen.nsec_mask = (^uint64(0)) >> (64 - value)
			bitpos -= uint8(value)
			tuidGen.nsec_shift = bitpos

		case 'C':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			tuidGen.counter_max = (^uint64(0)) >> (64 - value)
			tuidGen.counter = tuidGen.counter_max /* this forces initialization */
			bitpos -= uint8(value)
			tuidGen.counter_shift = uint8(bitpos)

		case 'R':
			if value > uint64(bitpos) {
				return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
			}
			tuidGen.random_mask = (^uint64(0)) >> (64 - value)
			bitpos -= uint8(value)
			tuidGen.random_shift = bitpos

		case 'I':
			tuidGen.id = value

		default:
			return nil, fmt.Errorf("TUID spec error at: %v%v", identifier, value)
		}

	}

    // the min increment is the LSB of the nanosecond component
    tuidGen.nsec_min_increment = ((tuidGen.nsec_mask ^ (tuidGen.nsec_mask - 1)) >> 1) + 1

	return tuidGen, nil
}

func (ctx *tuidCtx) Gen() uint64 {
	ctx.counter++
	if ctx.counter > ctx.counter_max {
		ctx.counter = 0

		nsec := uint64(time.Now().UnixNano())

		nsec <<= ctx.nsec_shift
		if ctx.nsec >= nsec {
			fmt.Printf("Collision, incrementing by %v\n", ctx.nsec_min_increment)
			ctx.nsec += ctx.nsec_min_increment
		} else {
			ctx.nsec = nsec
		}
	}

	tuid := ctx.nsec

	counter := ctx.counter
	counter <<= ctx.counter_shift
	tuid |= counter

	tuid |= ctx.id

	//    uint64 rnd = xorshift64(&(ctx->random));
	rnd := ctx.random

	rnd &= ctx.random_mask
	rnd <<= ctx.random_shift
	tuid |= rnd

	return tuid
}
