package hw04lrucache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct { // helper struct to store correspondence between value <> Key
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.RWMutex // added to work in concurrent mode
}

func NewCache(capacity int) Cache {
	return &lruCache{ // address
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set implements Cache.
func (l *lruCache) Set(key Key, value interface{}) bool {
	if Debug {
		fmt.Println("\nLIST items print on Starts: ")
		fmt.Println("\nLRU items print on Starts: ")
		l.PrintCache()
		fmt.Println("LRU queue = ", l.queue)
	}
	l.mu.Lock()         // Lock for write
	defer l.mu.Unlock() // Unlock when done
	wasInCache := false

	if node, found := l.items[key]; found { // если элемент присутствует в словаре,
		node.Value = cacheItem{key: key, value: value} // то обновить его значение
		l.queue.MoveToFront(node)                      // и переместить элемент в начало очереди;
		if Debug {
			fmt.Println("after LRU Update = ", l.queue)
		}
		wasInCache = true
	} else {
		newNode := l.queue.PushFront(cacheItem{key: key, value: value})
		l.items[key] = newNode // иначе добавить в словарь
		if Debug {
			fmt.Println("LRU after Insert = ", l.queue)
		}
		wasInCache = false
	}

	if l.queue.Len() > l.capacity { // если размер очереди больше ёмкости кэша,
		removed := l.queue.Back()
		if removed != nil {
			l.queue.Remove(removed) // то необходимо удалить последний элемент из очереди
			if item, ok := removed.Value.(cacheItem); ok {
				delete(l.items, item.key) // и его значение из словаря);
			}
		}
	}

	if Debug {
		fmt.Println("LRU queue = ", l.queue)
		fmt.Println("LRU items print: ")
		fmt.Println("\nLRU items print on fyn: ")
		l.PrintCache()
	}
	return wasInCache
}

// Get implements Cache.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.RLock()         // Read lock now
	defer l.mu.RUnlock() // Unlock after this method gets finished
	wasInCache := false
	var nodeValue interface{}
	if node, found := l.items[key]; found { // если элемент присутствует в словаре,
		l.queue.MoveToFront(node) // переместить элемент в начало очереди;
		if item, ok := node.Value.(cacheItem); ok {
			nodeValue = item.value
		}
		if Debug {
			fmt.Println("after LRU found in GET = ", l.queue)
		}
		wasInCache = true
	} else {
		wasInCache = false
	}

	return nodeValue, wasInCache
}

// Clear implements LRUCache.
func (l *lruCache) Clear() {
	l.mu.Lock()                                   // Lock both Read & Write
	defer l.mu.Unlock()                           // Unlock
	l.items = make(map[Key]*ListItem, l.capacity) // new clean Mapa
	l.queue = NewList()                           // new empty list
}

func (l *lruCache) PrintCache() {
	fmt.Printf("LRU Cache (Capacity: %d, Size: %d)\n", l.capacity, l.queue.Len())
	// Print the map contents
	fmt.Println("Map Contents:")
	if len(l.items) == 0 {
		fmt.Println("(empty)")
	} else {
		for key, node := range l.items {
			fmt.Printf("Key: %v => Value: %v\n", key, node.Value)
		}
	}
	// Print the queue order
	// c.PrintQueue()
}
