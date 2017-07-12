package conCurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	finish chan bool
)

func TestManager_AddWork(t *testing.T) {
	finish = make(chan bool, 1)
	m := NewManger()

	m.show()
	m.AddWork(println)
	m.show()

	<-finish

	m.show()
}

func println(ctx context.Context) {
	fmt.Println("123")
	time.Sleep(time.Second * 10)
	finish <- true
}
