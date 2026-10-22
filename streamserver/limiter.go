package main

import "fmt"

type ConnLimiter struct {
	currentConn int
	bucket      chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	conn := &ConnLimiter{
		currentConn: cc,
		bucket:      make(chan int, cc),
	}
	return conn
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.currentConn {
		fmt.Println("reached the rate limitation")
		return false
	}
	cl.bucket <- 1
	return true
}
func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	fmt.Println("New connection comming ", c)
}
