package db

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

type Task struct {
	Id   int
	Name string
}

// TODO: add 'rm' command to remove task from DB without completing
// TODO: add 'completed' command to list completed tasks

func Init() error {
	var err error
	db, err = bolt.Open("tasks.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Fprintf(os.Stderr, "task: error creating database: %v", err)
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
}

func AddTask(task string) (int, error) {
	var taskId int

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return fmt.Errorf("task: get bucket failed")
		}

		id, _ := b.NextSequence()
		taskId = int(id)
		return b.Put(itob(taskId), []byte(task))
	})
	if err != nil {
		return -1, err
	}

	return taskId, nil
}

func ListTasks() ([]Task, error) {
	var taskList []Task
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))

		if err := b.ForEach(func(k, v []byte) error {
			taskList = append(taskList, Task{Id: btoi(k), Name: string(v)})
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return taskList, err
	}
	return taskList, nil
}

func CompleteTask(id int) (string, error) {
	var task string
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return fmt.Errorf("task: get bucket failed")
		}

		task = string(b.Get(itob(id)))
		return b.Delete(itob(id))
	}); err != nil {
		return "", err
	}

	return task, nil
}

// itob returns an 8-byte big endian representation of an int
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi converts a []byte to int
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
