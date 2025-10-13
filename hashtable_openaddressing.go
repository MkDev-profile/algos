package main

import (
    "fmt"
    "hash/fnv"
)

// Entry represents a key-value pair in the hashtable
// We need to track if the slot was previously occupied (tombstone) for proper probing
type Entry struct {
    key   string
    value interface{}
    // occupied indicates if this slot currently holds a valid entry
    occupied bool
    // tombstone indicates if this slot was previously occupied but now deleted
    // This is crucial for open addressing to know when to stop probing
    tombstone bool
}

// Hashtable implements a hash table using open addressing collision resolution
type Hashtable struct {
    // buckets holds our key-value pairs
    buckets []Entry
    // size is the current number of elements in the hashtable
    size int
    // capacity is the total number of slots available
    capacity int
    // loadFactorThreshold determines when to resize the hashtable
    // When load factor (size/capacity) exceeds this, we resize
    loadFactorThreshold float64
}

// NewHashtable creates a new hashtable with initial capacity
func NewHashtable(initialCapacity int) *Hashtable {
    if initialCapacity <= 0 {
        initialCapacity = 16 // Default initial capacity
    }
    
    // Initialize all buckets as empty
    buckets := make([]Entry, initialCapacity)
    for i := range buckets {
        buckets[i] = Entry{occupied: false, tombstone: false}
    }
    
    return &Hashtable{
        buckets:             buckets,
        size:                0,
        capacity:            initialCapacity,
        loadFactorThreshold: 0.75, // Common load factor threshold
    }
}

// hashFunction computes the hash value for a given key
// We use FNV-1a hash function for good distribution
func (h *Hashtable) hashFunction(key string) int {
    hasher := fnv.New32a()
    hasher.Write([]byte(key))
    // Convert hash to positive integer and mod by capacity
    return int(hasher.Sum32()) % h.capacity
}

// probe returns the next index to check during probing
// This implements linear probing. Other strategies include:
// - Quadratic probing: (hash + i²) % capacity
// - Double hashing: (hash1 + i * hash2) % capacity
func (h *Hashtable) probe(initialIndex, attempt int) int {
    // Linear probing: check next consecutive slots
    return (initialIndex + attempt) % h.capacity
}

// findSlot finds the appropriate slot for a key using open addressing
// Returns the index where the key is found or where it should be inserted
// Also returns whether the key was found
func (h *Hashtable) findSlot(key string) (int, bool) {
    // Step 1: Compute initial hash index
    initialIndex := h.hashFunction(key)
    index := initialIndex
    attempt := 0
    
    // Track the first tombstone we encounter for potential insertion
    firstTombstone := -1
    
    // Step 2: Probe until we find an empty slot or the key
    for attempt < h.capacity {
        currentEntry := &h.buckets[index]
        
        // Case 1: Found the exact key (occupied and matching key)
        if currentEntry.occupied && currentEntry.key == key {
            return index, true
        }
        
        // Case 2: Found a tombstone - remember it for potential insertion
        if !currentEntry.occupied && currentEntry.tombstone && firstTombstone == -1 {
            firstTombstone = index
        }
        
        // Case 3: Found an empty slot (not tombstone) - search ends here
        if !currentEntry.occupied && !currentEntry.tombstone {
            // If we found a tombstone earlier, use that slot to maintain clustering
            if firstTombstone != -1 {
                return firstTombstone, false
            }
            return index, false
        }
        
        // Continue probing
        attempt++
        index = h.probe(initialIndex, attempt)
    }
    
    // If we've searched all slots and found no empty slot (shouldn't happen if resized properly)
    // But if it does, return the first tombstone or last checked index
    if firstTombstone != -1 {
        return firstTombstone, false
    }
    return index, false
}

// loadFactor calculates the current load factor of the hashtable
func (h *Hashtable) loadFactor() float64 {
    return float64(h.size) / float64(h.capacity)
}

// resize creates a new larger hashtable and rehashes all elements
// This is necessary to maintain performance as the table fills up
func (h *Hashtable) resize() {
    // Double the capacity (common growth strategy)
    newCapacity := h.capacity * 2
    if newCapacity < 16 {
        newCapacity = 16
    }
    
    // Create new temporary hashtable with larger capacity
    newBuckets := make([]Entry, newCapacity)
    for i := range newBuckets {
        newBuckets[i] = Entry{occupied: false, tombstone: false}
    }
    
    oldBuckets := h.buckets
    oldCapacity := h.capacity
    
    // Replace current buckets with new empty ones
    h.buckets = newBuckets
    h.capacity = newCapacity
    h.size = 0 // Reset size, will be rebuilt during rehashing
    
    // Rehash all valid entries from old buckets to new buckets
    for i := 0; i < oldCapacity; i++ {
        if oldBuckets[i].occupied {
            // Use Put method to rehash the entry
            h.Put(oldBuckets[i].key, oldBuckets[i].value)
        }
    }
    
    fmt.Printf("Resized hashtable: %d -> %d capacity, %d elements\n", 
        oldCapacity, newCapacity, h.size)
}

// Put inserts or updates a key-value pair in the hashtable
func (h *Hashtable) Put(key string, value interface{}) {
    // Check if resize is needed before insertion
    if h.loadFactor() >= h.loadFactorThreshold {
        h.resize()
    }
    
    // Find the appropriate slot for this key
    index, found := h.findSlot(key)
    
    if found {
        // Key exists - update the value
        h.buckets[index].value = value
    } else {
        // Key doesn't exist - insert new entry
        h.buckets[index] = Entry{
            key:      key,
            value:    value,
            occupied: true,
            tombstone: false,
        }
        h.size++
    }
}

// Get retrieves the value for a key, returns value and whether key exists
func (h *Hashtable) Get(key string) (interface{}, bool) {
    index, found := h.findSlot(key)
    if found {
        return h.buckets[index].value, true
    }
    return nil, false
}

// Delete removes a key-value pair from the hashtable
// Uses tombstone marking for proper probing behavior
func (h *Hashtable) Delete(key string) bool {
    index, found := h.findSlot(key)
    if found {
        // Mark as tombstone instead of truly emptying the slot
        // This ensures that probing for other keys won't break
        h.buckets[index].occupied = false
        h.buckets[index].tombstone = true
        h.buckets[index].key = "" // Clear key for garbage collection
        h.buckets[index].value = nil // Clear value for garbage collection
        h.size--
        return true
    }
    return false
}

// Contains checks if a key exists in the hashtable
func (h *Hashtable) Contains(key string) bool {
    _, found := h.findSlot(key)
    return found
}

// Size returns the number of elements in the hashtable
func (h *Hashtable) Size() int {
    return h.size
}

// Keys returns all keys in the hashtable
func (h *Hashtable) Keys() []string {
    keys := make([]string, 0, h.size)
    for i := 0; i < h.capacity; i++ {
        if h.buckets[i].occupied {
            keys = append(keys, h.buckets[i].key)
        }
    }
    return keys
}

// Values returns all values in the hashtable
func (h *Hashtable) Values() []interface{} {
    values := make([]interface{}, 0, h.size)
    for i := 0; i < h.capacity; i++ {
        if h.buckets[i].occupied {
            values = append(values, h.buckets[i].value)
        }
    }
    return values
}

// String provides a string representation of the hashtable for debugging
func (h *Hashtable) String() string {
    result := "Hashtable:\n"
    for i := 0; i < h.capacity; i++ {
        entry := h.buckets[i]
        status := "Empty"
        if entry.occupied {
            status = fmt.Sprintf("Occupied: %v -> %v", entry.key, entry.value)
        } else if entry.tombstone {
            status = "Tombstone"
        }
        result += fmt.Sprintf("  [%d]: %s\n", i, status)
    }
    result += fmt.Sprintf("Size: %d, Capacity: %d, Load Factor: %.2f\n", 
        h.size, h.capacity, h.loadFactor())
    return result
}

// Example usage and demonstration
func Main_hashtable_openaddressing() {
    // Create a new hashtable
    ht := NewHashtable(8)
    fmt.Println("Initial hashtable:")
    fmt.Println(ht)
    
    // Insert some key-value pairs
    fmt.Println("Inserting key-value pairs...")
    ht.Put("name", "Alice")
    ht.Put("age", 30)
    ht.Put("city", "New York")
    ht.Put("occupation", "Engineer")
    ht.Put("language", "Go")
    
    fmt.Println("After insertions:")
    fmt.Println(ht)
    
    // Demonstrate retrieval
    fmt.Println("Retrieval examples:")
    if value, found := ht.Get("name"); found {
        fmt.Printf("Found 'name': %v\n", value)
    }
    
    if value, found := ht.Get("age"); found {
        fmt.Printf("Found 'age': %v\n", value)
    }
    
    // Demonstrate update
    fmt.Println("\nUpdating 'age' to 31...")
    ht.Put("age", 31)
    if value, found := ht.Get("age"); found {
        fmt.Printf("Updated 'age': %v\n", value)
    }
    
    // Demonstrate deletion
    fmt.Println("\nDeleting 'city'...")
    ht.Delete("city")
    fmt.Printf("After deletion - Contains 'city': %v\n", ht.Contains("city"))
    fmt.Println(ht)
    
    // Demonstrate collision handling
    fmt.Println("Demonstrating collision handling...")
    
    // These keys might hash to the same index (depending on hash function)
    ht.Put("a", "value_a")
    ht.Put("b", "value_b") 
    ht.Put("c", "value_c")
    
    fmt.Println("After adding potential collision keys:")
    fmt.Println(ht)
    
    // Show all keys and values
    fmt.Printf("All keys: %v\n", ht.Keys())
    fmt.Printf("All values: %v\n", ht.Values())
    fmt.Printf("Total size: %d\n", ht.Size())
    
    // Demonstrate resizing by adding more elements
    fmt.Println("\nAdding more elements to trigger resize...")
    for i := 0; i < 10; i++ {
        key := fmt.Sprintf("key%d", i)
        ht.Put(key, i*100)
    }
    
    fmt.Println("After potential resize:")
    fmt.Println(ht)
}










