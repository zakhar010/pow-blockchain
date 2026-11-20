package blockchain

const (
	genesisData      = "Genesis Block"
	genesisTimestamp = int64(1767225600)
	genesisNonce     = uint64(97120)
)

var genesisHash = Hash{
	0x00, 0x00, 0xed, 0xe6, 0x41, 0x0c, 0x13, 0xe2,
	0x36, 0xc6, 0xd6, 0x68, 0x26, 0x38, 0xd3, 0x21,
	0xf2, 0x96, 0x4a, 0xdf, 0xfd, 0x42, 0x0f, 0x36,
	0x8b, 0x6d, 0xd9, 0x6b, 0x44, 0x2d, 0xf1, 0x11,
}

func genesisBlock() Block {
	return Block{
		Data: genesisData,
		Header: BlockHeader{
			Version:      blockVersion,
			PreviousHash: Hash{},
			DataHash:     hashBytes([]byte(genesisData)),
			Timestamp:    genesisTimestamp,
			Target:       InitialTarget(),
			Nonce:        genesisNonce,
		},
	}
}
