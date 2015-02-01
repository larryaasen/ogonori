package varint

//
// Follows the encoding documented at:
// https://developers.google.com/protocol-buffers/docs/encoding?csw=1
//

/* ---[ encode ]--- */

func zigzagEncodeUInt32(n int32) uint32 {
	return uint32((n >> 31) ^ (n << 1))
}

func zigzagEncodeUInt64(n int64) uint64 {
	return uint64((n >> 63) ^ (n << 1))
}

/* ---[ decode ]--- */

func zigzagDecodeInt32(n uint32) int32 {
	return int32((-(n & 1)) ^ (n >> 1))
}

func zigzagDecodeInt64(n uint64) int64 {
	return int64((-(n & 1)) ^ (n >> 1))
}
