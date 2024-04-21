package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/options"
)

func loginMongo(ctx context.Context, addr string, port int, user, password string) error {

	// dsn := "mongodb://dbmon:liujiadong@172.17.0.5:27017/admin?Timeout=4"
	// dinfo, err := mgo.ParseURL(dsn)

	// dinfo.Timeout = time.Second * 4
	// fmt.Printf("%v", dinfo)
	// adr := fmt.Sprintf("%s:%d", addr, port)
	// fmt.Println(adr)
	// go get go.mongodb.org/mongo-driver

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://dbmon:liujiadong@172.17.0.5:27017/admin"))
	if err != nil {
		return err
	}
	s := conn.Database("admin").RunCommand(ctx, map[string]string{"ping": "1"})
	if s.Err() != nil {
		return s.Err()
	}
	d, err := s.DecodeBytes()
	if err != nil {
		return err
	}
	fmt.Println(d.String())

	err = conn.Ping(ctx, nil)
	conn.Disconnect(ctx)
	return err
}

// sudo apt install mongodb-clients
func main() {

	user := "dbmon"
	password := "liujiadong"
	host := "172.17.0.5"
	port := 27017
	ctx, _ := context.WithTimeout(context.Background(), time.Second*200)
	err := loginMongo(ctx, host, port, user, password)
	if err != nil {
		fmt.Println("这不就是错误吗?", err)
	}
}

// mongodb://<username>:<password>@<host>:<port>/<database>
