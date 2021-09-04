// PoW algorithm
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int // bugint/bignum(integer)
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)                  // 00000001
	target.Lsh(target, uint(256-Difficulty)) // uint: unsigned integer / LSH: left shift

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 { // nonce < max of Int64
		data := pow.InitData(nonce)
		hash := sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) // hash [32]byte != pow.Target: big.Int(integer)

		if intHash.Cmp(pow.Target) == -1 { // intHash < pow.Target -> -1 / intHash > pow.Target -> 1
			break
		} else {
			nonce++
		}
	}
	fmt.Println() // br

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:]) // hash [32]byte != pow.Target: big.Int

	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num) // to byte from upper side / ex:1234ABCD(hexadecimal) -> [12 34 AB CD]
	Handle(err)

	return buff.Bytes()
}
