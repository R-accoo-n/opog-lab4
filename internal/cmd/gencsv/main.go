package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	file, err := os.Create("./internal/integration/data/products_10000.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	_ = writer.Write([]string{"name", "category_name", "category_tax", "price"})

	categories := []struct {
		Name string
		Tax  string
	}{
		{"Laptops", "20"},
		{"Smartphones", "15"},
		{"Accessories", "10"},
		{"TV", "18"},
	}

	for i := 1; i <= 10000; i++ {
		c := categories[r.Intn(len(categories))]
		price := 100 + r.Float64()*10000

		row := []string{
			fmt.Sprintf("Product_%d", i),
			c.Name,
			c.Tax,
			fmt.Sprintf("%.2f", price),
		}
		_ = writer.Write(row)
	}

	fmt.Println("CSV file generated: ./internal/integration/data/products_10000.csv")
}
