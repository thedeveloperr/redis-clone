package main

import (
	"testing"
	"time"
)

func TestSETCommand(t *testing.T) {
	db := CreateInMemStore()
	command := "SET k1 v1"
	result := db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}

	command = "SET k2 v1"
	result = db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}

	command = "SET k1 v2"
	result = db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}
}

func TestSETGETCommand(t *testing.T) {
	db := CreateInMemStore()
	command := "SET k1 v1"
	result := db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}
	command = "GET k1"
	result = db.ProcessCommand(command)
	if result != "v1" {
		t.Errorf("Got Wrong value GET for k1:" + result + " Expected: v1")
	}

	command = "SET k2 v1"
	result = db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}

	command = "GET k2"
	result = db.ProcessCommand(command)
	if result != "v1" {
		t.Errorf("Got Wrong value GET for k2:" + result + " Expected: v1")
	}

	command = "SET k1 v2"
	result = db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't run:" + command + "Got result:" + result)
	}

	command = "GET k1"
	result = db.ProcessCommand(command)
	if result != "v2" {
		t.Errorf("Got Wrong value GET for k1:" + result + " Expected: v2")
	}

	command = "GET kNone"
	result = db.ProcessCommand(command)
	if result != "(nil)" {
		t.Errorf("Got Wrong value GET for kNone:" + result + " Expected: (nil)")
	}

}

func Test_SET_EXPIRE_GET_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "SET k1 v1"
	result := db.ProcessCommand(command)
	if result != "OK" {
		t.Errorf("Couldn't r)n:" + command + " Got result:" + result)
	}
	command = "GET k1"
	result = db.ProcessCommand(command)
	if result != "v1" {
		t.Errorf("Got Wrong value GET for k1:" + result + " Expected: v1")
	}

	command = "EXPIRE knonexisting 2"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Got Wrong value EXPIRE for knonexisting:" + result + " Expected: 0")
	}

	command = "EXPIRE k1 2.3"
	result = db.ProcessCommand(command)
	if result != "COMMAND NOT VALID" {
		t.Errorf("Got Wrong value EXPIRE for k1:" + result + " Expected: COMMAND NOT VALID")
	}

	command = "EXPIRE k1 2"
	result = db.ProcessCommand(command)
	if result != "1" {
		t.Errorf("Got Wrong value EXPIRE for k1:" + result + " Expected: 1")
	}

	command = "GET k1"
	result = db.ProcessCommand(command)
	if result != "v1" {
		t.Errorf("Got Wrong value GET for k1:" + result + " Expected: v1")
	}

	time.Sleep(2 * time.Second)

	command = "GET k1"
	result = db.ProcessCommand(command)
	if result != "(nil)" {
		t.Errorf("Got Wrong value GET for k1:" + result + " Expected: (nil)")
	}
}

func Test_ZADD_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "ZADD k1 0.1 m1"
	result := db.ProcessCommand(command)
	if result != "1" {
		t.Errorf("Ran:" + command + ".Expected 1 but Got result:" + result)
	}

	command = "ZADD k1 1.1 m1"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 1 but Got result:" + result)
	}

	command = "ZADD k1 0.2 m2 3 m3 4 m4"
	result = db.ProcessCommand(command)
	if result != "3" {
		t.Errorf("Ran:" + command + ".Expected 3 but Got result:" + result)
	}

	command = "ZADD k1 0.2 m1 3 m3 4 m4"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}
}

func Test_EXPIRE_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "ZADD k1 0.1 m1"
	result := db.ProcessCommand(command)
	if result != "1" {
		t.Errorf("Ran:" + command + ".Expected 1 but Got result:" + result)
	}

	command = "ZRANK k1 m1"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}

	command = "EXPIRE k1 2"
	result = db.ProcessCommand(command)
	if result != "1" {
		t.Errorf("Got Wrong value EXPIRE for k1:" + "1" + " Expected: 1")
	}

	command = "ZRANK k1 m1"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}
	time.Sleep(2 * time.Second)
	command = "ZRANK k1 m1"
	result = db.ProcessCommand(command)
	if result != "(nil)" {
		t.Errorf("Ran:" + command + ".Expected (nil) but Got result:" + result)
	}

}

func Test_ZRANK_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "ZADD k1 0.1 m1"

	db.ProcessCommand(command)
	command = "ZADD k1 0.2 m2 3 m3 4 m4"
	db.ProcessCommand(command)

	command = "ZRANK k1 m1"
	result := db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}

	command = "ZRANK k1 m2"
	result = db.ProcessCommand(command)
	if result != "1" {
		t.Errorf("Ran:" + command + ".Expected 1 but Got result:" + result)
	}

	command = "ZRANK k1 m3"
	result = db.ProcessCommand(command)
	if result != "2" {
		t.Errorf("Ran:" + command + ".Expected 2 but Got result:" + result)
	}

	command = "ZRANK k1 m4"
	result = db.ProcessCommand(command)
	if result != "3" {
		t.Errorf("Ran:" + command + ".Expected 3 but Got result:" + result)
	}

	command = "ZRANK k1 m5"
	result = db.ProcessCommand(command)
	if result != "(nil)" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}

	command = "ZRANK k1"
	result = db.ProcessCommand(command)
	if result != "COMMAND NOT VALID" {
		t.Errorf("Ran:" + command + ".Expected COMMAND NOT VALID but Got result:" + result)
	}

	command = "ZADD k1 0.1 m0"
	result = db.ProcessCommand(command)

	command = "ZADD k1 50 m50"
	result = db.ProcessCommand(command)

	command = "ZRANK k1 m0"
	result = db.ProcessCommand(command)
	if result != "0" {
		t.Errorf("Ran:" + command + ".Expected 0 but Got result:" + result)
	}

	command = "ZRANK k1 m50"
	result = db.ProcessCommand(command)
	if result != "5" {
		t.Errorf("Ran:" + command + ".Expected 5 but Got result:" + result)
	}
}

func Test_ZRANGE_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "ZADD k1 0.1 m1"
	db.ProcessCommand(command)

	command = "ZADD k1 0.2 m2 3 m3 4 m4"
	db.ProcessCommand(command)

	command = "ZRANGE k1 0 0"
	result := db.ProcessCommand(command)
	if result != "1) 'm1'\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm1'\n ,but Got result:" + result)
	}

	command = "ZRANGE k1 1 3"
	result = db.ProcessCommand(command)
	if result != "1) 'm2'\n2) 'm3'\n3) 'm4'\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm2'\n2) 'm3'\n3) 'm4'\n but Got result:" + result)
	}

	command = "ZRANGE k1 -10 300"
	result = db.ProcessCommand(command)
	if result != "1) 'm1'\n2) 'm2'\n3) 'm3'\n4) 'm4'\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm1'\n2) 'm2'\n3) 'm3'\n4) 'm4'\n,but Got result:" + result)
	}

	command = "ZRANGE k1 -3 -2"
	result = db.ProcessCommand(command)
	if result != "1) 'm2'\n2) 'm3'\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm2'\n2) 'm3'\n,but Got result:" + result)
	}

	command = "ZRANGE k1 -10 -9"
	result = db.ProcessCommand(command)
	if result != "(empty list or set)" {
		t.Errorf("Ran:" + command + ".Expected (empty list or set) but Got result:" + result)
	}

}

func Test_ZRANGE_WITHSCORE_Command(t *testing.T) {
	db := CreateInMemStore()
	command := "ZADD k1 0.1 m1"
	db.ProcessCommand(command)

	command = "ZADD k1 0.2 m2 3 m3 4 m4"
	db.ProcessCommand(command)

	command = "ZRANGE k1 0 0 WITHSCORE"
	result := db.ProcessCommand(command)
	if result != "COMMAND NOT VALID" {
		t.Errorf("Ran:" + command + ".Expected:COMMAND NOT VALID but got result:\n" + result)
	}

	command = "ZRANGE k1 0 0 WITHSCORES"
	result = db.ProcessCommand(command)
	if result != "1) 'm1'\n2) 0.1\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm1'\n2) 0.1\n,but Got result:\n" + result)
	}

	command = "ZRANGE k1 1 3 WITHSCORES"
	result = db.ProcessCommand(command)
	if result != "1) 'm2'\n2) 0.2\n3) 'm3'\n4) 3\n5) 'm4'\n6) 4\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm2'\n2) 0.2\n3) 'm3'\n4) 3\n5) 'm4'\n6) 4\n,but Got result:\n" + result)
	}

	command = "ZRANGE k1 -10 300 WITHSCORES"
	result = db.ProcessCommand(command)
	if result != "1) 'm1'\n2) 0.1\n3) 'm2'\n4) 0.2\n5) 'm3'\n6) 3\n7) 'm4'\n8) 4\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm1'\n2) 0.1\n3) 'm2'\n4) 0.2\n5) 'm3'\n6) 3\n7) 'm4'\n8) 4\n,but Got result:\n" + result)
	}

	command = "ZRANGE k1 -3 -2 WITHSCORES"
	result = db.ProcessCommand(command)
	if result != "1) 'm2'\n2) 0.2\n3) 'm3'\n4) 3\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm2'\n2) 0.2\n3) 'm3'\n4) 3\n,but Got result:\n" + result)
	}

	command = "ZRANGE k1 -10 -9 WITHSCORES"
	result = db.ProcessCommand(command)
	if result != "(empty list or set)" {
		t.Errorf("Ran:" + command + ".Expected (empty list or set) but Got result:" + result)
	}

	command = "ZADD k1 50 m50 60 m60"
	db.ProcessCommand(command)

	command = "ZRANGE k1 -10 300 WITHSCORES"
	result = db.ProcessCommand(command)

	if result != "1) 'm1'\n2) 0.1\n3) 'm2'\n4) 0.2\n5) 'm3'\n6) 3\n7) 'm4'\n8) 4\n9) 'm50'\n10) 50\n11) 'm60'\n12) 60\n" {
		t.Errorf("Ran:" + command + ".Expected:\n1) 'm1'\n2) 0.1\n3) 'm2'\n4) 0.2\n5) 'm3'\n6) 3\n7) 'm4'\n8) 4\n,but Got result:\n" + result)
	}

}
