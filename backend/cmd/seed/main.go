// Seeds one user per role for local development. Bypasses the invite flow
// since invites require an existing super_admin to send them.
package main

import (
	"context"
	"fmt"
	"log"

	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/config"
	"team1blog/backend/internal/db"
	"team1blog/backend/internal/users"
)

type seedUser struct {
	name     string
	email    string
	password string
	role     users.Role
	wallet   string
}

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	usersRepo := users.NewRepository(pool)

	seedUsers := []seedUser{
		{"Ada Super Admin", "admin@team1.blog", "password123", users.RoleSuperAdmin, ""},
		{"Mira Moderator", "moderator@team1.blog", "password123", users.RoleModerator, ""},
		{"Diego Designer", "designer@team1.blog", "password123", users.RoleGraphicDesigner, ""},
		{"Pat Publisher", "publisher@team1.blog", "password123", users.RolePublisher, ""},
		{"Chidi Contributor", "contributor@team1.blog", "password123", users.RoleContributor, "0x1234567890123456789012345678901234567890"},
	}

	for _, su := range seedUsers {
		if existing, err := usersRepo.GetByEmail(ctx, su.email); err == nil && existing != nil {
			fmt.Printf("skip (exists): %s\n", su.email)
			continue
		}
		hash, err := auth.HashPassword(su.password)
		if err != nil {
			log.Fatalf("hash password for %s: %v", su.email, err)
		}
		u, err := usersRepo.Create(ctx, su.name, su.email, hash, su.role, nil)
		if err != nil {
			log.Fatalf("create user %s: %v", su.email, err)
		}
		if su.wallet != "" {
			if err := usersRepo.UpdateWallet(ctx, u.ID, su.wallet); err != nil {
				log.Fatalf("set wallet for %s: %v", su.email, err)
			}
		}
		fmt.Printf("created: %-12s %-25s password=%s\n", su.role, su.email, su.password)
	}

	fmt.Println("seed complete")
}
