# set env
```
PROJECT_ID=xxxx
REGION=asia-northeast1
```

# Create a Artifact Registry image repository
```
gcloud artifacts repositories create run-otel \
    --repository-format=docker \
    --location=$REGION \
    --project=$PROJECT_ID
```

# build and push image
```
make gcloud-build-and-push
```

# Deploy the code
```
make deploy
```

# check
```
SERVICE_URL=<the url of cloud run service>
curl -H \
    "Authorization: Bearer $(gcloud auth print-identity-token)" \
    $SERVICE_URL
```