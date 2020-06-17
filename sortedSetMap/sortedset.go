package sortedSetMap

import (
	"math/rand"
	"sync"
	"time"
)

const MAX_LEVEL = 25

type Level struct {
	nextNode         *SkiplistNode // next node
	distanceNextNode uint64        // distance to next node. Used for indexable skiplist
}

type SkiplistNode struct {
	member string
	score  float64
	levels []*Level // Array of Level stacked up on the Node
}

type Skiplist struct {
	header, tail *SkiplistNode // Start and end
	length       uint64        // number of nodes
	level        uint          // level at which current list is at
	randomSource rand.Source
}

func CreateSkiplist() *Skiplist {
	result := &Skiplist{
		header: &SkiplistNode{
			score:  0.0,
			member: "",
			levels: make([]*Level, MAX_LEVEL),
		},

		level:        1,
		length:       0,
		randomSource: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// initialise each level to empty list
	for j := 0; j < MAX_LEVEL; j++ {
		result.header.levels[j] = &Level{
			nextNode:         nil,
			distanceNextNode: 0,
		}
	}

	result.tail = nil
	return result
}

func CreateSkiplistNode(level uint, score float64, member string) *SkiplistNode {
	result := &SkiplistNode{
		score:  score,
		member: member,
		levels: make([]*Level, level),
	}
	// initialise each level to empty list
	for j := 0; j < int(level); j++ {
		result.levels[j] = &Level{
			nextNode:         nil,
			distanceNextNode: 0,
		}
	}

	return result
}

func (list *Skiplist) randomLevel() uint {

	// Float64() for some reason not buidling so used Int63 to get Float64
	r := float64(list.randomSource.Int63()) / (1 << 63)

	var level uint = 1

	// 1/2 ratio between nodes at each level for skiplist
	for level < MAX_LEVEL && r < 0.5 {
		level++
	}
	return level
}

// Returns previous nodes to insert after at every level
// calculate ranking and postion info
func (list *Skiplist) GetPreviousNodesAndRanks(score float64, member string) ([MAX_LEVEL]*SkiplistNode, [MAX_LEVEL]uint64) {
	var previousNodes [MAX_LEVEL]*SkiplistNode
	var rank [MAX_LEVEL]uint64
	iteratorNode := list.header

	// Start from top level and move down and right
	for i := int(list.level) - 1; i >= 0; i-- {
		level := iteratorNode.levels[i]
		if i == int(list.level)-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] // Get above level rank data
		}
		for (level.nextNode != nil) &&
			((level.nextNode.score < score) ||
				(level.nextNode.score == score && level.nextNode.member < member)) { // ZADD needs lexograohical sorting if score is same
			rank[i] += level.distanceNextNode
			iteratorNode = level.nextNode
			level = iteratorNode.levels[i]
		}
		// Found the correct node to insert after
		// at i level, continue to downward level
		previousNodes[i] = iteratorNode
	}
	return previousNodes, rank
}

func (list *Skiplist) Insert(score float64, member string) {
	previousNodes, rank := list.GetPreviousNodesAndRanks(score, member)

	randomLevelForNode := list.randomLevel()
	if randomLevelForNode > list.level {
		// increase list levels
		level := list.level
		for level < randomLevelForNode {
			rank[level] = 0
			previousNodes[level] = list.header
			previousNodes[level].levels[level].distanceNextNode = list.length
			level++
		}
		list.level = randomLevelForNode
	}

	nodeToInsert := CreateSkiplistNode(randomLevelForNode, score, member)
	var level_i uint = 0
	for level_i < randomLevelForNode {
		// insert node by swapping links
		nodeToInsert.levels[level_i].nextNode = previousNodes[level_i].levels[level_i].nextNode
		previousNodes[level_i].levels[level_i].nextNode = nodeToInsert

		// Now the length, distanceNextNode etc. changed, make necessarychanges
		nodeToInsert.levels[level_i].distanceNextNode = previousNodes[level_i].levels[level_i].distanceNextNode - (rank[0] - rank[level_i])
		previousNodes[level_i].levels[level_i].distanceNextNode = rank[0] - rank[level_i] + 1
		level_i++
	}

	// if new levels were added fix the distanceNextNode for them
	level_i = randomLevelForNode
	for level_i < list.level {
		// increment by one as a single node is added to right
		previousNodes[level_i].levels[level_i].distanceNextNode++
		level_i++
	}
	// if inserted at end
	if nodeToInsert.levels[0].nextNode == nil {
		list.tail = nodeToInsert
	}
	list.length++
}

// 0 if no rank found. Rank start from 1
func (list *Skiplist) GetRank(score float64, member string) uint64 {
	var rank uint64 = 0
	iteratorNode := list.header

	// Start from top level and move down and right
	for i := int(list.level) - 1; i >= 0; i-- {
		level := iteratorNode.levels[i]
		for (level.nextNode != nil) &&
			((level.nextNode.score < score) ||
				(level.nextNode.score == score && level.nextNode.member <= member)) {
			rank += level.distanceNextNode
			iteratorNode = level.nextNode
			level = iteratorNode.levels[i]
		}
		if iteratorNode.member == member {
			return rank
		}
	}

	// at level 0 and might have found correct node
	return 0
}

// 1 indexed rank
func (list *Skiplist) GetNodeAtRank(rank uint64) *SkiplistNode {
	var distanceTravelled uint64 = 0
	iteratorNode := list.header

	// Start from top level and move down and right
	for i := int(list.level) - 1; i >= 0; i-- {
		level := iteratorNode.levels[i]
		for (level.nextNode != nil) &&
			(distanceTravelled+level.distanceNextNode <= rank) {
			distanceTravelled += level.distanceNextNode
			iteratorNode = level.nextNode
			level = iteratorNode.levels[i]
		}
	}
	if rank == distanceTravelled {
		return iteratorNode
	}

	return nil
}

// pos is 0 based
func (list *Skiplist) GetMembersAndScoreInRange(posStart int64, posEnd int64) (members []string, scores []float64) {
	correctedStartPos := posStart
	correctedEndPos := posEnd
	if posStart < 0 {
		correctedStartPos = posStart + int64(list.length)
	}

	if posEnd < 0 {
		correctedEndPos = posEnd + int64(list.length)
	}

	if correctedStartPos < 0 {
		correctedStartPos = 0
	}

	if correctedStartPos > correctedEndPos {
		return members, scores
	}

	i := correctedStartPos
	iteratorNode := list.GetNodeAtRank(uint64(correctedStartPos + 1))
	for i <= correctedEndPos && iteratorNode != nil {
		members = append(members, iteratorNode.member)
		scores = append(scores, iteratorNode.score)
		iteratorNode = iteratorNode.levels[0].nextNode
		i++
	}

	return members, scores
}

type SortedsetMap interface {
	GetRank(key string, member string) uint64
	Add(key string, member string, score float64) int
	GetMembersAndScoreInRange(key string, start int64, end int64) ([]string, []float64)
	Expire(key string, timeoutSeconds int) int
}

type Value struct {
	value        *Sortedset
	setAt        time.Time
	expireAfter  time.Duration
	shouldExpire bool
}

type Sortedset struct {
	memberScoreMap map[string]float64
	skiplist       *Skiplist
}

type ConcurrentSortedsetMap struct {
	mutex sync.RWMutex
	data  map[string]*Value
}

func Create() *ConcurrentSortedsetMap {
	sortedsetmap := &ConcurrentSortedsetMap{
		data: make(map[string]*Value),
	}
	return sortedsetmap
}

func (c *ConcurrentSortedsetMap) Add(key string, member string, score float64) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	var skiplist *Skiplist
	if _, exists := c.data[key]; !exists {
		skiplist = CreateSkiplist()
		skiplist.Insert(score, member)
		sortedset := &Sortedset{
			memberScoreMap: make(map[string]float64),
			skiplist:       skiplist,
		}
		sortedset.memberScoreMap[member] = score
		c.data[key] = &Value{
			value:        sortedset,
			setAt:        time.Now(),
			expireAfter:  0,
			shouldExpire: false,
		}
		return 1
	}

	if _, exists := c.data[key].value.memberScoreMap[member]; exists {
		return 0
	}

	skiplist = c.data[key].value.skiplist
	c.data[key].value.memberScoreMap[member] = score
	skiplist.Insert(score, member)
	return 1
}

func (c *ConcurrentSortedsetMap) GetRank(key string, member string) (uint64, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	valueItem, exists := c.GetUnsafe(key)
	if !exists {
		return 0, false
	}
	score := valueItem.value.memberScoreMap[member]
	rank := valueItem.value.skiplist.GetRank(score, member)
	if rank == 0 {
		return 0, false
	}
	return rank - 1, true
}

func (c *ConcurrentSortedsetMap) GetMembersAndScoreInRange(key string, start int64, end int64) (members []string, scores []float64) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	valueItem, exists := c.GetUnsafe(key)
	if !exists {
		return members, scores
	}
	members, scores = valueItem.value.skiplist.GetMembersAndScoreInRange(start, end)
	return members, scores
}

func (c *ConcurrentSortedsetMap) GetUnsafe(key string) (*Value, bool) {

	valueItem, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// To improve accuracy of EXPIRE, in case time.AfterFunc runs later and a get call is made earlier
	if valueItem.shouldExpire && time.Now().Sub(valueItem.setAt) > valueItem.expireAfter {
		return nil, false
	}
	return valueItem, exists
}

func (c *ConcurrentSortedsetMap) Expire(key string, timeoutSeconds int) int {
	c.mutex.Lock()
	if val, ok := c.GetUnsafe(key); !ok {
		_ = val
		return 0
	}
	c.data[key].shouldExpire = true
	c.data[key].expireAfter = time.Duration(timeoutSeconds) * time.Second
	c.mutex.Unlock()
	time.AfterFunc(time.Duration(timeoutSeconds)*time.Second, func() {
		c.mutex.Lock()

		delete(c.data, key)
		c.mutex.Unlock()
	})
	return 1
}
