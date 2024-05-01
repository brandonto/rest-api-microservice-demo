package db

import (
    "encoding/binary"
)

func uint64ToBytes(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, v)
    return b
}
