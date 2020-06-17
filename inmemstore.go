package main

import (
	"bufio"
	"fmt"
	"github.com/thedeveloperr/redis-clone/hashmap"
	"github.com/thedeveloperr/redis-clone/sortedSetMap"
	"log"
	"os"
	"strconv"
	"time"
)

// Struct for handling of appending commands to AOF file
type AOFPersistor struct {
	queue  chan string
	ticker *time.Ticker
}

// Append string commands to the file with given filename
// runs fsync to make sure data reaches hard disk and reduce
// chance of data loss.
func FlushCommands(command string, filename string) error {
	if filename == "" {
		return nil
	}
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(command + "\n"); err != nil {
		return err
	}
	return f.Sync()
}

// Struct for the main In Memory db
type InMemoryStore struct {
	sortedSet     *sortedSetMap.ConcurrentSortedsetMap
	hashmap       *hashmap.ConcurrentMap
	dataPersistor *AOFPersistor
}

// First load all the data in AOF file if exists in memory
// Attach the AOF persistor to the in memory db so as to append future write
// commands
func CreateInMemStore(persistAfter int, AOFfilename string) *InMemoryStore {

	db := &InMemoryStore{
		sortedSet:     sortedSetMap.Create(),
		hashmap:       hashmap.Create(),
		dataPersistor: nil,
	}

	if AOFfilename != "" {
		file, err := os.Open(AOFfilename)
		if err == nil {

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				db.ProcessCommand(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	}

	dataPersistor := &AOFPersistor{
		ticker: time.NewTicker(time.Duration(persistAfter) * time.Second),
		queue:  make(chan string, 1000),
	}
	db.dataPersistor = dataPersistor

	// Runs every given amount of seconds
	// If the commands are there in buffered channel
	// consume them and write them.
	go func() {
		for {
			<-db.dataPersistor.ticker.C
			for elem := range db.dataPersistor.queue {
				err := FlushCommands(elem, AOFfilename)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return db
}

// Client's command is sent here, parsed and appropriate methods
// on hashmap and Ordered Set Map are called.
func (store *InMemoryStore) ProcessCommand(command string) string {
	comm := Command{
		fullText: command,
	}
	commType, key, args := comm.parse()
	switch commType {
	case "EXPIRE":
		ttl, _ := strconv.ParseInt(args[0][0], 10, 32)
		result := store.EXPIRE(key, int(ttl))
		if result != "0" && store.dataPersistor != nil {
			store.dataPersistor.queue <- command
		}
		return result
	case "GET":
		return store.GET(key)
	case "SET":
		result := store.SET(key, args[0][0])
		if result == "OK" && store.dataPersistor != nil {
			store.dataPersistor.queue <- command
		}
		return result
	case "ZRANGE":
		if len(args) == 2 {
			start, _ := strconv.ParseInt(args[0][0], 10, 64)
			end, _ := strconv.ParseInt(args[0][1], 10, 64)
			members, scores := store.ZRANGE_WITHSCR(key, start, end)
			result := ""
			num := 1
			for i := 0; i < len(members); i++ {
				result += strconv.Itoa(num) + ") " + "'" + members[i] + "'" + "\n"
				num++

				result += strconv.Itoa(num) + ") " + fmt.Sprintf("%g", scores[i]) + "\n"
				num++
			}
			if result == "" {
				return "(empty list or set)"
			}
			return result

		}
		if len(args) == 1 {
			start, _ := strconv.ParseInt(args[0][0], 10, 64)
			end, _ := strconv.ParseInt(args[0][1], 10, 64)
			members := store.ZRANGE(key, start, end)
			result := ""
			num := 1
			for i := 0; i < len(members); i++ {
				result += strconv.Itoa(num) + ") " + "'" + members[i] + "'" + "\n"
				num++
			}
			if result == "" {
				return "(empty list or set)"
			}
			return result
		}
	case "ZRANK":
		return store.ZRANK(key, args[0][0])
	case "ZADD":
		added := 0
		for i := 0; i < len(args); i++ {
			score, _ := strconv.ParseFloat(args[i][0], 64)
			added += store.ZADD(key, score, args[i][1])
		}
		if added > 0 && store.dataPersistor != nil {
			store.dataPersistor.queue <- command
		}
		result := fmt.Sprintf("%d", added)
		return result
	}
	return "COMMAND NOT VALID"
}

// Gets value of key if set otherwise (nil)
func (store *InMemoryStore) GET(key string) string {
	if val, exists := store.hashmap.Get(key); exists {
		return val
	}
	return "(nil)"
}

// Sets value of key returns "OK" if successful
func (store *InMemoryStore) SET(key string, value string) string {
	store.hashmap.Set(key, value)
	return "OK"
}

// Perform ZADD and Inserts a member element with a given score in sorted set backed by Skiplist and Hasmap
func (store *InMemoryStore) ZADD(key string, score float64, member string) int {
	return store.sortedSet.Add(key, member, score)
}

// Gets the members and Scores stored in Sorted Set within given range. Perform ZRANGE key start end WITHSCORES command
func (store *InMemoryStore) ZRANGE_WITHSCR(key string, start int64, end int64) (members []string, scores []float64) {
	members, scores = store.sortedSet.GetMembersAndScoreInRange(key, start, end)
	return
}

// Gets the members stored in Sorted Set within given range. Perform ZRANGE key start end command
func (store *InMemoryStore) ZRANGE(key string, start int64, end int64) (members []string) {
	members, _ = store.sortedSet.GetMembersAndScoreInRange(key, start, end)
	return
}

// Gets position of member inside sorted set. Perform ZRANK. It's 0 index based
func (store *InMemoryStore) ZRANK(key string, member string) string {
	if rank, exists := store.sortedSet.GetRank(key, member); exists {
		return fmt.Sprintf("%d", rank)
	}
	return "(nil)"
}

// Expire and remove key after some given ttl seconds. Perform EXPIRE key ttl command
func (store *InMemoryStore) EXPIRE(key string, ttl int) string {
	canExpire := store.hashmap.Expire(key, ttl)
	if canExpire == 1 {
		return "1"
	}

	canExpire = store.sortedSet.Expire(key, ttl)
	if canExpire == 1 {
		return "1"
	}
	return "0"
}
