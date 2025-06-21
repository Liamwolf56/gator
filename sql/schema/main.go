package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/google/uuid"
    _ "github.com/lib/pq"

    "gator/config"
    "gator/internal/database"
)

type state struct {
    db  *database.Queries
    cfg *config.Config
}

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    db, err := sql.Open("postgres", cfg.DBUrl)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }

    defer db.Close()

    st := &state{
        db:  database.New(db),
        cfg: cfg,
    }

    if len(os.Args) < 2 {
        fmt.Println("expected 'register' or 'login'")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "register":
        st.handleRegister()
    case "login":
        st.handleLogin()
    default:
        fmt.Println("unknown command")
        os.Exit(1)
    }
}

func (s *state) handleRegister() {
    if len(os.Args) < 3 {
        fmt.Println("name required")
        os.Exit(1)
    }

    name := os.Args[2]
    ctx := context.Background()
    now := time.Now()

    user, err := s.db.CreateUser(ctx, database.CreateUserParams{
        ID:        uuid.New(),
        CreatedAt: now,
        UpdatedAt: now,
        Name:      name,
    })

    if err != nil {
        fmt.Printf("failed to create user: %v\n", err)
        os.Exit(1)
    }

    s.cfg.CurrentUserName = name
    if err := config.SaveConfig(s.cfg); err != nil {
        fmt.Printf("failed to save config: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("✅ user registered:", user.Name)
}

func (s *state) handleLogin() {
    if len(os.Args) < 3 {
        fmt.Println("name required")
        os.Exit(1)
    }

    name := os.Args[2]
    ctx := context.Background()

    user, err := s.db.GetUser(ctx, name)
    if err != nil {
        fmt.Println("❌ user not found")
        os.Exit(1)
    }

    s.cfg.CurrentUserName = name
    if err := config.SaveConfig(s.cfg); err != nil {
        fmt.Printf("failed to save config: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("✅ logged in as:", user.Name)
}

