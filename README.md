# me.infra

`me` システムをAWS CDK for Go で管理。

## v1 スコープ

- 初期スコープは `DataStack`、`ApiStack`、`WebStack`
- `DomainStack` は本体が固まってから追加する optional 扱い
- frontend / backend の成果物はこのリポジトリでは build せず、別リポジトリ側で作成してデプロイする

## 前提

- Go 1.26.2
- Node.js 24.x
- AWS CLI
- AWS CDK CLI

`mise` でランタイムを入れる:

```sh
mise install
```

CDK CLI が未導入なら入れる:

```sh
npm install -g aws-cdk
```

初回デプロイ前に対象環境を bootstrap する:

```sh
cdk bootstrap aws://<account-id>/<region>
```

## 設定

CDK アプリは以下の context key を読みます。

- `envName`（default: `dev`）
- `account`
- `region`
- `prefix`
- `tableName`（default: `me.`）
- `meId`（default: `replace-me`）
- `zennUsername`（default: `replace-me`）
- `logLevel`（default: `info`）

`account` と `region` は `-c` で渡すか、`CDK_DEFAULT_ACCOUNT` / `CDK_DEFAULT_REGION` から読みます。

`DataStack` で作るもの:

- DynamoDB table (`PK` / `SK`)
- GSI (`GSI1`, `GSI2`, `GSI3`, `GSI_EMAIL`)
- TTL attribute (`ttl`)
- Secrets Manager (`jwtSecret`, `qiitaToken`)
- SSM Parameter Store (`meId`, `zennUsername`, `logLevel`)

`ApiStack` で作るもの:

- `me-<env>-api` という名前の Lambda
- Lambda Function URL
- CloudWatch Logs の log group
- `DYNAMODB_TABLE_NAME`, `JWT_SECRET`, `QIITA_TOKEN`, `ME_ID`, `ZENN_USERNAME`, `LOG_LEVEL` の env 配線

`WebStack` で作るもの:

- frontend 配信用の private S3 bucket
- S3 / Lambda Function URL 向け OAC を含む CloudFront distribution
- `/api/*` を API Function URL へ向ける behavior
- SPA rewrite 用の CloudFront Function
- frontend bucket 名、distribution ID、distribution domain name の outputs

## 基本コマンド

テスト:

```sh
go test ./...
```

template synth:

```sh
cdk synth -c envName=dev
```

placeholder API は synth 前に自動 package されます。

diff:

```sh
AWS_PROFILE=<profile> cdk diff -c envName=dev -c account=<account-id> -c region=<region>
```

単一 stack の deploy:

```sh
AWS_PROFILE=<profile> cdk deploy me-dev-data -c envName=dev -c account=<account-id> -c region=<region>
```

core stack をまとめて deploy:

```sh
AWS_PROFILE=<profile> cdk deploy \
  me-dev-data \
  me-dev-api \
  me-dev-web \
  -c envName=dev \
  -c account=<account-id> \
  -c region=<region>
```

## デプロイフロー

### 1. 初回インフラ deploy

以下の順で deploy します。

1. `me-<env>-data`
2. `me-<env>-api`
3. `me-<env>-web`

初回の `ApiStack` では `cmd/placeholder-api/` の placeholder Lambda を使います。これにより、app repo 側が本物の backend code を publish する前に、Lambda・IAM・Function URL・log group を先に作れます。

例:

```sh
AWS_PROFILE=<profile> cdk deploy \
  me-dev-data \
  me-dev-api \
  me-dev-web \
  -c envName=dev \
  -c account=<account-id> \
  -c region=<region>
```

### 2. IaC だけ変更した場合

IaC だけ変更したときは通常の CDK フローを使います。

```sh
AWS_PROFILE=<profile> cdk diff me-dev-web -c envName=dev -c account=<account-id> -c region=<region>
AWS_PROFILE=<profile> cdk deploy me-dev-web -c envName=dev -c account=<account-id> -c region=<region>
```

必要に応じて `me-dev-web` を `me-dev-data` や `me-dev-api` に読み替えます。

### 3. app repo から backend code を deploy する

backend の build と code rollout は app repo 側の責務です。zip の直下に `bootstrap` が入る Lambda zip を作ったあと、GitHub Actions などから直接 `update-function-code` します。

```sh
aws lambda update-function-code \
  --function-name me-dev-api \
  --zip-file fileb://<path-to-backend-zip> \
  --publish
```

契約:

- function name 形式は `me-<env>-api`
- infra repo は Lambda の設定、IAM、Function URL、logs、env wiring を持つ
- app repo は backend artifact の build と `update-function-code` を持つ

### 4. app repo から frontend を deploy する

frontend の build と配信は app repo 側の責務です。`frontend/dist/` を `WebStack` が作った bucket に sync します。

```sh
aws s3 sync frontend/dist/ s3://<frontend-bucket-name>/ --delete
aws cloudfront create-invalidation --distribution-id <distribution-id> --paths '/*'
```

契約:

- 成果物は `frontend/dist/`
- app repo は `aws s3 sync` と invalidation を持つ
- infra repo は S3 bucket と CloudFront distribution を持つ

### 5. app repo 側が必要とする契約値

環境ごとに app repo workflow が必要とする値は以下です。

| 値                                  | 例                                                        |
| ----------------------------------- | --------------------------------------------------------- |
| Lambda function name                | `me-dev-api`                                              |
| Frontend bucket name                | `FrontendBucketNameOutput` で export される値             |
| CloudFront distribution ID          | `FrontendDistributionIdOutput` で export される値         |
| CloudFront distribution domain name | `FrontendDistributionDomainNameOutput` で export される値 |

これらは CDK stack が出力するので、app repo workflow に deploy input として渡します。

## リポジトリ構成

```text
cmd/app/              CDK application entrypoint
cmd/placeholder-api/  初回 deploy 用 placeholder Lambda
lib/config/           共通設定
lib/stacks/           stack 定義
scripts/              package / CDK 実行補助
```

bootstrap 段階では placeholder を含む `DataStack`、`ApiStack`、`WebStack` を synth / deploy できる状態までをこのリポジトリで持ちます。

## API デプロイ契約

- Lambda 本体と周辺 AWS リソースは infra repo が作る
- 初回コードは `cmd/placeholder-api/` から作る
- 生成された placeholder zip は `.artifacts/` 配下に置き、commit しない
- 初回 deploy 後は app repo の GitHub Actions が `aws lambda update-function-code` でコード更新する

## Frontend デプロイ契約

- app repo が `frontend/dist/` を `aws s3 sync` で配信する
- app repo が `aws cloudfront create-invalidation` を実行する
- infra repo は bucket 名と CloudFront distribution ID を app repo workflow に渡す前提で管理する
