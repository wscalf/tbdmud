package storage

import (
	"bytes"
	"encoding/json"
	"path/filepath"

	"fmt"
	"log/slog"

	"github.com/wscalf/tbdmud/internal/game"
	bolt "go.etcd.io/bbolt"
)

const (
	PlayerBucket  = "players"
	AccountBucket = "accounts"
)

type BoltDBStore struct {
	db         *bolt.DB
	writeQueue chan writeJob
	nameToId   map[string][]byte
}

func NewBoltDBStore() *BoltDBStore {
	return &BoltDBStore{
		writeQueue: make(chan writeJob, 10),
		nameToId:   map[string][]byte{},
	}
}

func (b *BoltDBStore) Initialize(path string) error {
	path = filepath.Join(path, "world.db")
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}

	b.db = db

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(PlayerBucket))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(AccountBucket))
		return err
	})
	if err != nil {
		return err
	}

	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PlayerBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", PlayerBucket)
		}

		saveData := &game.PlayerSaveData{}
		return bucket.ForEach(func(id, data []byte) error {
			err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&saveData)
			if err != nil {
				slog.Error("error decoding player save data", "err", err, "id", string(id))
				return err
			}

			b.nameToId[saveData.Name] = id
			return nil
		})
	})
}

func (b *BoltDBStore) Process() {
	for j := range b.writeQueue {
		err := j.Run()
		if err != nil {
			slog.Error("error processing write job", "err", err)
		}
	}
}

func (b *BoltDBStore) CreateOrUpdateAccount(account *game.Account) error {
	data := account.GetSaveData()
	b.writeQueue <- writeJob{id: []byte(account.Login), bucketRef: []byte(AccountBucket), data: data, store: b}
	return nil
}

func (b *BoltDBStore) FindAccount(name string) (*game.Account, error) {
	saveData := &game.AccountSaveData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(AccountBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", AccountBucket)
		}

		data := bucket.Get([]byte(name))
		if data != nil {
			return json.NewDecoder(bytes.NewBuffer(data)).Decode(&saveData)
		} else {
			return nil
		}
	})

	if err != nil {
		return nil, err
	}

	if saveData != nil {
		return game.AccountFromSaveData(saveData), nil
	} else {
		return nil, nil
	}
}

func (b *BoltDBStore) CreateOrUpdatePlayer(data *game.PlayerSaveData) error {
	if data == nil {
		return fmt.Errorf("player save data cannot be nil")
	}

	b.nameToId[data.Name] = []byte(data.ID)
	b.writeQueue <- writeJob{id: []byte(data.ID), bucketRef: []byte(PlayerBucket), data: data, store: b}
	return nil
}

func (b *BoltDBStore) FindPlayer(id string) (*game.PlayerSaveData, error) {
	saveData := &game.PlayerSaveData{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PlayerBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", PlayerBucket)
		}

		data := bucket.Get([]byte(id))
		if data != nil {
			return json.NewDecoder(bytes.NewBuffer(data)).Decode(saveData)
		} else {
			saveData = nil
			return nil
		}
	})

	return saveData, err
}

func (b *BoltDBStore) FindPlayerByName(name string) (*game.PlayerSaveData, error) {
	id, ok := b.nameToId[name]
	if !ok {
		return nil, nil
	}

	return b.FindPlayer(string(id))
}

type writeJob struct {
	id        []byte
	bucketRef []byte
	data      any
	store     *BoltDBStore
}

func (j writeJob) Run() error {
	return j.store.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(j.bucketRef)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", j.bucketRef)
		}

		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(j.data)
		if err != nil {
			return err
		}

		return bucket.Put(j.id, buf.Bytes())
	})
}
