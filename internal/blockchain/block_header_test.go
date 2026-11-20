package blockchain

import (
	"bytes"
	"testing"
)

func testBlockHeader() BlockHeader {
	var previousHash Hash
	var dataHash Hash
	var target Target

	for i := 0; i < 32; i++ {
		previousHash[i] = byte(i)
		dataHash[i] = byte(i + 32)
		target[i] = byte(i + 64)
	}

	return BlockHeader{
		Version:      0x01020304,
		PreviousHash: previousHash,
		DataHash:     dataHash,
		Timestamp:    0x0102030405060708,
		Target:       target,
		Nonce:        0x1112131415161718,
	}
}

func TestBlockHeaderSerializeSize(t *testing.T) {
	header := testBlockHeader()

	got := header.serialize()

	if len(got) != headerSize {
		t.Fatalf(
			"serialized header size: got %d, want %d",
			len(got),
			headerSize,
		)
	}
}

func TestBlockHeaderSerialize(t *testing.T) {
	header := testBlockHeader()

	got := header.serialize()

	var expected bytes.Buffer

	expected.Write([]byte{0x01, 0x02, 0x03, 0x04})
	expected.Write(header.PreviousHash[:])
	expected.Write(header.DataHash[:])
	expected.Write([]byte{
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
	})
	expected.Write(header.Target[:])
	expected.Write([]byte{
		0x11, 0x12, 0x13, 0x14,
		0x15, 0x16, 0x17, 0x18,
	})

	if !bytes.Equal(got, expected.Bytes()) {
		t.Fatalf(
			"incorrect serialization: got %x, want %x",
			got,
			expected.Bytes(),
		)
	}
}

func TestBlockHeaderHashCommitsToEveryField(t *testing.T) {
	base := testBlockHeader()
	baseHash := base.Hash()

	tests := []struct {
		name   string
		mutate func(*BlockHeader)
	}{
		{name: "version", mutate: func(h *BlockHeader) { h.Version++ }},
		{name: "previous hash", mutate: func(h *BlockHeader) { h.PreviousHash[0]++ }},
		{name: "data hash", mutate: func(h *BlockHeader) { h.DataHash[0]++ }},
		{name: "timestamp", mutate: func(h *BlockHeader) { h.Timestamp++ }},
		{name: "target", mutate: func(h *BlockHeader) { h.Target[0]++ }},
		{name: "nonce", mutate: func(h *BlockHeader) { h.Nonce++ }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed := base
			tt.mutate(&changed)

			if got := changed.Hash(); got == baseHash {
				t.Fatalf("changing %s did not change header hash %x", tt.name, got)
			}
		})
	}
}
