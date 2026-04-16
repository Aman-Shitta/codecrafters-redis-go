package server

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// RDB_OP_CODES
const (
	// End of the RDB file
	EOF = 0xFF
	// Database Selector
	SELECTDB = 0xFE
	// Expire time in seconds, see Key Expiry Timestamp
	EXPIRETIME = 0xFD
	// Expire time in milliseconds, see Key Expiry Timestamp
	EXPIRETIMEMS = 0xFC
	// Hash table sizes for the main keyspace and expires, see Resizedb information
	RESIZEDB = 0xFB
	// Auxiliary fields. Arbitrary key-value settings, see Auxiliary fields
	AUX = 0xFA
)

type Config struct {
	Dir        string
	Dbfilename string
}

func NewConfig() *Config {
	return &Config{
		Dir:        "",
		Dbfilename: "",
	}
}

func (c *Config) UpdateConfig(dir string, dbfilename string) error {

	c.Dir = dir
	c.Dbfilename = dbfilename

	return nil
}

type Item struct {
	Data interface{}
	Type string
}

type Store struct {
	sync.Mutex
	Data    map[string]Item
	Channel map[net.Conn][]string
}

// hold all data for the current session
var SessionStore = &Store{
	Data:    make(map[string]Item),
	Channel: make(map[net.Conn][]string),
}

var mu sync.Mutex

func Get(key string) (Item, bool) {
	SessionStore.Lock()
	defer SessionStore.Unlock()
	item, ok := SessionStore.Data[key]
	return item, ok
}

func Set(key string, val Item) {
	SessionStore.Lock()
	SessionStore.Data[key] = val
	SessionStore.Unlock()
}

func RemoveKey(key string) {
	SessionStore.Lock()
	delete(SessionStore.Data, key)
	SessionStore.Unlock()
}

// Loads the data from givenn RDB file
// adding persistence
func (c *Config) AutoLoad() error {

	dump_file := fmt.Sprintf("%s/%s", c.Dir, c.Dbfilename)

	data, err := os.ReadFile(dump_file)

	if err != nil {
		return err
	}

	ParseRDB(data)

	return nil
}
