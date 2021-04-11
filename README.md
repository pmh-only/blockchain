# pmh-only/blockchain
비효율적일지도 모를 블럭체인

## 블럭 구조
각 블럭은 `Head`, `Body`, `Tail`로 구성됩니다.

### Head: BlockHead
블럭의 대한 정보가 포함되는 곳입니다.\
총 `75 Bytes`로 구성되어 있습니다.

```go
type BlockHead struct {
	Index      uint16 // 2 Bytes
	CreatedAt  uint32 // 4 Bytes
	PrevHash   []byte // 64 Bytes limit
	Nonce      uint32 // 4 Bytes
	Difficulty uint8 // 1 Byte
}
```

### Body: BlockBody
블럭의 데이터가 저장됩니다.\
길이에는 제한이 없습니다.

```go
type BlockBody []byte // no limit
```

### Tail: BlockTail
블럭의 Head와 Body를 해시한 값이 저장됩니다.

```go
type BlockTail struct {
	CurrHash []byte // 64 Bytes limit
}
```

### 예시
Block Structure:
```go
Block{
  Head: {
    Index: 0,
    CreatedAt: 1618114690
    PrevHash: 0
    Nonce: 3, // there is no pow yet
    Difficulty: 4
  },
  Body: []byte("hi"),
  Tail: {
    CurrHash: SHAKE256(Head + Body, 512)
  }
}
```

Serialization:\
`
00000000 00000000 01100000 01110010 01111001 11111111 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000011 00000100 01101000 01101001 00100110 10000000 01001010 10000101 10001101 00101101 11100010 00110101 00011010 11001001 00001010 00110100 10111000 01111110 11000100 01110110 10000110 11000111 10110111 01101101 01010111 00000111 11100010 11000000 00010110 00101000 01101100 11111100 11010110 00100111 00011100 11110000 10100000 01111101 10100100 00011101 00111101 10000011 11100011 10011111 10001000 01100110 10110001 10010011 11101011 11010001 01011111 10111000 10110111 11100111 00011100 01111101 10000111 11001100 10111101 10111101 10100111 11001101 10101111 01100010 01100100 10010000 01001100 01011100`

처음 2바이트는 일련번호\
다음 4바이트는 생성시각\
다음 64바이트는 이전 블럭 해시\
다음 4바이트는 Nonce값\
다음 1바이트는 채굴난이도\
다음부터 마지막 64바이트 전까지는 데이터\
마지막 64바이트는 현재 블럭 해시
