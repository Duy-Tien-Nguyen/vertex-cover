package vc

import (
	"fmt"
	"math/bits"
	"strings"
)

type BitMask struct {
	n    int
	data []uint64
}

func NewBitMask(n int) *BitMask {
	numBlocks := (n + 63) / 64
	return &BitMask{
		n:    n,
		data: make([]uint64, numBlocks),
	}
}

func (bm *BitMask) Set(i int) {
	bm.data[i/64] |= 1 << (i % 64)
}

func (bm *BitMask) Unset(i int) {
	bm.data[i/64] &^= 1 << (i % 64)
}

func (bm *BitMask) Get(i int) bool {
	return bm.data[i/64]&(1<<(i%64)) != 0
}

func (bm *BitMask) Count() int {
	count := 0
	for _, block := range bm.data {
		count += bits.OnesCount64(block)
	}
	return count
}

func (bm *BitMask) Clone() *BitMask {
	newData := make([]uint64, len(bm.data))
	copy(newData, bm.data)
	return &BitMask{n: bm.n, data: newData}
}

func (bm *BitMask) And(other *BitMask) {
	for i := range bm.data {
		bm.data[i] &= other.data[i]
	}
}

func (bm *BitMask) Or(other *BitMask) {
	for i := range bm.data {
		bm.data[i] |= other.data[i]
	}
}

func (bm *BitMask) Xor(other *BitMask) {
	for i := range bm.data {
		bm.data[i] ^= other.data[i]
	}
}

func (bm *BitMask) Equals(other *BitMask) bool {
	if bm.n != other.n {
		return false
	}
	for i := range bm.data {
		if bm.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

func (bm *BitMask) ShiftLeft(k int) {
	if k <= 0 {
		return
	}
	totalBits := bm.n
	newMask := NewBitMask(totalBits)
	for i := totalBits - 1; i >= k; i-- {
		if bm.Get(i - k) {
			newMask.Set(i)
		}
	}
	copy(bm.data, newMask.data)
}

func (bm *BitMask) ShiftRight(k int) {
	if k <= 0 {
		return
	}
	totalBits := bm.n
	newMask := NewBitMask(totalBits)
	for i := 0; i < totalBits-k; i++ {
		if bm.Get(i + k) {
			newMask.Set(i)
		}
	}
	copy(bm.data, newMask.data)
}

// Toggle flips bit i
func (bm *BitMask) Toggle(i int) {
	bm.data[i/64] ^= 1 << uint(i%64)
}

func (bm *BitMask) Print() {
	var sb strings.Builder
	for i := 0; i < bm.n; i++ {
		if bm.Get(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	fmt.Println(sb.String())
}
