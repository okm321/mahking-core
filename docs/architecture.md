# Architecture

## 全体構成

```mermaid
graph TB
    User([ユーザー])

    subgraph Cloudflare
        DNS[DNS<br/>mahking-api.okmkm.dev]
    end

    subgraph GCP["GCP (mahking-{env})"]
        subgraph LB["Load Balancer"]
            GlobalIP[Global IP]
            HTTPS_FWD["Forwarding Rule<br/>:443"]
            HTTP_FWD["Forwarding Rule<br/>:80"]
            HTTPS_Proxy["HTTPS Proxy<br/>+ Certificate Manager"]
            HTTP_Proxy["HTTP Proxy"]
            URLMap["URL Map"]
            Redirect["URL Map<br/>301 → HTTPS"]
            Backend["Backend Service"]
            NEG["Serverless NEG"]
        end

        subgraph VPC["VPC (mahking-{env}-vpc)"]
            subgraph Subnet["Subnet (10.0.0.0/24)"]
                CloudRun["Cloud Run<br/>mahking-{env}-api<br/>Go App :8080"]
            end

            subgraph Google["Google Managed VPC"]
                CloudSQL["Cloud SQL<br/>PostgreSQL 18<br/>Private IP: 10.x.x.x"]
            end

            PSC["Private Service<br/>Connection"]
        end

        AR["Artifact Registry<br/>mahking-{env}"]
        SM["Secret Manager<br/>PG_PASS"]
    end

    User -->|HTTPS| DNS
    DNS -->|A Record| GlobalIP
    GlobalIP --> HTTPS_FWD
    GlobalIP --> HTTP_FWD
    HTTPS_FWD --> HTTPS_Proxy
    HTTP_FWD --> HTTP_Proxy
    HTTPS_Proxy --> URLMap
    HTTP_Proxy --> Redirect
    URLMap --> Backend
    Backend --> NEG
    NEG --> CloudRun
    CloudRun -->|Direct VPC Egress<br/>PRIVATE_RANGES_ONLY| CloudSQL
    CloudRun -.->|PG_PASS| SM
    Subnet ---|VPC Peering| PSC
    PSC --- Google
```

## ローカル開発（マイグレーション）

```mermaid
graph LR
    subgraph Local["ローカル PC"]
        Atlas["Atlas<br/>migrate apply/diff"]
        Proxy["Cloud SQL Auth Proxy<br/>localhost:15432"]
        pgcli["pgcli"]
        gcloud["gcloud CLI"]
    end

    subgraph GCP["GCP (mahking-dev)"]
        SM3["Secret Manager<br/>db-password"]
        SQL2["Cloud SQL<br/>Public IP<br/>authorized_networks = ∅"]
    end

    gcloud -->|"パスワード取得"| SM3
    Atlas -->|"--var db_password"| Proxy
    pgcli --> Proxy
    Proxy -->|"IAM 認証 + TLS"| SQL2

    style SQL2 fill:#4285f4,color:#fff
    style Proxy fill:#0f9d58,color:#fff
```

- dev のみ Public IP を有効化（`authorized_networks = []` で直接アクセスは全拒否）
- Auth Proxy が IAM 認証で Cloud SQL に接続
- prd は Private IP のみ（ローカルからのアクセス不可）

## ネットワーク構成

```mermaid
graph LR
    subgraph Internet
        User([ユーザー])
    end

    subgraph VPC["VPC: mahking-{env}-vpc"]
        subgraph Subnet["Subnet: mahking-{env}-subnet<br/>10.0.0.0/24"]
            CR[Cloud Run<br/>Direct VPC Egress]
        end

        subgraph PSC["Private Service Connection<br/>/16 予約済み"]
            SQL[(Cloud SQL<br/>Private IP<br/>10.x.x.x)]
        end
    end

    User -->|HTTPS :443| LB["Load Balancer<br/>Global IP"]
    LB --> CR
    CR -->|Private IP| SQL

    style SQL fill:#4285f4,color:#fff
    style CR fill:#0f9d58,color:#fff
    style LB fill:#f4b400,color:#fff
```

## CI/CD パイプライン

```mermaid
graph LR
    subgraph GitHub
        PR["Pull Request"]
        DevBranch["dev branch push"]
        MainBranch["main branch push"]
    end

    subgraph CI["CI Job"]
        Vet["go vet"]
        Lint["golangci-lint"]
        Test["go test"]
    end

    subgraph CD["Deploy Job (_deploy.yml)"]
        WIF["WIF 認証<br/>OIDC Token"]
        Build["docker build<br/>+ push"]
        Deploy["gcloud run deploy<br/>--image のみ"]
    end

    subgraph GCP
        AR2["Artifact Registry"]
        CR_Dev["Cloud Run<br/>dev"]
        CR_Prd["Cloud Run<br/>prd"]
    end

    PR --> CI
    DevBranch --> CI
    MainBranch --> CI

    CI --> Vet --> Lint --> Test

    DevBranch -.->|needs: ci| CD
    MainBranch -.->|needs: ci<br/>environment: production| CD

    CD --> WIF --> Build --> Deploy
    Build --> AR2
    Deploy -->|dev| CR_Dev
    Deploy -->|prd 手動承認| CR_Prd
```

## Terraform モジュール依存関係

```mermaid
graph TD
    PS["project_services<br/>GCP API 有効化"]

    PS --> VPC["vpc<br/>VPC + Subnet<br/>+ Private Service Connection"]
    PS --> SM["secrets<br/>Secret Manager<br/>(db_password)"]
    PS --> AR["artifact_registry<br/>Docker Registry"]
    PS --> WIF["github_actions_wif<br/>Workload Identity<br/>Federation"]

    VPC --> SQL["cloud_sql<br/>PostgreSQL 18<br/>Private IP"]
    SM --> SQL

    VPC --> CR["cloud_run<br/>Go App<br/>Direct VPC Egress"]
    SQL --> CR
    SM -.->|PG_PASS| CR

    CR --> LB["load_balancer<br/>HTTPS LB<br/>+ Certificate Manager"]
```

## IAM 構成

```mermaid
graph TD
    subgraph SA["Service Accounts"]
        CR_SA["Cloud Run デフォルト SA<br/>{number}-compute@<br/>developer.gserviceaccount.com"]
        GA_SA["github-actions-deploy@<br/>mahking-{env}.iam.gserviceaccount.com"]
    end

    subgraph WIF["Workload Identity Federation"]
        Pool["WIF Pool<br/>github-actions-{env}"]
        Provider["OIDC Provider<br/>token.actions.githubusercontent.com"]
        Condition["attribute_condition<br/>repository == okm321/mahking-core"]
    end

    subgraph Roles["IAM Roles"]
        R1["roles/secretmanager.secretAccessor"]
        R2["roles/run.developer"]
        R3["roles/artifactregistry.writer"]
        R4["roles/iam.serviceAccountUser"]
    end

    subgraph Resources
        SM2["Secret Manager<br/>PG_PASS"]
        CR2["Cloud Run"]
        AR3["Artifact Registry"]
    end

    GHA["GitHub Actions<br/>OIDC Token"] --> Pool
    Pool --> Provider
    Provider --> Condition
    Condition -->|workloadIdentityUser| GA_SA

    CR_SA -->|R1| SM2
    GA_SA -->|R2| CR2
    GA_SA -->|R3| AR3
    GA_SA -->|R4| CR_SA
```

## ディレクトリ構成

```
mahking-core/
├── go/                           # Go アプリケーション
│   ├── cmd/server/               # エントリーポイント
│   ├── internal/                 # 内部パッケージ
│   ├── pkg/                      # 共有パッケージ
│   ├── Dockerfile                # マルチステージビルド
│   ├── Makefile                  # build / deploy / release
│   └── .golangci.yml             # linter 設定
├── .github/workflows/
│   ├── ci.yml                    # PR → lint + test
│   ├── _deploy.yml               # reusable deploy workflow
│   ├── deploy-dev.yml            # dev push → CI + deploy
│   └── deploy-prd.yml            # main push → CI + deploy (手動承認)
├── ops/db/                       # DB マイグレーション
│   ├── atlas.hcl                 # Atlas 設定 (local / dev)
│   ├── schema/                   # スキーマ定義
│   ├── migrations/               # マイグレーションファイル
│   └── Makefile                  # proxy / migrate / db-connect
└── ops/terraform/gcp/
    ├── dev/                      # dev 環境
    ├── prd/                      # prd 環境
    └── modules/
        ├── artifact_registry/    # Docker Registry
        ├── cloud_run/            # Cloud Run
        ├── cloud_sql/            # Cloud SQL (PostgreSQL)
        ├── github_actions_wif/   # Workload Identity Federation
        ├── load_balancer/        # HTTPS Load Balancer
        ├── project_services/     # GCP API 有効化
        ├── secret_manager/       # Secret Manager
        └── vpc/                  # VPC + Subnet + PSC
```

## 環境比較

| 設定 | dev | prd |
|-----|-----|-----|
| **Project** | mahking-dev | mahking-prd |
| **Cloud SQL** | db-f1-micro / 10GB | db-custom-1-3840 / 20GB |
| **Cloud SQL Public IP** | ON（Auth Proxy 用） | OFF |
| **Cloud SQL Backup** | OFF | ON |
| **Cloud SQL 削除保護** | OFF | ON |
| **Cloud Run CPU/Mem** | 1 / 512Mi | 2 / 1Gi |
| **Cloud Run Instances** | 0-10 | 1-100 |
| **CDN** | OFF | ON |
| **Deploy 承認** | 自動 | 手動 (environment: production) |
| **Subnet CIDR** | 10.0.0.0/24 | 10.1.0.0/24 |
| **Domain** | mahking-api.okmkm.dev | TBD |
