package db

import (
    "encoding/binary"
)

// Simple helper function to convert a uint64 into an array of bytes
//
func uint64ToBytes(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, v)
    return b
}
