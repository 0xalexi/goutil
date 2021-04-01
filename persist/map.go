// Data structures backed by persistent storage for easily retaining data across runs.
package persist

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	util "github.com/alexi/goutil"
	log "github.com/alexi/goutil/log"

	badger "github.com/dgraph-io/badger"
)

func RemoveStore(key string) {
	os.Remove(pathname(key))
}

func pathname(key string) string {
	return fmt.Sprintf(".store/%s", key)
}

type MarshalUnmarshaller interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	Copy() MarshalUnmarshaller
}

type PersistentStringMap struct {
	key                string
	mu                 sync.RWMutex
	objectTypeInstance interface{}
}

// func getReflectValue(value interface{}) (bool, reflect.Value) {
// 	v := reflect.ValueOf(value)
// 	if v.Kind() == reflect.Ptr {
// 		el := v.Elem()
// 		if !el.IsValid() {
// 			return false, v
// 		}
// 		v = reflect.ValueOf(el.Interface())
// 	}
// 	return true, v
// }

func unmarshalStruct(result interface{}) interface{} {
	outValue := reflect.Indirect(reflect.New(reflect.ValueOf(result).Type()))
	return outValue.Addr().Interface()
}

func getReflectValue(value interface{}) (bool, reflect.Value) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		el := v.Elem()
		if !el.IsValid() {
			return false, v
		}
		v = reflect.ValueOf(el.Interface())
	}
	return true, v
}

func (m *PersistentStringMap) unmarshal(b []byte) (interface{}, error) {
	if encoder, ok := m.objectTypeInstance.(MarshalUnmarshaller); ok {
		v := encoder.Copy()
		err := v.Unmarshal(b)
		return v, err
	}
	v := unmarshalStruct(m.objectTypeInstance)
	err := util.DecodeBytes(b, v)
	ok, r := getReflectValue(v)
	if !ok {
		return nil, err
	}
	return r.Interface(), err
}

func NewPersistentStringMap(key string, otype interface{}) *PersistentStringMap {
	m := &PersistentStringMap{
		key:                key,
		objectTypeInstance: otype,
	}
	return m
}

func getdb(key string) (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions(pathname(key)))
	if err != nil {
		if _, err := os.Stat(".store"); os.IsNotExist(err) {
			if err = os.Mkdir(".store", 0777); err != nil {
				log.LogError(err)
				return nil, err
			}
		}
		return badger.Open(badger.DefaultOptions(pathname(key)))
	}
	return db, err
}

func (m *PersistentStringMap) Write(k string, v interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	db, err := getdb(m.key)
	if err != nil {
		log.LogError(err)
		return
	}
	defer db.Close()
	if err = db.Update(func(txn *badger.Txn) error {
		var _v []byte
		if encoder, ok := v.(MarshalUnmarshaller); ok {
			_v, _ = encoder.Marshal()
		} else {
			_v, err = util.GetBytes(v)
			if err != nil {
				log.LogError("get-bytes error:", err)
			}
		}
		return txn.Set([]byte(k), _v)
	}); err != nil {
		log.LogError(err)
	}
	return
}

func (m *PersistentStringMap) Delete(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	db, err := badger.Open(badger.DefaultOptions(pathname(m.key)))
	if err != nil {
		log.LogError(err)
		return
	}
	defer db.Close()
	if err = db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(k))
	}); err != nil {
		log.LogError(err)
	}
	return
}

func (m *PersistentStringMap) Read(k string) (v interface{}) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, err := badger.Open(badger.DefaultOptions(pathname(m.key)))
	if err != nil {
		log.LogError(err)
		return
	}
	defer db.Close()
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			log.LogError(err)
			return err
		}
		if err = item.Value(func(val []byte) error {
			var _err error
			v, _err = m.unmarshal(val)
			return _err
		}); err != nil {
			log.LogError(err)
		}
		return nil
	}); err != nil {
		log.LogError(err)
	}
	return
}

func (m *PersistentStringMap) ReadOk(k string) (v interface{}, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, err := badger.Open(badger.DefaultOptions(pathname(m.key)))
	if err != nil {
		log.LogError(err)
		return
	}
	defer db.Close()
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		// if err == badger.ErrKeyNotFound {
		// 	log.LogWarn(k, "not found")
		// 	return err
		// } else if err != nil {
		// 	log.LogError(err)
		// 	return err
		// }
		if err != nil {
			log.LogError(err)
			return err
		}
		ok = true
		if err = item.Value(func(val []byte) error {
			var _err error
			v, _err = m.unmarshal(val)
			return _err
		}); err != nil {
			log.LogError(err)
		}
		return nil
	}); err != nil {
		log.LogError(err)
	}
	return
}
