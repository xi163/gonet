package keepalive

import (
	"time"

	"github.com/cwloo/gonet/core/base/cc"
	"github.com/cwloo/gonet/core/net/keepalive/bucket"
)

var (
	flag    = cc.NewAtomFlag()
	buckets bucket.Buckets
)

func Init(timeout, d time.Duration) {
	if flag.TestSet() {
		second := int32(int64(timeout) / int64(time.Second))
		buckets = bucket.NewBuckets(second, d)
	}
}
