package helpers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

func InItClient () (*supabase.Client, error) {
	err := godotenv.Load()

	if err != nil {return nil,err}

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_ANON_KEY")

	client := supabase.CreateClient(url,key)

	return client,nil
}