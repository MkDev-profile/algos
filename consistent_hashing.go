// in-memory distributed/sharded(using consistent hashing) cache.

package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
	//"strconv"
)

// CNode represents a physical node in the cluster
type CNode struct {
	Name string
}

// ConsistentHash implements consistent hashing
type ConsistentHash struct {
	virtualNodeToNodeMap     map[uint32]CNode // mapping: virtualNodeHash -> node
	sortedVirtualNodesHashes []uint32        // hashes ring = sorted hashes for binary search
	virtualNodesNum int           // number of virtual nodes per physical node
}

// NewConsistentHash creates a new ConsistentHash instance
func NewConsistentHash(virtualNodes int) *ConsistentHash {
	return &ConsistentHash{
		virtualNodeToNodeMap:       make(map[uint32]CNode),
		sortedVirtualNodesHashes:  make([]uint32, 0),
		virtualNodesNum: virtualNodes,
	}
}

// hashKey generates a hash for a given key
func (ch *ConsistentHash) hashKey(key string) uint32 {
	h := sha1.New()
	h.Write([]byte(key))
	hash := h.Sum(nil)
	
	// Take first 4 bytes for uint32 hash
	return uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
}

// AddNode adds a physical node to the hash ring
func (ch *ConsistentHash) AddNode(node CNode, nodeIdForExample int) {
	for i := 0; i < ch.virtualNodesNum; i++ {
		//virtualNodeKey := node.Name + "#" + strconv.Itoa(i)
		//virtualNodeHash := ch.hashKey(virtualNodeKey)
		virtualNodeHash := uint32(i*3 + nodeIdForExample)
		
		if _, exists := ch.virtualNodeToNodeMap[virtualNodeHash]; !exists {
			ch.virtualNodeToNodeMap[virtualNodeHash] = node
			ch.sortedVirtualNodesHashes = append(ch.sortedVirtualNodesHashes, virtualNodeHash)
		}
	}
	
	// Sort the keys for binary search
	sort.Slice(ch.sortedVirtualNodesHashes, func(i, j int) bool {
		return ch.sortedVirtualNodesHashes[i] < ch.sortedVirtualNodesHashes[j]
	})

	fmt.Printf("Adding Node: %s\n", node.Name)
	fmt.Printf("  ch.virtualNodes: %v\n", ch.virtualNodeToNodeMap)
	fmt.Printf("  ch.sortedVirtualNodesHashes: %v\n", ch.sortedVirtualNodesHashes)
}

// RemoveNode removes a physical node from the hash ring
func (ch *ConsistentHash) RemoveNode(nodeName string) {
	keysToRemove := make([]uint32, 0)
	
	// Find all virtual nodes for this physical node
	for virtualNodeHash, node := range ch.virtualNodeToNodeMap {
		if node.Name == nodeName {
			keysToRemove = append(keysToRemove, virtualNodeHash)
		}
	}
	
	// Remove from nodes map
	for _, hash := range keysToRemove {
		delete(ch.virtualNodeToNodeMap, hash)
	}
	
	// Rebuild sorted keys
	ch.sortedVirtualNodesHashes = make([]uint32, 0, len(ch.virtualNodeToNodeMap))
	for hash := range ch.virtualNodeToNodeMap {
		ch.sortedVirtualNodesHashes = append(ch.sortedVirtualNodesHashes, hash)
	}
	sort.Slice(ch.sortedVirtualNodesHashes, func(i, j int) bool {
		return ch.sortedVirtualNodesHashes[i] < ch.sortedVirtualNodesHashes[j]
	})
}

// GetNode returns the node responsible for the given key
func (ch *ConsistentHash) GetNode(item string) (CNode, bool) {
	if len(ch.virtualNodeToNodeMap) == 0 {
		return CNode{}, false
	}
	
	//hash := ch.hashKey(key)
	intHash,_ := strconv.Atoi(item) // for testing
	itemHash := uint32(intHash)
	
	// Binary search for the first node with hash >= key's hash
	// p.s. search for i from 0 to n, n=len(ch.sortedVirtualNodesHashes)
	// p.s. if i for condition not found, returns n
	virtualNodeIndex := sort.Search(len(ch.sortedVirtualNodesHashes), func(i int) bool {
		return ch.sortedVirtualNodesHashes[i] >= itemHash
	})
	
	// If we reached the end (i=n), wrap around to the first node
	if virtualNodeIndex == len(ch.sortedVirtualNodesHashes) {
		virtualNodeIndex = 0
	}
	
	node, exists := ch.virtualNodeToNodeMap[ch.sortedVirtualNodesHashes[virtualNodeIndex]]
	return node, exists
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type Cache struct {
	items map[string]Item
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]Item),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.Expired() {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Item)
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
		}
	}
}

func Main_ConsistentHash() {
	ch := NewConsistentHash(3)
	
	// Add some nodes
	nodes := []CNode{
		{Name: "node0"},
		{Name: "node1"},
		{Name: "node2"},
		{Name: "node3"},
	}
	
	for i, node := range nodes {
		ch.AddNode(node, i*10)
	}
	
	// Test key lookups
	testKeys := []string{
		/*
		"user_12345",
		"product_67890",
		"order_11111",
		"session_22222",
		"cache_33333",
		*/
		"5",
		"8",
		"15",
		"18",
		"37",
		"100",
	}
	
	fmt.Println("\nItem to Node mapping:")
	for _, key := range testKeys {
		node, _ := ch.GetNode(key)
		fmt.Printf("Key: %s -> Node: %s\n", key, node.Name)
	}
	
	// Test node removal
	fmt.Println("\nRemoving node1...")
	ch.RemoveNode("node1")
	
	fmt.Println("\nItem to Node mapping after removal:")
	for _, key := range testKeys {
		node, _ := ch.GetNode(key)
		fmt.Printf("Key: %s -> Node: %s\n", key, node.Name)
	}
}

















