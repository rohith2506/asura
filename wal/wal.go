/*
Implementation of write ahead log. Design inspired by rocksdb

1. Each logfile consists of 32kB blocks. 
2. Each block can be one or many records. But if it's less than 6 bytes of space left for next record, we fill that with padding. Why? Because of next step
3. Each record is made up of
    4 bytes => hash of payload using CRC
    2 bytes => size of payload
    1 byte  => type of payload
    remaining number of bytes => payload

type of payload => full, first, middle, last

full => all payload fits into the log
first, middle and last => payload which don't fit into the block and broken into other blocks

@Author: Rohith Uppala
*/

package wal

import (
    "io"
    "os"
)

// All constants
const (
    BLOCKSIZE = 32 * 1024
    MIN_BLOCK_SPACE_NEEDED = 6
    CRC_HASH_SIZE = 4
    PAYLOAD_SIZE = 2
    PAYLOAD_TYPE_SIZE = 1
)

type Record struct {
    crc_hash [CRC_HASH_SIZE]byte{}
    payload_size [PAYLOAD_SIZE]byte{}
    payload_type [PAYLOAD_TYPE_SIZE]byte{}
    payload []byte{}
}

type RecordWriter struct {
    record *Record
    filePointer *FilePointer
}

func (r *Record) Payload() []byte{} {
    return r.payload
}

func (rWriter *RecordWriter) Write() (error) {

}
