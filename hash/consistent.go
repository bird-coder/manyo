/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2022-06-28 16:04:53
 * @LastEditTime: 2022-07-11 15:15:32
 * @LastEditors: yuanshisan
 */
package hash

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

var ErrEmptyCircle = errors.New("empty circle")

type Node struct {
	Id      int
	Addr    string
	numReps int
}

const defaultReplicas = 4

func NewNode(id int, addr string, num int) *Node {
	if num <= 0 {
		num = defaultReplicas
	}
	return &Node{
		Id:      id,
		Addr:    addr,
		numReps: num,
	}
}

type Consistent struct {
	circle  map[uint32]string
	members map[string]bool
	count   int
	ring    HashRing
	sync.RWMutex
}

func NewConsistent() *Consistent {
	c := new(Consistent)
	c.circle = make(map[uint32]string)
	c.members = make(map[string]bool)
	return c
}

func (c *Consistent) AddNode(node *Node) {
	c.Lock()
	defer c.Unlock()

	c.add(node)
}

func (c *Consistent) add(node *Node) {
	if _, ok := c.members[node.Addr]; ok {
		return
	}

	for i := 0; i < node.numReps; i++ {
		c.circle[c.hashKey(c.eltKey(i, node))] = node.Addr
	}
	c.members[node.Addr] = true
	c.sortHashRing()
	c.count++
}

func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	c.remove(node)
}

func (c *Consistent) remove(node *Node) {
	if _, ok := c.members[node.Addr]; !ok {
		return
	}

	for i := 0; i < node.numReps; i++ {
		delete(c.circle, c.hashKey(c.eltKey(i, node)))
	}
	delete(c.members, node.Addr)
	c.sortHashRing()
	c.count--
}

// func (c *Consistent) Set(nodes []*Node) {
// 	c.Lock()
// 	defer c.Unlock()

// 	for k := range c.members {
// 		found := false
// 		for _, node := range nodes {
// 			if k == node.Addr {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			c.remove(nodes[0])
// 		}
// 	}
// 	for _, node := range nodes {
// 		c.add(node)
// 	}
// }

func (c *Consistent) Members() []string {
	c.RLock()
	defer c.RUnlock()

	var m []string
	for k := range c.members {
		m = append(m, k)
	}
	return m
}

func (c *Consistent) PickNode(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.ring) == 0 {
		return "", ErrEmptyCircle
	}

	hash := c.hashKey(key)
	i := c.search(hash)

	return c.circle[c.ring[i]], nil
}

func (c *Consistent) search(hash uint32) int {
	n := len(c.ring)
	i := sort.Search(n, func(i int) bool {
		return c.ring[i] >= hash
	})
	if i >= n {
		return 0
	}
	return i
}

func (c *Consistent) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.circle {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

func (c *Consistent) eltKey(idx int, node *Node) string {
	return node.Addr + "|" + strconv.Itoa(idx)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], []byte(key))
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}
