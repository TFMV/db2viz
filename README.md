# db2viz

db2viz is a data pipeline project that demonstrates how to move data from an on-premises database (Postgres) to Google Cloud BigQuery for visualization in Power BI.

## Project Structure

- `cmd/main.go`: The entry point for the application.
- `config/config.go`: Configuration loader.
- `internal/db/postgres_connector.go`: Connects to the Postgres database.
- `internal/data/loader.go`: Loads data from Postgres.
- `internal/data/transformer.go`: Transforms data before uploading.
- `internal/gcp/bigquery.go`: Uploads data to Google BigQuery.
- `scripts/run_postgres_docker.sh`: Script to run a Postgres container.

## Getting Started

1. Run the Postgres Docker container

    ```sh
    ./scripts/run_postgres_docker.sh
    ```

2. Configure your `config/config.yaml` with the necessary details for Postgres and Google BigQuery.

3. Build and run the Go application

    ```sh
    docker build -t db2viz .
    docker run --rm db2viz
    ```

## Description

This project demonstrates a simple ETL (Extract, Transform, Load) pipeline:

1. Extract data from a Postgres database.
2. Transform the data as necessary.
3. Load the data into Google BigQuery for visualization in Power BI.
