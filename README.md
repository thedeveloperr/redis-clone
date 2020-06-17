# Redis clone.

## Steps to run
1. Install go version 1.14 You can see here how to install: https://golang.org/doc/install
2. Clone the repo. `git clone git@github.com:thedeveloperr/redis-clone.git`  or  `git clone https://github.com/thedeveloperr/redis-clone.git`
3. In the root directory run:
   ```go run ./```
4. On Mac if any popup asking for "Do you want the application “redis-clone” to accept incoming network connections?" click Yes
5. Server will start running at ```http://localhost:8080/```


## Steps to run commands
1. HTTP posts are used to send commands to the server
2. Eg. Copy and paste the following command while server is running `curl -d "command=SET edtech=awesome" http://localhost:8080/`
3. Similarly run other commands just pass the commands as POST data `command=GET edtech` that is: `curl -d "command=GET edtech" http://localhost:8080/` 
4. You can close the server too and rerun the program and send the HTTP command GET edtech again to see the last set valued. This is done by simulating redis's `Append Only File Persistance` technique.

## Steps to run test
1. Go to the root of repo.
2. Run `go test`.

## NOTE:
- Don't delete AOF_test_read.log as it's req. for automated testing.

## Code base overview
- Written in wiki: https://github.com/thedeveloperr/redis-clone/wiki


## Decisions

### Why Golang ?
   Golang: It's a type safe compiled language with first class support for concurrency. Easy concurrency via goroutines. Low memory footprint and less verbose than Java. Concurrency can improve throughput. Golang seemed right tool for job.

### What are the further improvements that can be made to make it efficient ?
  Future Improvements:-
  - Right now AOF file persistance (similar to what redis does) is rudimentary and can grow large as it's append only. So will need to add some techniques to rewrite AOF just like redis do once the file reaches certain size.
  - Many commands are missing and only following commands are there:
    - GET, SET, ZRANK, ZADD, ZRANGE, EXPIRE are the only supported commands right now

  - Stress testing and benchmarking can further provide insights into bottlenecks
  - Concurrency for Data structures like SkipList used in ordered set can be further improved by sharding/bucketing the write request and locking that bucket only to reduce lock contention when a write is happening.


### Data structures used and why?
  - Thread safe Hashmap: For basic operations like GET SET and for maintaining inner mapping of score and members in ordered set. Also making sure key already exists or not efficiently.
    * All these things can be done in Avg. O(1) time.
    * Golang doesn't have a map which provide thread safety for both read and write (sync.Map is optimised for Read and suffers on repeated write). Used sync.RWMutex to implement thread safe Map.
  - Thread safe Skiplist: SortedSet etc. are usually implemented using LinkedList or BalancedTrees etc. but to make Insert (ZADD), and Query (ZRANGE and ZRANK) happens in order O(log(N)) a different datastructre is needed.
  - Skiplist does Insert, Search etc. All in avg. O(log(N))

### Does it supports multithreading ?
 - Supports Multithreading via Goroutines. Used thread safe data structures via RWMutex as it [solves Reader Writer Problem](https://en.wikipedia.org/wiki/Readers%E2%80%93writer_lock). 
 - For better Concurrent Reading.
 - Write on Hashmap Doesn't block write or read on Sorted Set and vice verca. Eg. ZADD doesn't blocks GET or SET.
 - For background deletion of key after n Seconds timeout of EXPIRE command requires threading so other read operations on other data structures keep on happening.
