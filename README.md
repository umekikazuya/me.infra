# me.infra

AWS CDK for Go repository for the `me` system infrastructure.

## v1 scope

- `DataStack`, `ApiStack`, and `WebStack` are the initial delivery scope.
- `DomainStack` stays optional until the core stacks are working.
- Frontend and backend artifacts are built outside this repository and handed off to the IaC flow.

## Prerequisites

- Go 1.26.2
- Node.js 24.x
- AWS CLI with a configured profile
- AWS CDK CLI

Install the managed runtimes with mise:

```sh
mise install
```

Install the CDK CLI if it is not already available:

```sh
npm install -g aws-cdk
```

Bootstrap the target AWS environment before the first deploy:

```sh
cdk bootstrap aws://<account-id>/<region>
```

## Configuration

The CDK app reads the following context keys:

- `envName` (default: `dev`)
- `account`
- `region`
- `prefix`
- `tableName` (default: `me.`)
- `meId` (default: `replace-me`)
- `zennUsername` (default: `replace-me`)
- `logLevel` (default: `info`)

`account` and `region` can be passed with `-c` or inherited from `CDK_DEFAULT_ACCOUNT` and `CDK_DEFAULT_REGION`.

`DataStack` creates:

- DynamoDB table with `PK` / `SK`
- GSIs `GSI1`, `GSI2`, `GSI3`, `GSI_EMAIL`
- TTL attribute `ttl`
- Secrets for `jwtSecret` and `qiitaToken`
- SSM parameters for `meId`, `zennUsername`, and `logLevel`

`ApiStack` creates:

- Lambda function named `me-<env>-api`
- Lambda Function URL
- CloudWatch log group with retention
- Environment wiring for `DYNAMODB_TABLE_NAME`, `JWT_SECRET`, `QIITA_TOKEN`, `ME_ID`, `ZENN_USERNAME`, `LOG_LEVEL`

`WebStack` creates:

- Private S3 bucket for frontend assets
- CloudFront distribution with OAC for S3 and Lambda Function URL origins
- `/api/*` behavior to the API Function URL
- CloudFront Function for SPA route rewrite
- Outputs for frontend bucket name, distribution ID, and distribution domain name

## Commands

Run unit tests:

```sh
go test ./...
```

Synthesize the CloudFormation templates:

```sh
cdk synth -c envName=dev
```

The CDK entrypoint packages the placeholder API automatically before synthesis.

Diff against a target environment:

```sh
AWS_PROFILE=<profile> cdk diff -c envName=dev -c account=<account-id> -c region=<region>
```

Deploy a single stack:

```sh
AWS_PROFILE=<profile> cdk deploy me-dev-data -c envName=dev -c account=<account-id> -c region=<region>
```

Deploy the core stacks in order:

```sh
AWS_PROFILE=<profile> cdk deploy \
  me-dev-data \
  me-dev-api \
  me-dev-web \
  -c envName=dev \
  -c account=<account-id> \
  -c region=<region>
```

## Deployment flow

### 1. Initial infrastructure deploy

Deploy in this order:

1. `me-<env>-data`
2. `me-<env>-api`
3. `me-<env>-web`

The first `ApiStack` deploy uses the placeholder Lambda from `cmd/placeholder-api/` so the function, IAM policy, log group, and Function URL can be created before the app repo starts publishing the real backend code.

Example:

```sh
AWS_PROFILE=<profile> cdk deploy \
  me-dev-data \
  me-dev-api \
  me-dev-web \
  -c envName=dev \
  -c account=<account-id> \
  -c region=<region>
```

### 2. Infra-only changes

When only the IaC changes, use the normal CDK flow:

```sh
AWS_PROFILE=<profile> cdk diff me-dev-web -c envName=dev -c account=<account-id> -c region=<region>
AWS_PROFILE=<profile> cdk deploy me-dev-web -c envName=dev -c account=<account-id> -c region=<region>
```

Replace `me-dev-web` with `me-dev-data` or `me-dev-api` as needed.

### 3. Backend code deploy from app repo

The app repo owns backend build and code rollout. After building a Lambda zip whose root contains `bootstrap`, the app repo GitHub Actions workflow updates the function code directly:

```sh
aws lambda update-function-code \
  --function-name me-dev-api \
  --zip-file fileb://<path-to-backend-zip> \
  --publish
```

Deployment contract:

- Function name format: `me-<env>-api`
- Infra repo owns Lambda configuration, IAM, Function URL, logs, and environment wiring
- App repo owns backend artifact build and `update-function-code`

### 4. Frontend deploy from app repo

The app repo owns frontend build and static asset sync. Deploy `frontend/dist/` to the bucket created by `WebStack`:

```sh
aws s3 sync frontend/dist/ s3://<frontend-bucket-name>/ --delete
aws cloudfront create-invalidation --distribution-id <distribution-id> --paths '/*'
```

Deployment contract:

- Artifact shape: `frontend/dist/`
- App repo owns `aws s3 sync` and invalidation
- Infra repo owns the S3 bucket and CloudFront distribution

### 5. Contract values the app repo needs

The app repo workflow needs these values per environment:

| Value                               | Example                                            |
| ----------------------------------- | -------------------------------------------------- |
| Lambda function name                | `me-dev-api`                                       |
| Frontend bucket name                | exported by `FrontendBucketNameOutput`             |
| CloudFront distribution ID          | exported by `FrontendDistributionIdOutput`         |
| CloudFront distribution domain name | exported by `FrontendDistributionDomainNameOutput` |

These values are produced by the CDK stacks and should be surfaced to the app repo workflow as deployment inputs.

## Repository layout

```text
cmd/app/         CDK application entrypoint
cmd/placeholder-api/ initial placeholder Lambda source
lib/config/      shared environment and naming config
lib/stacks/      stack definitions
scripts/         packaging and CDK entrypoint helpers
```

The first bootstrap PR creates placeholder `DataStack`, `ApiStack`, and `WebStack` so the repository can synthesize before AWS resources are added incrementally.

## API deployment contract

- The infra repo creates the Lambda function and its surrounding AWS resources.
- The initial function code comes from `cmd/placeholder-api/`.
- The generated placeholder zip is kept under `.artifacts/` and is not committed.
- After the first deploy, the app repo GitHub Actions pipeline updates code with `aws lambda update-function-code`.

## Frontend deployment contract

- The app repo deploys `frontend/dist/` to the S3 bucket with `aws s3 sync`.
- The app repo invalidates CloudFront after deploy with `aws cloudfront create-invalidation --distribution-id <distribution-id> --paths '/*'`.
- The infra repo provides the bucket name and CloudFront distribution ID as deployment contract values for the app repo workflow.
