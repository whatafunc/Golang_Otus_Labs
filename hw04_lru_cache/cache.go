package hw04lrucache

import (
	"fmt"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
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
		l.queue.PrintAll()
		fmt.Println("\nLRU items print on Starts: ")
		l.PrintCache()
		fmt.Println("LRU queue = ", l.queue)
	}
	wasInCache := false

	if node, found := l.items[key]; found { // если элемент присутствует в словаре,
		node.Value = value        // то обновить его значение
		l.queue.MoveToFront(node) // и переместить элемент в начало очереди;
		if Debug {
			fmt.Println("after LRU Update = ", l.queue)
		}
		wasInCache = true
	} else {
		newNode := l.queue.PushFront(value) // иначе добавить в словарь
		newNode.Key = key                   // add reference
		l.items[key] = newNode              // и в начало очереди
		if Debug {
			fmt.Println("LRU after Insert = ", l.queue)
		}
		wasInCache = false
	}

	if l.queue.Len() > l.capacity { // если размер очереди больше ёмкости кэша,
		removed := l.queue.Back()
		if removed != nil {
			l.queue.Remove(removed)      // то необходимо удалить последний элемент из очереди
			delete(l.items, removed.Key) // и его значение из словаря);
		}
	}

	if Debug {
		fmt.Println("LRU queue = ", l.queue)
		fmt.Println("LRU items print: ")
		l.queue.PrintAll()
		fmt.Println("\nLRU items print on fyn: ")
		l.PrintCache()
	}
	return wasInCache
}

// Get implements Cache.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	wasInCache := false
	var nodeValue interface{}
	if node, found := l.items[key]; found { // если элемент присутствует в словаре,
		l.queue.MoveToFront(node) // переместить элемент в начало очереди;
		nodeValue = node.Value    // и вернуть его значение и true;
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
