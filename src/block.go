package blockchain

import (
	"encoding/binary"
	"reflect"
	"strings"
	"time"

	"github.com/pmh-only/blockchain/utils"
	"golang.org/x/crypto/sha3"
)

func CreateBlock(index uint16, prevHash []byte, flag BodyFlags, message []byte) Block {
	block := Block{}

	block.Head.Index = index
	block.Head.CreatedAt = uint32(time.Now().Unix())
	block.Head.PrevHash = prevHash
	block.Head.Nonce = 0
	block.Head.Difficulty = START_DIFFICULTY + uint8(index/DIFFICULTY_INCREASE_STEP)

	block.Body.Flag = flag
	block.Body.Message = message

	block.Tail.CurrHash = block.CaculateHash()

	return block
}

func (block Block) CaculateHash() []byte {
	serial := block.SerializationWithoutTail()
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, serial)
	return hash
}

func (block Block) SerializationWithoutTail() []byte {
	head := make([]byte, 0, SIZE_OF_HEAD)
	body := append([]byte{byte(block.Body.Flag)}, block.Body.Message...)

	index := make([]byte, 2)
	createdAt := make([]byte, 4)
	nonce := make([]byte, 4)

	prevHash := utils.PadOrTrim(block.Head.PrevHash, 64)
	difficulty := block.Head.Difficulty

	binary.BigEndian.PutUint16(index, block.Head.Index)
	binary.BigEndian.PutUint32(createdAt, block.Head.CreatedAt)
	binary.BigEndian.PutUint32(nonce, block.Head.Nonce)

	head = append(head, index...)
	head = append(head, createdAt...)
	head = append(head, prevHash...)
	head = append(head, nonce...)
	head = append(head, difficulty)

	return append(head, body...)
}

func (block Block) SerializationWithTail() []byte {
	return append(block.SerializationWithoutTail(), block.CaculateHash()...)
}

func (block *Block) MineBlock() {
	for !block.IsMined() {
		block.Head.Nonce++
		block.Tail.CurrHash = block.CaculateHash()
	}
}

func (block Block) IsValid() bool {
	if !block.IsMined() {
		return false
	}

	if START_DIFFICULTY+uint8(block.Head.Index/DIFFICULTY_INCREASE_STEP) != block.Head.Difficulty {
		return false
	}

	serial := block.SerializationWithoutTail()
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, serial)

	return reflect.DeepEqual(block.Tail.CurrHash, hash)
}

func (block Block) IsMined() bool {
	hash := utils.BytesToBinString(block.Tail.CurrHash)
	diff := strings.Repeat("0", int(block.Head.Difficulty))

	return strings.HasPrefix(hash, diff)
}

func Deserialization(serial []byte) Block {
	block := Block{}

	block.Head.Index = binary.BigEndian.Uint16(serial[:2])
	block.Head.CreatedAt = binary.BigEndian.Uint32(serial[2:6])
	block.Head.PrevHash = serial[6:70]
	block.Head.Nonce = binary.BigEndian.Uint32(serial[70:74])
	block.Head.Difficulty = serial[74]

	block.Body.Flag = BodyFlags(serial[75])
	block.Body.Message = serial[76 : len(serial)-64]

	block.Tail.CurrHash = serial[len(serial)-64:]

	return block
}
