package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	Epoch int64 = 1671698366302

	TimeBits uint8 = 41
	NodeBits uint8 = 10 // 0 ~ 1023
	StepBits uint8 = 12 // 0 ~ 4095 => 4096/ms

	NodeMax int64 = -1 ^ (-1 << NodeBits)
)

type Node struct {
	mu    sync.Mutex
	epoch time.Time

	time int64
	node int64
	step int64

	timeShift uint8
	nodeShift uint8
	nodeMax   int64
	nodeMask  int64
	stepMask  int64
}

type ID int64

func NewNode(node int64) (*Node, error) {
	if TimeBits+NodeBits+StepBits > 63 {
		return nil, errors.New("you have a total 63 bits to share between Time/Node/Step")
	}
	// todo node 重复
	if node < 0 || node > NodeMax {
		return nil, errors.New("node number must be between 0 and " + strconv.FormatInt(NodeMax, 10))
	}
	return &Node{
		epoch: time.Unix(Epoch/1000, Epoch%1000*1000000),
		node:  node,

		timeShift: NodeBits + StepBits,
		nodeShift: StepBits,
		nodeMax:   ^(-1 << NodeBits),
		nodeMask:  NodeMax << StepBits,
		stepMask:  ^(-1 << StepBits),
	}, nil
}

func (n *Node) GenerateID() ID {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Milliseconds()
	if now <= n.time { // 防止时钟回拨异常
		n.step = (n.step + 1) & n.stepMask
		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}
	n.time = now

	return ID(now<<n.timeShift | n.node<<n.nodeShift | n.step)
}
