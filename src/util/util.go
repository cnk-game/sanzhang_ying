package util

import (
	"os"
	"strings"
	"fmt"
	"reflect"
	"unsafe"
	"time"
)

func Hostname() string {
	h, err := os.Hostname()
	if err == nil {
		return ShortHostname(h)
	}
	return ""
}

func ShortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}


func ByteSlice(slice interface{}) (data []byte) {
    sv := reflect.ValueOf(slice)
    if sv.Kind() != reflect.Slice {
        panic(fmt.Sprintf("ByteSlice called with non-slice value of type %T", slice))
    }
    h := (*reflect.SliceHeader)((unsafe.Pointer(&data)))
    h.Cap = sv.Cap() * int(sv.Type().Elem().Size())
    h.Len = sv.Len() * int(sv.Type().Elem().Size())
    h.Data = sv.Pointer()
    return
}


func GetDayZero() int64 {
    now := time.Now()
    zero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

    return zero.Unix()
}