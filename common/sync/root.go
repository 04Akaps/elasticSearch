package sync

import "sync"

/*
	sync.Pool

	1. 메모리 할당과, 해제에 대한 비용을 줄인다. -> 값을 재사용한다.
	2. GC 부담을 줄인다. -> 주기적으로 반환되면, GC가 회수하는 횟수가 줄어든다.
	3. 동시성 관리를 내부적으로 잘 해준다.
	4. 단기적인 데이터 처리를 위해 설계되었다. -> GC가 언제 회수 할 지 모르기 떄문에
	--> 결론 : 요청이 많은 곳에서 사용이 되어야 한다.
*/

type Pool[T any] struct {
	sync.Pool
	new func() T
}

func (p *Pool[T]) Get() T {
	v := p.Pool.Get()

	if v == nil {
		return p.new()
	}

	return v.(T)
}

func (p *Pool[T]) Put(x T) {
	p.Pool.Put(x)
}

func NewPool[T any](newF func() T) *Pool[T] {
	return &Pool[T]{
		Pool: sync.Pool{
			New: func() interface{} { return newF() },
		},
		new: newF,
	}
}
