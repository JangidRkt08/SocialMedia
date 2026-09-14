package main

import (
	"fmt"
	"log"

	"github.com/JangidRkt08/SocialMedia/internal/db"
	"github.com/JangidRkt08/SocialMedia/internal/env"
	"github.com/JangidRkt08/SocialMedia/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "")
	fmt.Println(addr)
	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
