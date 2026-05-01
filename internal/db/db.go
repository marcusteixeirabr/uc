package db

import (
	"errors"
	"fmt"
	"os"

	supabase "github.com/supabase-community/supabase-go"
)

func Connect() (*supabase.Client, error) {
	url := os.Getenv("DATABASE_URL")
	key := os.Getenv("DATABASE_KEY")
	if url == "" || key == "" {
		return nil, errors.New("DATABASE_URL e DATABASE_KEY são obrigatórias no .env")
	}
	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return nil, fmt.Errorf("db.Connect: criar cliente supabase: %w", err)
	}
	return client, nil
}
