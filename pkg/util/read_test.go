package util

import (
	"bufio"
	"bytes"
	"strings"
	"sync"
	"testing"
)

func BenchmarkReadFunc(t *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	for i := 0; i < t.N; i++ {

		arr := []byte(strings.Repeat("a", 1024*1024))
		reader := bufio.NewReader(bytes.NewBuffer(arr))

		res := pool.Get().([]byte)
		ans := bytes.NewBuffer([]byte{})
		w := bufio.NewWriter(ans)
		for {
			num, err := reader.Read(res)
			if err != nil {
				break
			}
			w.Write(res[:num])
			w.Flush()
		}
		pool.Put(res)
	}
}
func BenchmarkReadFunc2(t *testing.B) {
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 10000)
		},
	}
	for i := 0; i < t.N; i++ {

		arr := []byte(strings.Repeat("a", 1024*1024))
		reader := bufio.NewReader(bytes.NewBuffer(arr))
		res := pool.Get().([]byte)
		ans := bytes.NewBuffer([]byte{})
		w := bufio.NewWriter(ans)
		for {
			num, err := reader.Read(res)
			if err != nil {
				break
			}
			w.Write(res[:num])
			w.Flush()
		}
		pool.Put(res)
	}
}
