package generator

type Generator interface {
	Next() uint64
	Release(uint64)
}
