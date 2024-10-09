package main

import (
	cacophony "cacophony/proto"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	sqlDBHost := os.Getenv("DB_HOST")
	sqlDBPort := os.Getenv("DB_PORT")
	sqlDBUser := os.Getenv("DB_USER")
	sqlDBPassword := os.Getenv("DB_PASSWORD")
	sqlDBName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sqlDBUser, sqlDBPassword, sqlDBHost, sqlDBPort, sqlDBName)
	fmt.Println(dsn)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer sqlDB.Close()

	// Configure connection pooling settings
	sqlDB.SetMaxOpenConns(5)    // Maximum number of open connections
	sqlDB.SetMaxIdleConns(3)    // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(0) // Max connection lifetime (0 means unlimited)
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping to the database: %v", err)
	}

	ctx := context.Background()
	rHost := os.Getenv("REDIS_HOST")
	rPort := os.Getenv("REDIS_PORT")
	rPass := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rHost, rPort),
		Password: rPass, // Redis password here
		DB:       0,     // Use default DB
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cacophony.RegisterChatServiceServer(s, &server{db: sqlDB})
	log.Printf("server listening at %v", lis.Addr())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
