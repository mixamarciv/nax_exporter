package app_fnc

import (
	"crypto/md5"
	"fmt"
	crc32 "hash/crc32"
	"strconv"
)

func StrMd5(text []byte) string {
	d := md5.Sum(text)
	s := fmt.Sprintf("%x", d)
	return s
}

func StrCrc32(text []byte) string {
	h := crc32.NewIEEE()
	h.Write(text)
	v := h.Sum32()
	s := strconv.FormatUint(uint64(v), 32)
	//s := fmt.Sprintf("%d", v)
	return s
}
