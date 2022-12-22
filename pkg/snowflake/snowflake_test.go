package snowflake

import "testing"

func TestNewNode(t *testing.T) {
	_, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewNode(1024)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}
}

func TestGenerateID(t *testing.T) {

	node, _ := NewNode(1)

	var old, new ID
	for i := 0; i < 1000000; i++ {
		new = node.GenerateID()
		if old == new {
			t.Errorf("old(%d) & new(%d) are the same", old, new)
		}

		if old > new {
			t.Errorf("old value:(%d) is greater then the new value:(%d)", old, new)
		}
		old = new
	}
}
