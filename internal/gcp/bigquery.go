package gcp

import (
	"context"

	"github.com/TFMV/db2viz/config"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type BigQueryUploader struct {
	client    *bigquery.Client
	datasetID string
	tableID   string
}

func NewBigQueryUploader(cfg config.BigQueryConfig) *BigQueryUploader {
	client, err := bigquery.NewClient(context.Background(), cfg.ProjectID, option.WithCredentialsFile(cfg.Credentials))
	if err != nil {
		panic(err)
	}
	return &BigQueryUploader{
		client:    client,
		datasetID: cfg.DatasetID,
		tableID:   cfg.TableID,
	}
}

func (u *BigQueryUploader) Upload(ctx context.Context, data []map[string]interface{}) error {
	inserter := u.client.Dataset(u.datasetID).Table(u.tableID).Inserter()

	items := make([]bigquery.ValueSaver, len(data))
	for i, record := range data {
		items[i] = &bigquery.StructSaver{
			Schema:   nil, // You can define the schema here if needed
			Struct:   record,
			InsertID: "", // Optional: an idempotent identifier
		}
	}

	if err := inserter.Put(ctx, items); err != nil {
		return err
	}

	return nil
}
