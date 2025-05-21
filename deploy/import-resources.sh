#!/bin/bash

terraform import google_cloud_run_service.disappr locations/us-central1/namespaces/disappr-io/services/disappr

echo "Cloud Build repository already exists in state"

terraform import 'google_firestore_database.default' 'projects/disappr-io/databases/(default)' || echo "Firestore database doesn't exist yet"