package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		if Debug {
			fmt.Println("\033[H\033[2J")
			fmt.Println("quee = ", c)

			wasInCache := c.Set("aaa", 100)
			require.False(t, wasInCache)
			wasInCache = c.Set("aaa", 101)
			require.True(t, wasInCache)
			wasInCache = c.Set("bbb", 777)
			require.False(t, wasInCache)
			wasInCache = c.Set("bbb", 778)
			require.True(t, wasInCache)

			fmt.Println("\n wasInCache = ", wasInCache)
			fmt.Println("quee after = ", c)
		}

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		c.Clear() // delete all items
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	// t.Skip() // Remove me if task with asterisk completed.
	// test results shared here:
	// https://docs.google.com/document/d/1xdr1-q-F1b3uvxZ5FBxlinSdBMvEzDPrnuj-ycRsMWU/edit?usp=sharing

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			require.NotPanics(t, func() { // catch no panic or race condition
				c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
			})
		}
	}()

	wg.Wait()
}
