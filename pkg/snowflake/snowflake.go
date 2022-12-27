package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type Node struct {
	mu sync.Mutex

	epoch     int64
	epochTime time.Time

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

func NewNode(node int64, optFuns ...optionsFunc) (*Node, error) {
	opts := defaultOptions
	for _, optFun := range optFuns {
		err := optFun(opts)
		if err != nil {
			return nil, err
		}
	}

	n := &Node{
		node: node,

		epoch:     opts.epoch,
		epochTime: time.Unix(opts.epoch/1000, opts.epoch%1000*1000000),

		timeShift: opts.nodeBits + opts.stepBits,
		nodeShift: opts.stepBits,
		nodeMax:   ^(-1 << opts.nodeBits),
		nodeMask:  ^(-1 << opts.nodeBits) << opts.stepBits,
		stepMask:  ^(-1 << opts.stepBits),
	}

	if node < 0 || node > n.nodeMax {
		return nil, errors.New("node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}
	return n, nil
}

func (n *Node) GenerateID() ID {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epochTime).Milliseconds()
	if now <= n.time { // 防止时钟回拨异常
		n.step = (n.step + 1) & n.stepMask
		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epochTime).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}
	n.time = now

	return ID(now<<n.timeShift | n.node<<n.nodeShift | n.step)
}

func (n *Node) ParseMSTime(id ID) int64 {
	return int64(id)>>n.timeShift + n.epoch
}

func (n *Node) ParseNodeID(id ID) int64 {
	return int64(id) & n.nodeMask >> n.nodeShift
}

func (n *Node) ParseStep(id ID) int64 {
	return int64(id) & n.stepMask
}

func (n *Node) getTime() int64 {
	return n.time + n.epoch
}
