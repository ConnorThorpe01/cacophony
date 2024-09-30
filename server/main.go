package main

import (
	cacophony "cacophony/proto"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	dbHost := os.Getenv("CACOPHONY_DB_HOST")
	dbPort := os.Getenv("CACOPHONY_DB_PORT")
	dbUser := os.Getenv("CACOPHONY_DB_USER")
	dbPassword := os.Getenv("CACOPHONY_DB_PASSWORD")
	dbName := os.Getenv("CACOPHONY_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Configure connection pooling settings
	db.SetMaxOpenConns(5)    // Maximum number of open connections
	db.SetMaxIdleConns(3)    // Maximum number of idle connections
	db.SetConnMaxLifetime(0) // Max connection lifetime (0 means unlimited)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	cacophony.RegisterChatServiceServer(s, &server{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
