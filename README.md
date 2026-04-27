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
