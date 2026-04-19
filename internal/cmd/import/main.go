package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/R-accoo-n/opog-lab3/internal"
	"github.com/R-accoo-n/opog-lab3/internal/adapters/postgres"
	"github.com/R-accoo-n/opog-lab3/internal/ports/ftp"
)

func main() {
	filePath := flag.String("file", "./internal/integration/data/products_10000.csv", "path to csv file")
	flag.Parse()

	dsn := "postgres://postgres:postgres@127.0.0.1:5432/travellers?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	productsClient := postgres.NewClient(db)
	productsService := internal.NewProducts(productsClient)
	parser := ftp.NewParser(productsService, 500)

	if err = parser.Run(context.Background(), *filePath); err != nil {
		log.Fatalf("import failed: %v", err)
	}

	fmt.Println("CSV import completed successfully")
}
