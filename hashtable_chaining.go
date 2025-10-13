package main

import (
	"bytes"
	"fmt"
)

// HashTableNode represents a key-value pair in the hash table chain
type HashTableNode struct {
	key   string      // The key used for hashing and lookup
	value interface{} // The value associated with the key
	next  *HashTableNode       // Pointer to the next node in the chain (for collision resolution)
}

// HashTable represents the hash table data structure using separate chaining
type HashTable struct {
	buckets []*HashTableNode  // Array of pointers to the head of each chain
	size    int      // Current number of key-value pairs in the table
	capacity int     // Total number of buckets available
}

// NewHashTable creates and returns a new hash table with the specified capacity
func NewHashTable(capacity int) *HashTable {
	return &HashTable{
		buckets: make([]*HashTableNode, capacity),
		size:    0,
		capacity: capacity,
	}
}

// hashFunction computes the bucket index for a given key
// It uses a simple polynomial rolling hash for good distribution
func (ht *HashTable) hashFunction(key string) int {
	// Prime number base for polynomial rolling hash
	const prime = 31
	
	hash := 0
	for _, char := range key {
		// Multiply current hash by prime and add character code
		hash = (hash*prime + int(char)) % ht.capacity
	}
	
	// Ensure non-negative index
	return (hash + ht.capacity) % ht.capacity
}

// Put inserts a key-value pair into the hash table
// If the key already exists, it updates the value
func (ht *HashTable) Put(key string, value interface{}) {
	// Calculate the bucket index for the key
	index := ht.hashFunction(key)
	
	// Check if key already exists in the chain
	current := ht.buckets[index]
	for current != nil {
		if current.key == key {
			// Key found - update the value
			current.value = value
			return
		}
		current = current.next
	}
	
	// Key not found - create new node and add to front of chain
	newNode := &HashTableNode{
		key:   key,
		value: value,
		next:  ht.buckets[index], // Point to current head of chain
	}
	
	ht.buckets[index] = newNode // Make new node the head of chain
	ht.size++
}

// Get retrieves the value associated with a key
// Returns the value and true if key exists, nil and false otherwise
func (ht *HashTable) Get(key string) (interface{}, bool) {
	index := ht.hashFunction(key)
	
	// Traverse the chain to find the key
	current := ht.buckets[index]
	for current != nil {
		if current.key == key {
			return current.value, true
		}
		current = current.next
	}
	
	// Key not found
	return nil, false
}

// Delete removes a key-value pair from the hash table
// Returns true if key was found and deleted, false otherwise
func (ht *HashTable) Delete(key string) bool {
	index := ht.hashFunction(key)
	
	// Handle case where node to delete is head of chain
	current := ht.buckets[index]
	if current != nil && current.key == key {
		ht.buckets[index] = current.next
		ht.size--
		return true
	}
	
	// Traverse chain to find node to delete
	var prev *HashTableNode
	for current != nil {
		if current.key == key {
			// Found the node - remove it from chain
			prev.next = current.next
			ht.size--
			return true
		}
		prev = current
		current = current.next
	}
	
	// Key not found
	return false
}

// Contains checks if a key exists in the hash table
func (ht *HashTable) Contains(key string) bool {
	_, exists := ht.Get(key)
	return exists
}

// Size returns the number of key-value pairs in the hash table
func (ht *HashTable) Size() int {
	return ht.size
}

// IsEmpty checks if the hash table is empty
func (ht *HashTable) IsEmpty() bool {
	return ht.size == 0
}

// Keys returns a slice of all keys in the hash table
func (ht *HashTable) Keys() []string {
	keys := make([]string, 0, ht.size)
	
	// Iterate through all buckets
	for i := 0; i < ht.capacity; i++ {
		current := ht.buckets[i]
		// Traverse each chain
		for current != nil {
			keys = append(keys, current.key)
			current = current.next
		}
	}
	
	return keys
}

// Values returns a slice of all values in the hash table
func (ht *HashTable) Values() []interface{} {
	values := make([]interface{}, 0, ht.size)
	
	for i := 0; i < ht.capacity; i++ {
		current := ht.buckets[i]
		for current != nil {
			values = append(values, current.value)
			current = current.next
		}
	}
	
	return values
}

// LoadFactor returns the current load factor (size/capacity)
// Higher load factor indicates higher probability of collisions
func (ht *HashTable) LoadFactor() float64 {
	return float64(ht.size) / float64(ht.capacity)
}

// String provides a string representation of the hash table
// Useful for debugging and visualization
func (ht *HashTable) String() string {
	var buffer bytes.Buffer
	
	buffer.WriteString(fmt.Sprintf("HashTable (size: %d, capacity: %d, load factor: %.2f)\n", 
		ht.size, ht.capacity, ht.LoadFactor()))
	
	for i := 0; i < ht.capacity; i++ {
		buffer.WriteString(fmt.Sprintf("Bucket %d: ", i))
		
		current := ht.buckets[i]
		if current == nil {
			buffer.WriteString("empty")
		} else {
			for current != nil {
				buffer.WriteString(fmt.Sprintf("[%s: %v]", current.key, current.value))
				if current.next != nil {
					buffer.WriteString(" -> ")
				}
				current = current.next
			}
		}
		buffer.WriteString("\n")
	}
	
	return buffer.String()
}

// Resize increases the capacity of the hash table and rehashes all elements
// This helps maintain good performance by keeping load factor low
func (ht *HashTable) Resize(newCapacity int) {
	if newCapacity <= ht.capacity {
		return // Only allow resizing to larger capacity
	}
	
	// Create new buckets with increased capacity
	newBuckets := make([]*HashTableNode, newCapacity)
	oldBuckets := ht.buckets
	oldCapacity := ht.capacity
	
	// Update hash table properties
	ht.buckets = newBuckets
	ht.capacity = newCapacity
	ht.size = 0 // Reset size, will be rebuilt during rehashing
	
	// Rehash all existing elements
	for i := 0; i < oldCapacity; i++ {
		current := oldBuckets[i]
		for current != nil {
			ht.Put(current.key, current.value)
			current = current.next
		}
	}
}

// Clear removes all elements from the hash table
func (ht *HashTable) Clear() {
	ht.buckets = make([]*HashTableNode, ht.capacity)
	ht.size = 0
}

// BucketStats returns statistics about bucket distribution
// Useful for analyzing the effectiveness of the hash function
func (ht *HashTable) BucketStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	emptyBuckets := 0
	maxChainLength := 0
	totalChainLength := 0
	nonEmptyBuckets := 0
	
	for i := 0; i < ht.capacity; i++ {
		chainLength := 0
		current := ht.buckets[i]
		
		for current != nil {
			chainLength++
			current = current.next
		}
		
		if chainLength == 0 {
			emptyBuckets++
		} else {
			nonEmptyBuckets++
			totalChainLength += chainLength
			if chainLength > maxChainLength {
				maxChainLength = chainLength
			}
		}
	}
	
	stats["empty_buckets"] = emptyBuckets
	stats["non_empty_buckets"] = nonEmptyBuckets
	stats["max_chain_length"] = maxChainLength
	if nonEmptyBuckets > 0 {
		stats["avg_chain_length"] = float64(totalChainLength) / float64(nonEmptyBuckets)
	} else {
		stats["avg_chain_length"] = 0.0
	}
	stats["load_factor"] = ht.LoadFactor()
	
	return stats
}

func Main_hashtable_chaining() {
	// Create a new hash table with initial capacity of 5
	ht := NewHashTable(5)
	
	fmt.Printf("=== HashTable Implementation using Separate Chaining ===\n\n")
	
	// Demonstrate basic operations
	fmt.Println("1. Adding elements to the hash table:")
	ht.Put("apple", 10)
	ht.Put("banana", 20)
	ht.Put("orange", 30)
	ht.Put("grape", 40)
	ht.Put("kiwi", 50)
	
	fmt.Println(ht.String())
	
	// Demonstrate collision handling
	fmt.Println("2. Adding more elements (will cause collisions):")
	ht.Put("melon", 60)
	ht.Put("pear", 70)
	ht.Put("peach", 80)
	
	fmt.Println(ht.String())
	
	// Demonstrate retrieval
	fmt.Println("3. Retrieving values:")
	if value, exists := ht.Get("banana"); exists {
		fmt.Printf("Value for 'banana': %v\n", value)
	}
	
	if value, exists := ht.Get("watermelon"); exists {
		fmt.Printf("Value for 'watermelon': %v\n", value)
	} else {
		fmt.Println("'watermelon' not found in hash table")
	}
	
	// Demonstrate update
	fmt.Println("\n4. Updating existing key:")
	ht.Put("apple", 100) // Update existing key
	if value, exists := ht.Get("apple"); exists {
		fmt.Printf("Updated value for 'apple': %v\n", value)
	}
	
	// Demonstrate deletion
	fmt.Println("\n5. Deleting a key:")
	if ht.Delete("orange") {
		fmt.Println("Successfully deleted 'orange'")
	}
	fmt.Println(ht.String())
	
	// Demonstrate utility methods
	fmt.Println("6. Utility methods:")
	fmt.Printf("All keys: %v\n", ht.Keys())
	fmt.Printf("Hash table size: %d\n", ht.Size())
	fmt.Printf("Is empty: %t\n", ht.IsEmpty())
	fmt.Printf("Contains 'grape': %t\n", ht.Contains("grape"))
	
	// Demonstrate statistics
	fmt.Println("\n7. Hash table statistics:")
	stats := ht.BucketStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}
	
	// Demonstrate resizing
	fmt.Println("\n8. Resizing the hash table:")
	ht.Resize(10)
	fmt.Println(ht.String())
	
	// Demonstrate clearing
	fmt.Println("9. Clearing the hash table:")
	ht.Clear()
	fmt.Printf("Size after clear: %d, Is empty: %t\n", ht.Size(), ht.IsEmpty())
}











