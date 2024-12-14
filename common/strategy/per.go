package strategy

import "go.uber.org/atomic"

/*
	PER 알고리즘을 활용하여 Cache Hit을 유도하고
	Cache Stampede 상황을 방지하기 최대한 방어하기 위한 함수
*/

func PERComputeAndSet(
	needToSet bool,
	atomicKey atomic.Int32,
	fn func(),
) {
	if needToSet && atomicKey.CAS(0, 1) {
		go func() {
			defer atomicKey.Store(0)
			fn()
		}()
	}
}
