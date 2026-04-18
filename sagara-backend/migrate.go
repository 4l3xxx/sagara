package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres.zdjoycjhelpizeutevhp:indomieseleraku@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Sedang menambahkan kolom 'status' dan 'admin_notes'...")

	queries := []string{
		"ALTER TABLE consultation_requests ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'New';",
		"ALTER TABLE consultation_requests ADD COLUMN IF NOT EXISTS admin_notes TEXT DEFAULT '';",
	}

	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("Gagal menjalankan query: %v\n", err)
		} else {
			fmt.Println("Berhasil!")
		}
	}

	fmt.Println("\nMigrasi Selesai! Sekarang Anda bisa restart server Go-nya.")
}
