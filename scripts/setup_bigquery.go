#!/bin/bash

PROJECT_ID=your-gcp-project-id
DATASET_ID=your-bigquery-dataset-id

bq --project_id $PROJECT_ID mk $DATASET_ID
