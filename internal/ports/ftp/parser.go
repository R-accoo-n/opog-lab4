package ftp

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/R-accoo-n/opog-lab3/internal"
)

type ProductsService interface {
	BulkCreateProducts(ctx context.Context, params []internal.CreateProductPayload) (int, error)
}

type Parser struct {
	service   ProductsService
	batchSize int
}

func NewParser(service ProductsService, batchSize int) Parser {
	if batchSize <= 0 {
		batchSize = 500
	}
	return Parser{
		service:   service,
		batchSize: batchSize,
	}
}

func (p Parser) Run(ctx context.Context, filePath string) error {
	start := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 4
	reader.ReuseRecord = true

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read csv header: %w", err)
	}
	log.Printf("csv header: %v", header)

	var (
		totalRows    int
		importedRows int
		failedRows   int
		batch        = make([]internal.CreateProductPayload, 0, p.batchSize)
	)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		count, err := p.service.BulkCreateProducts(ctx, batch)
		if err != nil {
			return err
		}
		importedRows += count
		batch = batch[:0]
		return nil
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		totalRows++

		if err != nil {
			failedRows++
			log.Printf("failed to parse row %d: %v", totalRows, err)
			continue
		}

		tax, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err != nil {
			failedRows++
			log.Printf("failed to parse tax in row %d: %v", totalRows, err)
			continue
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(record[3]), 64)
		if err != nil {
			failedRows++
			log.Printf("failed to parse price in row %d: %v", totalRows, err)
			continue
		}

		product := internal.CreateProductPayload{
			Name: strings.TrimSpace(record[0]),
			Category: internal.Category{
				Name: strings.TrimSpace(record[1]),
				Tax:  tax,
			},
			Price: price,
		}

		batch = append(batch, product)

		if len(batch) >= p.batchSize {
			if err = flush(); err != nil {
				return fmt.Errorf("failed to flush batch: %w", err)
			}
		}
	}

	if err = flush(); err != nil {
		return fmt.Errorf("failed to flush final batch: %w", err)
	}

	log.Printf(
		"import finished: total=%d imported=%d failed=%d duration=%s",
		totalRows, importedRows, failedRows, time.Since(start),
	)

	return nil
}
