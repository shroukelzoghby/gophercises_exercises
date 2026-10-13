package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	taskBucket      = []byte("tasks")
	completedBucket = []byte("completed")
)

// Task is an incomplete TODO item.
type Task struct {
	Key  int
	Value string
}

// CompletedTask is a finished TODO item with a completion timestamp.
type CompletedTask struct {
	Key         int
	Value       string
	CompletedAt time.Time
}

type completedRecord struct {
	Value       string    `json:"value"`
	CompletedAt time.Time `json:"completed_at"`
}

var db *bolt.DB

// Init opens (or creates) the BoltDB file and ensures required buckets exist.
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("opening db: %w", err)
	}
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(taskBucket); err != nil {
			return fmt.Errorf("creating tasks bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(completedBucket); err != nil {
			return fmt.Errorf("creating completed bucket: %w", err)
		}
		return nil
	})
}

// Close closes the database connection.
func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

// CreateTask adds a new incomplete task and returns its key.
func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, err := b.NextSequence()
		if err != nil {
			return err
		}
		id = int(id64)
		return b.Put(itob(id), []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

// AllTasks returns all incomplete tasks.
func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.ForEach(func(k, v []byte) error {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// DeleteTask removes an incomplete task by its database key.
func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

// CompleteTask moves a task from the incomplete bucket to the completed bucket.
func CompleteTask(key int) (string, error) {
	var taskText string
	err := db.Update(func(tx *bolt.Tx) error {
		tasks := tx.Bucket(taskBucket)
		k := itob(key)
		v := tasks.Get(k)
		if v == nil {
			return fmt.Errorf("task %d not found", key)
		}
		taskText = string(v)

		if err := tasks.Delete(k); err != nil {
			return err
		}

		completed := tx.Bucket(completedBucket)
		id64, err := completed.NextSequence()
		if err != nil {
			return err
		}
		record, err := json.Marshal(completedRecord{
			Value:       taskText,
			CompletedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		return completed.Put(itob(int(id64)), record)
	})
	if err != nil {
		return "", err
	}
	return taskText, nil
}

// TasksCompletedToday returns tasks completed on the current calendar day (local time).
func TasksCompletedToday() ([]CompletedTask, error) {
	var result []CompletedTask
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(completedBucket)
		return b.ForEach(func(k, v []byte) error {
			var record completedRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}
			if !record.CompletedAt.Before(startOfDay) && record.CompletedAt.Before(endOfDay) {
				result = append(result, CompletedTask{
					Key:         btoi(k),
					Value:       record.Value,
					CompletedAt: record.CompletedAt,
				})
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
