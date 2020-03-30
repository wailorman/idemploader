# idemploader

[OpenAPI DOCS](https://idemploader.wailorman.ru/docs/api/v1/)

## Description

Microservice for "idempotent" file uploading. It means that you can upload file to S3 (the only available provider for now)
without worrying about:

* sharing S3 credentials to API clients (for example, mobile apps)
* storing copies of the same file
* possible overwriting of files with the same name or 
* accidental file removing by clients
* migration to another cloud provider

## Usage

### ENV reference
* **`IDEMPLOADER_API_IMAGE`** — docker image name (default is `wailorman/idemploader_api:latest`)
* **`IDEMPLOADER_S3_HOST`**
* **`IDEMPLOADER_S3_ACCESS_KEY`**
* **`IDEMPLOADER_S3_ACCESS_SECRET`**
* **`IDEMPLOADER_S3_BUCKET`**
* **`IDEMPLOADER_S3_PATH`**
* **`IDEMPLOADER_HOST`** — Host of your's hosted idemploader service (for example: `http://idemploader.example.com/api`)
* **`IDEMPLOADER_ALLOWED_ACCESS_TOKEN`** — random string (key) for authenticating in API (`X-Access-Token` header)

### Local (with docker-compose)

1. ```
   git clone git@github.com:wailorman/idemploader.git
   cd idemploader/
   cp sample.env .env
   ```
2. Fill your credentials above & write them to `.env` file
4. ```
   docker-compose up -d
   ```
   
### Kubernetes

1. Install [envsubst](https://command-not-found.com/envsubst)
2. Load credentials to your ENV (as in reference). If you already filled `.env` file, you can use: `export $(cat k8s.env | xargs)`
3. Ensure kubernetes yaml filled correctly
   ```
   cat .k8s/idemploader-api.yml | envsubst
   ```
4. Apply kubernetes configuration
   ```
   envsubst < .k8s/idemploader-api.yml | kubectl apply -f -
   ```
