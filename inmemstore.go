package main

import (
	"fmt"
	"github.com/thedeveloperr/redis-clone/hashmap"
	"github.com/thedeveloperr/redis-clone/sortedSetMap"
	"log"
	"os"
	"strconv"
	"time"
)

type AOFPersistor struct {
	queue  chan string
	ticker *time.Ticker
}

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

type InMemoryStore struct {
	sortedSet     *sortedSetMap.ConcurrentSortedsetMap
	hashmap       *hashmap.ConcurrentMap
	dataPersistor *AOFPersistor
}

func CreateInMemStore(persistAfter int, AOFfilename string) *InMemoryStore {
	dataPersistor := &AOFPersistor{
		ticker: time.NewTicker(time.Duration(persistAfter) * time.Second),
		queue:  make(chan string, 1000),
	}

	go func() {
		for {
			<-dataPersistor.ticker.C
			for elem := range dataPersistor.queue {
				err := FlushCommands(elem, AOFfilename)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return &InMemoryStore{
		sortedSet:     sortedSetMap.Create(),
		hashmap:       hashmap.Create(),
		dataPersistor: dataPersistor,
	}
}

func (store *InMemoryStore) ProcessCommand(command string) string {
	comm := Command{
		fullText: command,
	}
	commType, key, args := comm.parse()
	switch commType {
	case "EXPIRE":
		ttl, _ := strconv.ParseInt(args[0][0], 10, 32)
		result := store.EXPIRE(key, int(ttl))
		if result != "0" {
			store.dataPersistor.queue <- command
		}
		return result
	case "GET":
		return store.GET(key)
	case "SET":
		result := store.SET(key, args[0][0])
		if result == "OK" {
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
		if added > 0 {
			store.dataPersistor.queue <- command
		}
		result := fmt.Sprintf("%d", added)
		return result
	}
	return "COMMAND NOT VALID"
}

func (store *InMemoryStore) GET(key string) string {
	if val, exists := store.hashmap.Get(key); exists {
		return val
	}
	return "(nil)"
}

func (store *InMemoryStore) SET(key string, value string) string {
	store.hashmap.Set(key, value)
	return "OK"
}

func (store *InMemoryStore) ZADD(key string, score float64, member string) int {
	return store.sortedSet.Add(key, member, score)
}

func (store *InMemoryStore) ZRANGE_WITHSCR(key string, start int64, end int64) (members []string, scores []float64) {
	members, scores = store.sortedSet.GetMembersAndScoreInRange(key, start, end)
	return
}

func (store *InMemoryStore) ZRANGE(key string, start int64, end int64) (members []string) {
	members, _ = store.sortedSet.GetMembersAndScoreInRange(key, start, end)
	return
}

func (store *InMemoryStore) ZRANK(key string, member string) string {
	if rank, exists := store.sortedSet.GetRank(key, member); exists {
		return fmt.Sprintf("%d", rank)
	}
	return "(nil)"
}

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
