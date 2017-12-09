package blockchain

type chain struct {
	tip *Block
}

func NewBlockchain(b *Block) *chain {
	return &chain{
		tip: b,
	}
}

func (b *chain) Tip() *Block {
	return b.tip
}

func (b *chain) Mine(d []byte) {
}
