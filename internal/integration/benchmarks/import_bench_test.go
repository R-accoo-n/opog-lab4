package benchmarks

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/R-accoo-n/opog-lab3/internal"
	"github.com/R-accoo-n/opog-lab3/internal/adapters/postgres"
	"github.com/R-accoo-n/opog-lab3/internal/ports/ftp"
)

func BenchmarkImportCSV(b *testing.B) {
	dsn := "postgres://postgres:postgres@127.0.0.1:5432/travellers?sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		b.Fatalf("db connect failed: %v", err)
	}

	client := postgres.NewClient(db)
	service := internal.NewProducts(client)
	parser := ftp.NewParser(service, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = parser.Run(context.Background(), "../data/products_10000.csv"); err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
