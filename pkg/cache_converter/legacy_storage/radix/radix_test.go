package radix

import (
	"strconv"
	"testing"
	"time"

	"github.com/minotar/imgd/storage"
	"github.com/minotar/imgd/storage/util/helper"
)

func TestInsertAndRetrieve(t *testing.T) {
	cache, _ := New(RedisConfig{
		Network: "tcp",
		Address: "127.0.0.1:6379",
		Size:    1,
	})
	preSize := cache.Size()
	for i := 0; i < 10; i++ {
		str := helper.RandString(32)
		cache.Insert(str, []byte(strconv.Itoa(i)), time.Minute)
		item, _ := cache.Retrieve(str)
		if string(item) != strconv.Itoa(i) {
			t.Fail()
		}
	}
	if cache.Len() != 10 {
		t.Fail()
	}
	if cache.Size() < preSize {
		t.Fail()
	}
}

func TestRetrieveMiss(t *testing.T) {
	cache, _ := New(RedisConfig{
		Network: "tcp",
		Address: "127.0.0.1:6379",
		Size:    1,
	})
	str := helper.RandString(32)
	item, err := cache.Retrieve(str)
	if err.Error() != storage.ErrNotFound.Error() {
		t.Log("ErrNotFound was instead:", err.Error())
		t.Fail()
	}
	if string(item) != "" {
		t.Fail()
	}

}

type testBucket struct {
	keys  []string
	cache *RedisCache
	size  int
}

var largeBucket testBucket

func initLargeBucket(n int) {
	if largeBucket.size > n {
		return
	}
	keys := make([]string, n)
	cache, _ := New(RedisConfig{
		Network: "tcp",
		Address: "127.0.0.1:6379",
		Size:    1,
	})
	for i := 0; i < n; i++ {
		keys[i] = helper.RandString(32)
		cache.Insert(keys[i], []byte(strconv.Itoa(i)), time.Minute)
	}

	largeBucket = testBucket{keys, cache, n}
}

func BenchmarkInsert(b *testing.B) {
	initLargeBucket(b.N)
	cache, _ := New(RedisConfig{
		Network: "tcp",
		Address: "127.0.0.1:6379",
		Size:    1,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Insert(largeBucket.keys[i], []byte(strconv.Itoa(i)), time.Minute)
	}
}

func BenchmarkLookup(b *testing.B) {
	initLargeBucket(b.N)

	iter := 10
	if b.N < 10 {
		iter = b.N
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < iter; k++ {
			largeBucket.cache.Retrieve(largeBucket.keys[k])
		}
	}
}
