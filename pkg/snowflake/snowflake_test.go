package snowflake

import (
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	_, err := NewNode(0)
	if err != nil {
		t.Errorf("error creating NewNode, %s", err)
	}

	_, err = NewNode(1024)
	if err == nil {
		t.Errorf("no error creating NewNode")
	}
}

func TestGenerateID(t *testing.T) {

	node, _ := NewNode(0)

	var old, new ID
	for i := 0; i < 1000000; i++ {
		new = node.GenerateID()
		if old == new {
			t.Errorf("old:(%d) & new:(%d) are the same", old, new)
		}

		if old > new {
			t.Errorf("old value:(%d) is greater then the new value:(%d)", old, new)
		}
		old = new
	}
}

func TestParse(t *testing.T) {
	node, _ := NewNode(0)
	id := node.GenerateID()
	nodeID := node.ParseNodeID(id)
	if nodeID != 0 {
		t.Errorf("ParseNodeID error nodeID:(%d) want:(0)", nodeID)
	}
	step := node.ParseStep(id)
	if step != 0 {
		t.Errorf("ParseStep error step:(%d) want:(0)", step)
	}
	time := node.ParseMSTime(id)
	nodeTime := node.getTime()
	if time != nodeTime {
		t.Errorf("ParseMSTime error time:(%d) want:(%d)", time, nodeTime)
	}
}

func TestCustomBit(t *testing.T) {
	// node timeBits nodeBits stepBits result(err == nil => 0)
	useCase := [][]int{
		{0, 21, 21, 21, 0},
		{255, 39, 8, 16, 0},

		// error case
		{0, 1, 2, 3, 1},
		{256, 39, 8, 3, 1},
		{0, 0, 32, 31, 1},
		{0, 31, 0, 32, 1},
		{0, 31, 32, 0, 1},
	}

	for _, uc := range useCase {
		_, err := NewNode(int64(uc[0]), CustomBit(uint8(uc[1]), uint8(uc[2]), uint8(uc[3])))
		if (err != nil && uc[4] == 0) || (err == nil && uc[4] == 1) {
			t.Errorf("error NewNode CustomBit, case %v err %s", uc, err)
		}
	}
}

func TestSetEpoch(t *testing.T) {
	_, err := NewNode(0, SetEpoch(time.Now().UnixMilli()+1000))
	if err == nil {
		t.Errorf("no error NewNode SetEpoch")
	}
	_, err = NewNode(0, SetEpoch(time.Now().UnixMilli()))
	if err != nil {
		t.Errorf("error NewNode SetEpoch, %s", err)
	}
}

func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(0)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.GenerateID()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {

	node, _ := NewNode(0, CustomBit(32, 1, 30))

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.GenerateID()
	}
}
