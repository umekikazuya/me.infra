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
- `meId`（default: `replace-me`）
- `zennUsername`（default: `replace-me`）
- `logLevel`（default: `info`）
- `jwtSecret`（optional: 未指定なら synth / deploy ごとに自動生成）
- `qiitaToken`（default: `replace-me`）
- `appDomain`（optional: 例 `www.example.com`）
- `appCertificateArn`（optional: CloudFront 用 ACM certificate ARN。`appDomain` とセット）
- `apiDomain`（optional: 例 `api.example.com`）
- `apiCertificateArn`（optional: API Gateway 用 ACM certificate ARN。`apiDomain` とセット）

`account` と `region` は `-c` で渡すか、`CDK_DEFAULT_ACCOUNT` / `CDK_DEFAULT_REGION` から読みます。
`appDomain` / `appCertificateArn`、`apiDomain` / `apiCertificateArn` はそれぞれセットで渡します。DNS と証明書はこの repo では作らず、外部管理を前提にしています。

`DataStack` で作るもの:

- DynamoDB table (`PK` / `SK`)
- GSI (`GSI1`, `GSI2`, `GSI3`, `GSI_EMAIL`)
- TTL attribute (`ttl`)

DynamoDB table 名は code/context で固定せず、CloudFormation の自動命名に任せます。実名は `TableNameOutput` から参照します。

`jwtSecret` と `qiitaToken` は CDK context から直接 Lambda 環境変数へ渡します。Secrets Manager は使いません。`jwtSecret` を省略すると synth / deploy のたびに新しい値が生成され、既存 JWT は無効になります。
`meId`、`zennUsername`、`logLevel` も deploy 入力から Lambda 環境変数へ直接渡します。

`ApiStack` で作るもの:

- `me-<env>-api` という名前の Lambda
- API Gateway HTTP API
- optional の custom domain (`apiDomain`, `apiCertificateArn` を渡した場合)
- CloudWatch Logs の log group
- `DYNAMODB_TABLE_NAME`, `JWT_SECRET`, `QIITA_TOKEN`, `ME_ID`, `ZENN_USERNAME`, `LOG_LEVEL` の env 配線
- frontend custom domain がある場合は `https://<appDomain>` からの CORS を許可

`apiDomain` と `apiCertificateArn` を渡すと、API Gateway の default `execute-api` endpoint は無効化され、公開 URL は custom domain に切り替わります。

### CORS の責務分担（API Gateway / アプリ）

このリポジトリでは、CORS は以下の境界で運用します。

- API Gateway (HTTP API) 側で管理するもの
  - preflight (`OPTIONS`) 応答
  - 許可 origin / method / header
  - `maxAge` など preflight キャッシュ設定
- アプリ（Lambda 実装）側で管理するもの
  - 認証・認可（CORS とは別責務）
  - 業務レスポンス本体とエラーレスポンスの内容
  - API Gateway で表現できない個別ヘッダー制御が必要な場合の追加対応

現行の `ApiStack` では `appDomain` がある場合に `CorsPreflight` を有効化し、`https://<appDomain>` のみを `AllowOrigins` に設定します。`AllowHeaders` は `Authorization` / `Content-Type`、`AllowMethods` は `GET,HEAD,OPTIONS,POST,PUT,PATCH,DELETE` を許可します。

運用ルール:

- CORS ヘッダー設定は API Gateway 側を正とし、アプリ側で同じヘッダーを二重付与しない
- `allowCredentials` を使う場合は wildcard origin (`*`) を使わず、明示 origin を列挙する
- frontend domain を変更した場合は `appDomain` を更新して再 deploy する

`WebStack` で作るもの:

- frontend 配信用の private S3 bucket
- S3 向け OAC を含む CloudFront distribution
- CloudFront は `PriceClass 200` に固定
- SPA rewrite 用の CloudFront Function
- optional の frontend custom domain (`appDomain`, `appCertificateArn` を渡した場合)
- frontend bucket 名、distribution ID、distribution domain name、frontend endpoint の outputs

frontend と API は別 origin です。frontend は CloudFront/S3、API は API Gateway custom domain で公開します。

このリポジトリの v1 default は検証環境の後片付けコストを下げるため destroy 寄りです。stack を削除すると DynamoDB、CloudWatch Logs、frontend bucket 内の object も保持せず削除されます。

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

custom domain ありで synth する例:

```sh
cdk synth \
  -c envName=dev \
  -c appDomain=www.example.com \
  -c appCertificateArn=<cloudfront-certificate-arn> \
  -c apiDomain=www.example.com \
  -c apiCertificateArn=<api-certificate-arn>
```

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
  -c appDomain=<frontend-domain> \
  -c appCertificateArn=<cloudfront-certificate-arn> \
  -c apiDomain=<api-domain> \
  -c apiCertificateArn=<api-certificate-arn> \
  -c meId=<me-id> \
  -c zennUsername=<zenn-username> \
  -c qiitaToken=<qiita-token> \
  -c account=<account-id> \
  -c region=<region>
```

## デプロイフロー

### 1. 初回インフラ deploy

以下の順で deploy します。

1. `me-<env>-data`
2. `me-<env>-api`
3. `me-<env>-web`

初回の `ApiStack` では `cmd/placeholder-api/` の placeholder Lambda を使います。これにより、app repo 側が本物の backend code を publish する前に、Lambda・IAM・API Gateway・log group を先に作れます。
custom domain を使う場合は、frontend / API 用の certificate ARN を context で渡します。DNS record 自体は外部で管理します。

例:

```sh
AWS_PROFILE=<profile> cdk deploy \
  me-dev-data \
  me-dev-api \
  me-dev-web \
  -c envName=dev \
  -c appDomain=<frontend-domain> \
  -c appCertificateArn=<cloudfront-certificate-arn> \
  -c apiDomain=<api-domain> \
  -c apiCertificateArn=<api-certificate-arn> \
  -c meId=<me-id> \
  -c zennUsername=<zenn-username> \
  -c qiitaToken=<qiita-token> \
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
- infra repo は Lambda の設定、IAM、API Gateway、custom domain、logs、env wiring を持つ
- app repo は backend artifact の build と `update-function-code` を持つ
- app repo は API base URL として `ApiEndpointOutput` の値を使う

### 4. app repo から frontend を deploy する

frontend の build と配信は app repo 側の責務です。`frontend/dist/` を `WebStack` が作った bucket に sync します。

```sh
aws s3 sync frontend/dist/ s3://<frontend-bucket-name>/ --delete
aws cloudfront create-invalidation --distribution-id <distribution-id> --paths '/*'
```

契約:

- 成果物は `frontend/dist/`
- app repo は `aws s3 sync` と invalidation を持つ
- infra repo は S3 bucket、CloudFront distribution、frontend custom domain の受け口を持つ
- app repo は API base URL を build-time か runtime config で注入する

### 5. app repo 側が必要とする契約値

環境ごとに app repo workflow が必要とする値は以下です。

| 値                         | 例                                                |
| -------------------------- | ------------------------------------------------- |
| Lambda function name       | `me-dev-api`                                      |
| API endpoint               | `ApiEndpointOutput` で export される値            |
| Frontend bucket name       | `FrontendBucketNameOutput` で export される値     |
| CloudFront distribution ID | `FrontendDistributionIdOutput` で export される値 |
| Frontend endpoint          | `FrontendEndpointOutput` で export される値       |

これらは CDK stack が出力するので、app repo workflow に deploy input として渡します。

### 6. 外部 DNS 管理側が必要とする契約値

DNS record はこの repo では作りません。外部の DNS 管理系では以下の出力値を使って record を張ります。

| 値                                          | 用途                                       |
| ------------------------------------------- | ------------------------------------------ |
| `ApiCustomDomainRegionalNameOutput`         | `api.<domain>` の CNAME / Alias target     |
| `ApiCustomDomainRegionalHostedZoneIdOutput` | Route 53 Alias を使う場合の hosted zone ID |
| `FrontendDistributionDomainNameOutput`      | `www.<domain>` の CNAME / Alias target     |

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
- API Gateway / DynamoDB / deploy role は env ごとに分離する

## Frontend デプロイ契約

- app repo が `frontend/dist/` を `aws s3 sync` で配信する
- app repo が `aws cloudfront create-invalidation` を実行する
- infra repo は bucket 名と CloudFront distribution ID を app repo workflow に渡す前提で管理する
- frontend は `www.<domain>`、API は `api.<domain>` の別 origin 構成
