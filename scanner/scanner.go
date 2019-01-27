package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"net"
	"os"
	"strings"
	"time"
)

var c = 65535

func scan(ip string, port int, timeout time.Duration, ports chan int) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			scan(ip, port, timeout, ports)
		}
		c--
		return
	}
	conn.Close()
	ports <- port
	c--
	if c == 1 {
		close(ports)
	}
}

func dbWrite(ip []byte, ports []int, db *bolt.DB) {
	pJson, err := json.Marshal(ports)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("portScanner"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(ip), pJson)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func dbRead(ip []byte, db *bolt.DB) []int {
	var ports []int
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("portScanner"))
		if bucket == nil {
			return fmt.Errorf("Bucket portScanner not found!")
		}
		val := bucket.Get(ip)
		err := json.Unmarshal(val, &ports)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ports
}

func main() {
	ip := os.Args[1]
	dbFound := true
	db, err := bolt.Open("scanner.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		dbFound = false
	}
	defer db.Close()
	ports := make(chan int)
	for port := 1; port <= 65535; port++ {
		go scan(ip, port, 100*time.Millisecond, ports)
	}
	var p []int
	for port := range ports {
		fmt.Println(port)
		p = append(p, port)
	}
	if dbFound == true {
		oldPorts := dbRead([]byte(ip), db)
		for _, port := range p {
			found := false
			for _, oldPort := range oldPorts {
				if port == oldPort {
					found = true
					break
				}
			}
			if found == false {
				fmt.Println(port, "open")
			}
		}
	} else {
		dbWrite([]byte(ip), p, db)
		for _, port := range p {
			fmt.Println(port, "open")
		}
	}
	fmt.Println(p)
}
