package snowflake

import (
	"errors"
	"time"
)

type options struct {
	epoch    int64
	timeBits uint8
	nodeBits uint8
	stepBits uint8
}

type optionsFunc func(*options) error

// Milliseconds
func SetEpoch(epoch int64) optionsFunc {
	return func(o *options) error {
		if epoch > time.Now().UnixMilli() {
			return errors.New("cannot be in the future")
		}
		o.epoch = epoch
		return nil
	}
}

func CustomBit(timeBits, nodeBits, stepBits uint8) optionsFunc {

	return func(o *options) error {
		if timeBits+nodeBits+stepBits != 63 || timeBits == 0 || nodeBits == 0 || stepBits == 0 {
			return errors.New("you have a total 63 bits to share between Time/Node/Step")
		}
		o.timeBits = timeBits
		o.nodeBits = nodeBits
		o.stepBits = stepBits
		return nil
	}
}

var defaultOptions = &options{
	epoch: 1671698366302, // 2022-12-22 16:39:26

	timeBits: 41,
	nodeBits: 10, // 0 ~ 1023
	stepBits: 12, // 0 ~ 4095 => 4096/ms
}
