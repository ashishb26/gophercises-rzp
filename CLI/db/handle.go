package db

import (
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

var dbName = "my.db"

/*
	bucketName bucket stores id,task_name of pending tasks
*/
var bucketName = "task_list"

/*
	compBucket bucket stores id,task_name of completed tasks
	here task has the same id as it had in bucketName
*/
var compBucket = "comp_list"

/*
	logBucket bucket stores id,completion_time of completed tasks
	here the id is the same as that of compBucket
*/
var logBucket = "log_list"

/*
	Task struct is used to return the (id,task) values
*/
type Task struct {
	Key   int
	Value string
}

/*
	getPath() returns the path of the home directory
*/
func getPath() string {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, dbName)
	return dbPath
}

/*
	OpenDb() opens/creates the db in the home directory.
	It also creates necessaary buckets if they donot exist
	and returns a *bolt.DB for other functions to use
*/
func OpenDb() *bolt.DB {
	filePath := getPath()
	db, err := bolt.Open(filePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	check(err)
	check(db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			log.Fatalln(err)
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(compBucket))
		if err != nil {
			log.Fatalln(err)
		}
		return nil
	}))

	check(db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(logBucket))
		if err != nil {
			log.Fatalln(err)
		}
		return nil
	}))

	return db
}

/*
	AddTask adds a task to bucketName bucket
*/
func AddTask(taskName string) error {
	db := OpenDb()
	defer db.Close()
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		key, _ := b.NextSequence()

		return b.Put(itob(int(key)), []byte(taskName))
	})
	return err
}

/*
	ListTasks() lists all the pending tasks
*/
func ListTasks() ([]Task, error) {
	db := OpenDb()
	defer db.Close()
	var taskList []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()
		var temp Task
		for k, v := c.First(); k != nil; k, v = c.Next() {
			temp.Key = btoi(k)
			temp.Value = string(v)
			taskList = append(taskList, temp)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return taskList, nil
}

/*
	DeleteTask() removes a pending task before completion from
	the bucketName bucket
*/
func DeleteTask(key int) error {
	db := OpenDb()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		return b.Delete(itob(key))
	})
	return err
}

/*
	CompleteTask() completes a task by removing it from the bucketName
	bucket and making corresponding entries in compBucket and logBucket
*/
func CompleteTask(key int) error {
	db := OpenDb()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		cBucket := tx.Bucket([]byte(compBucket))
		lBucket := tx.Bucket([]byte(logBucket))

		id := itob(key)
		task := b.Get(id)
		err := b.Delete(id)
		check(err)

		err = cBucket.Put(id, task)
		check(err)

		curTime := time.Now().Format("2006-01-02 15:04:05")
		err = lBucket.Put(id, []byte(curTime))
		check(err)

		return nil
	})

	return err
}

/*
	GetCompletedTasks() returns all tasks that were completed in the last
	"dur" hours
*/
func GetCompletedTasks(dur time.Duration) error {
	db := OpenDb()
	return db.Update(func(tx *bolt.Tx) error {
		cBucket := tx.Bucket([]byte(compBucket))
		lBucket := tx.Bucket([]byte(logBucket))
		c1 := cBucket.Cursor()
		var cnt = 0
		for key, task := c1.First(); key != nil; key, task = c1.Next() {
			byteTime := lBucket.Get(key)
			curstrTime := time.Now().Format("2006-01-02 15:04:05")
			strTime := string(byteTime)
			compTime, _ := time.Parse("2006-01-02 15:04:05", strTime)
			curTime, _ := time.Parse("2006-01-02 15:04:05", curstrTime)
			diff := time.Duration(curTime.Sub(compTime).Hours()) * time.Hour

			if diff < dur {
				cnt++
				fmt.Printf("%d. %s\n", cnt, string(task))
			}
		}
		if cnt == 0 {
			fmt.Println("You haven't completed any tasks in the last", dur, " hours")
		}
		return nil
	})
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
