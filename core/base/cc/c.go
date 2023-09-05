package cc

// 锁计数器
type Counter interface {
	Up()
	Down()
	Count() int
	Reset()
	Wait()
}
