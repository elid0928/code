package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Event 事件
type Event struct {
	ID      uint
	Name    string
	Payload string
	RunAt   string
}
type Storage interface {
	Save(event string, payload string, runAt string) error
	CheckDataEvent() []Event
	Delete(Event) error
}

// redis 可以使用hash

type MysqlStorage struct {
	// db *sql.DB
	db *sql.DB
}

func NewMysqlStorage(db *sql.DB) MysqlStorage {
	return MysqlStorage{
		db: db,
	}
}

func (m MysqlStorage) Save(event string, payload string, runAt string) error {
	log.Print("🚀 Scheduling event", event, "to run at ", runAt)

	sql := fmt.Sprintf(`INSERT INTO jobs (name, payload, runAt) VALUES ("%s", "%s", "%s")`, event, payload, runAt)
	log.Print("🚀 Schedule SQL: ", sql)
	_, err := m.db.Exec(sql)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
	return err
}
func (m MysqlStorage) CheckDataEvent() []Event {
	events := []Event{}
	rows, err := m.db.Query(fmt.Sprintf(`SELECT id, name, payload FROM jobs WHERE runAt < "%s"`, time.Now().Format("2006-01-02 15:04:05")))
	if err != nil {
		log.Print("💀 checkDataEvent query error: ", err)
		return nil
	}
	// 遍历查询结果
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.Payload)
		events = append(events, evt)
	}
	return events
}
func (m MysqlStorage) Delete(evt Event) error {
	sql := fmt.Sprintf(`DELETE FROM jobs WHERE id = %d`, evt.ID)
	_, err := m.db.Exec(sql)
	return err
}

type RedisStorage struct {
	// redis
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) RedisStorage {
	return RedisStorage{
		client: client,
	}
}
func (r RedisStorage) Save(event string, payload string, runAt string) error {
	log.Print("🚀 Scheduling event", event, "to run at ", runAt)
	// r.client.HSet()
	key := fmt.Sprintf("event:%s:%s", event, runAt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r.client.HSet(ctx, key, "name", event)
	r.client.HSet(ctx, key, "payload", payload)
	r.client.HSet(ctx, key, "runAt", runAt)
	log.Print("🚀 Schedule SQL: payload, runAt ")
	return nil
}
func (r RedisStorage) CheckDataEvent() []Event {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	// 将所有的event:*的key都取出来
	keys := r.client.Keys(ctx, "event:*")
	events := []Event{}
	for _, key := range keys.Val() {
		payload := r.client.HGet(ctx, key, "payload")
		runAt := r.client.HGet(ctx, key, "runAt")
		name := r.client.HGet(ctx, key, "name")
		if runAt.Val() < time.Now().Format("2006-01-02 15:04:05") {
			events = append(events, Event{
				Name:    name.Val(),
				Payload: payload.Val(),
				RunAt:   runAt.Val(),
			})
		}

	}

	return events
}
func (r RedisStorage) Delete(evt Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	key := fmt.Sprintf("event:%s:%s", evt.Name, evt.RunAt)
	log.Printf("🚀 Deleting event %s", key)
	return r.client.Del(ctx, key).Err()
}
