package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"halo0201.live/code/golang/event/storage"
)

func initDBConnection() *sql.DB {
	// db, err := sql.Open("mysql", "mysql://ljd:liu.620904@mysql-service.mysql.database.azure.com:3306/go?sslmode=disable")
	db, err := sql.Open("mysql", "ljd:liu.620904@tcp(mysql-service.mysql.database.azure.com:3306)/go?timeout=5s&tls=skip-verify")
	if err != nil {
		log.Fatal("🔥 Error opening database: ", err)
	}
	return db
}

func initRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "ec2-52-54-21-105.compute-1.amazonaws.com:23519",
		Password: "p8a74041154a39e47a9536db38fbb51062875e8edef0d0e57c0209e3e4de84f17",
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	return rdb
}

type Listeners map[string]ListenFunc

type ListenFunc func(string)

type Scheduler struct {
	storage   storage.Storage
	listeners Listeners
}

// 创建一个新调度器
func NewMySQLScheduler(db *sql.DB, lst Listeners) Scheduler {
	return Scheduler{
		storage:   storage.NewMysqlStorage(db), // storage.NewRedisStorage(
		listeners: lst,
	}
}

// 创建一个新调度器
func NewRedisScheduler(rds *redis.Client, lst Listeners) Scheduler {
	return Scheduler{
		storage:   storage.NewRedisStorage(rds), // storage.NewRedisStorage(
		listeners: lst,
	}
}

// 注册监听器/ 处理函数
func (s Scheduler) AddListener(event string, lis ListenFunc) {
	s.listeners[event] = lis
}

// 检查storage中
// func (s Scheduler) checkDataEvent() []storage.Event {
// 	events := []storage.Event{}
// 	rows, err := s.db.Query(fmt.Sprintf(`SELECT id, name, payload FROM jobs WHERE runAt < "%s"`, time.Now().Format("2006-01-02 15:04:05")))
// 	if err != nil {
// 		log.Print("💀 checkDataEvent query error: ", err)
// 		return nil
// 	}
// 	// 遍历查询结果
// 	for rows.Next() {
// 		evt := storage.Event{}
// 		rows.Scan(&evt.ID, &evt.Name, &evt.Payload)
// 		events = append(events, evt)
// 	}
// 	return events
// }

func (s Scheduler) callListener(event storage.Event) {
	if lis, ok := s.listeners[event.Name]; ok {
		go lis(event.Payload)
		err := s.storage.Delete(event)
		if err != nil {
			log.Print("💀 callListener delete error: ", err)
		}
	} else {
		log.Print("💀 callListener listener not found. Name: ", event.Name)
	}

}

// CheckEventsInInterval checks the event in given interval
func (s Scheduler) CheckEventsInInterval(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("⏰ Ticks Received...")
				events := s.storage.CheckDataEvent()
				for _, e := range events {
					s.callListener(e)
				}
			}

		}
	}()
}

func PayBills(payload string) {
	log.Println("💸 Paying bills: ", payload)
}

var eventListeners = Listeners{
	"PayBills": PayBills,
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// db := initDBConnection()
	db := initRedisConnection()
	scheduler := NewRedisScheduler(db, eventListeners)
	scheduler.CheckEventsInInterval(ctx, time.Second*10)

	evt := storage.Event{
		Name:    "PayBills",
		Payload: "Electricity Bill",
		RunAt:   time.Now().Add(1 * time.Minute).Format("2006-01-02 15:04:05"),
	}
	scheduler.storage.Save(evt.Name, evt.Payload, evt.RunAt)

	go func() {
		for range interrupt {
			log.Println("\n❌ Interrupt received closing...")
			cancel()
		}
	}()

	<-ctx.Done()
}
