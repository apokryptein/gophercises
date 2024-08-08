package db

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB

type Task struct {
	Name string
	Id   int
}

// TODO: add 'completed' command to list completed tasks

func Init(path string) error {
	var err error
	db, err = bolt.Open(path+"tasks.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Fprintf(os.Stderr, "task: error creating database: %v", err)
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte("tasks"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("completed"))
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
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))

		if err := b.ForEach(func(k, v []byte) error {
			taskList = append(taskList, Task{Name: string(v), Id: btoi(k)})
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

func ListCompleted() ([]string, error) {
	date := strings.Split(getDateTime(), " ")[0]
	var completed []string

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("completed"))

		if err := b.ForEach(func(k, v []byte) error {
			tDate := strings.Split(string(k), " ")[0]
			if date == tDate {
				completed = append(completed, string(v))
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return completed, nil
}

func CompleteTask(id int) (string, error) {
	if err := addToCompleted(id); err != nil {
		return "", err
	}
	return RemoveTask(id)
}

func RemoveTask(id int) (string, error) {
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

func addToCompleted(id int) error {
	var compTask string
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		if b == nil {
			return fmt.Errorf("task: get bucket failed")
		}
		compTask = string(b.Get(itob(id)))
		return nil
	}); err != nil {
		return err
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("completed"))
		if b == nil {
			return fmt.Errorf("task: get bucket failed")
		}

		date := getDateTime()
		return b.Put([]byte(date), []byte(compTask))
	})
	if err != nil {
		return err
	}

	return nil
}

func getDateTime() string {
	t := strings.Split(time.Now().String(), " ")
	dateTime := t[0:2]
	return strings.Join(dateTime, " ")
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
