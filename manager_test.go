package conCurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var (
	finish chan bool
	count  int
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
	time.Sleep(time.Second * 5)
	finish <- true
}

func TestManager_PanicHandler(t *testing.T) {
	m := NewManger()

	m.show()
	m.AddWork(panicPrintln)
	m.show()

	m.show()
}

func panicPrintln(ctx context.Context) {
	fmt.Println("123")
	time.Sleep(time.Second * 2)
	panic("123")

}

func TestManager_AddWork2(t *testing.T) {
	m := NewManger()
	m.upperLimit = 5

	for i := 0; i < 10; i++ {
		m.AddWork(func(ctx context.Context) {
			count++
			fmt.Println(count)
			time.Sleep(time.Second * 3)
		})
	}
	//go func() {
	//	for i := 0; i < 12; i++ {
	//		m.show()
	//		time.Sleep(time.Second)
	//	}
	//}()
	time.Sleep(time.Second * 15)
}

func simple(ctx context.Context) {
	fmt.Println("123")
	time.Sleep(time.Second * 3)
}
