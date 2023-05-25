package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

// Interface for the DB Client
type IDBClient interface {
	Initialize(filepath string)
	Open()
	Close()
	CreateUser(u User) (bool, error)
	Login(email, password string) (*TokenPair, error)
	CheckUserIsNew(email string) (bool, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id string) (User, error)
}

// Struct to handle the DB Connection
type DBClient struct {
	client   *bolt.DB
	filepath string
}

// The DB Client for the application
var DB IDBClient

const (
	usersBucketName string = "Users"
)

func (db *DBClient) Initialize(filepath string) {
	db.filepath = filepath
	db.Open()
	defer db.Close()
	err := db.client.Update(func(txn *bolt.Tx) error {
		// Initialize Users Bucket
		_, err := txn.CreateBucketIfNotExists([]byte(usersBucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (db *DBClient) Open() {
	if db.filepath == "" {
		log.Fatal("Filepath required for Database")
	}
	d, err := bolt.Open(db.filepath, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	db.client = d
}

func (db *DBClient) Close() {
	db.client.Close()
}

func (db *DBClient) CreateUser(u User) (bool, error) {
	isNew, err := db.CheckUserIsNew(u.Email)
	if err != nil || !isNew {
		return false, err
	}
	db.Open()
	defer db.Close()
	err = db.client.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(usersBucketName))
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		u.ID = strconv.Itoa(int(id))
		u.HashPassword()
		userBytes, err := json.Marshal(u)
		if err != nil {
			return err
		}
		err = b.Put([]byte(u.Email), userBytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (db *DBClient) CheckUserIsNew(email string) (bool, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	if user.ID != "" {
		return false, fmt.Errorf("User Exists with email %s", email)
	}
	return true, nil
}

func (db *DBClient) Login(email, password string) (*TokenPair, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	err = user.CheckPassword(password)
	if err != nil {
		return nil, err
	}
	token, err := TS.GenerateTokenPairForUser(user)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (db *DBClient) GetUserByEmail(email string) (User, error) {
	db.Open()
	defer db.Close()
	user := User{}
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		userBytes := b.Get([]byte(email))
		if userBytes == nil {
			return nil
		}
		err := json.Unmarshal(userBytes, &user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DBClient) GetUserById(id string) (User, error) {
	db.Open()
	defer db.Close()
	user := User{}
	err := db.client.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(usersBucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}
			if user.ID == id {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return User{}, err
	}
	return user, nil
}
