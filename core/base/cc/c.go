package cc

// <summary>
// Counter 锁计数器
// <summary>
type Counter interface {
	Up()
	Down()
	Count() int
	Reset()
	Wait()
}
