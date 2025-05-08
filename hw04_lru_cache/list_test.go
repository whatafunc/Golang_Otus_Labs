package hw04lrucache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
		if Debug {
			fmt.Println("\033[H\033[2J")
			l.PushFront("a")
			l.PushBack("bb")
			l.PushBack("ccc")
			fmt.Println(l.Front().Value) // "a"
			fmt.Println(l.Back().Value)  // "ccc"
			fmt.Println("========================1.done")
		}
	})

	t.Run("complex", func(t *testing.T) {
		if Debug {
			fmt.Println("========================2.start")
		}
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		if Debug {
			fmt.Println(l)
		}
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		if Debug {
			fmt.Println("after removal Len =", l.Len())
		}
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]
		if Debug {
			l.PrintAll() // OK [80, 60, 40, 10, 30, 50, 70]
		}

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)
		if Debug {
			l.PrintAll() //
		}
		l.MoveToFront(l.Front()) // send head elem to start  // [80, 60, 40, 10, 30, 50, 70]
		if Debug {
			l.PrintAll() //
		}
		l.MoveToFront(l.Back()) // send tail elem to start  // [70, 80, 60, 40, 10, 30, 50]
		if Debug {
			l.PrintAll() //
		}

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
		if Debug {
			l.PrintAll() //
		}
	})
}
