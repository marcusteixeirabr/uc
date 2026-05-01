package db

import (
	"fmt"
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

func Connect() (*supabase.Client, error) {
	url := os.Getenv("DATABASE_URL")
	key := os.Getenv("DATABASE_KEY")
	if url == "" || key == "" {
		return nil, fmt.Errorf("DATABASE_URL e DATABASE_KEY são obrigatórias no .env")
	}
	return supabase.NewClient(url, key, nil)
}
