# Reaction Bot - アーキテクチャ

## アーキテクチャ図

```mermaid
graph TB
    subgraph "Presentation Layer"
        MAIN[main.go<br/>エントリーポイント<br/>DI・起動・シャットダウン]
    end

    subgraph "Interface Layer"
        HANDLER[DiscordHandler<br/>Discord イベントハンドリング<br/>HandleReactionAdd<br/>HandleReactionRemove]
    end

    subgraph "Use Cases Layer"
        TRANSFER[TransferMessageUseCase<br/>メッセージ転送ロジック<br/>TransferMessage<br/>DeleteTransferredMessage]
    end

    subgraph "Entities Layer"
        CONFIG[Config<br/>設定・ドメインモデル<br/>LoadConfig]
    end

    subgraph "External"
        DISCORD[Discord API<br/>discordgo]
        ENV[環境変数<br/>.env]
    end

    MAIN -->|creates & injects| HANDLER
    MAIN -->|creates & injects| TRANSFER
    MAIN -->|loads| CONFIG

    HANDLER -->|uses| TRANSFER
    HANDLER -->|references| CONFIG

    TRANSFER -->|references| CONFIG

    CONFIG -->|reads| ENV
    HANDLER <-->|API calls| DISCORD
    TRANSFER <-->|API calls| DISCORD

    classDef presentation fill:#ff6b6b,stroke:#c92a2a,color:#fff
    classDef interface fill:#ffd93d,stroke:#f08700,color:#000
    classDef usecase fill:#4ecdc4,stroke:#0a9396,color:#fff
    classDef entity fill:#95e1d3,stroke:#38b000,color:#000
    classDef external fill:#e0e0e0,stroke:#666,color:#000

    class MAIN presentation
    class HANDLER interface
    class TRANSFER usecase
    class CONFIG entity
    class DISCORD,ENV external
```

## レイヤー責務

### Presentation Layer
- **main.go**: アプリケーションのエントリーポイント、依存関係の注入（DI）、Bot起動・シャットダウン

### Interface Layer
- **DiscordHandler**: Discordイベントのハンドリング、外部ライブラリ（discordgo）との接続

### Use Cases Layer
- **TransferMessageUseCase**: メッセージ転送のビジネスロジック、Discord API操作のカプセル化

### Entities Layer
- **Config**: 設定とドメインモデル、環境変数の読み込みとバリデーション

## 依存関係の方向

```
Presentation → Interface → Use Cases → Entities
```

各レイヤーは内側（下位）のレイヤーのみに依存し、外側（上位）のレイヤーに依存しない（クリーンアーキテクチャの原則）