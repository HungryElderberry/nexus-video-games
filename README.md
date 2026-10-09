# 🎮 Nexus Video Games API

A high-performance backend API built with Go (Golang), implementing Clean Architecture principles. It features a video games catalog, an in-game wallet, Xendit payment gateway integration, Resend email alerts, and single-game session enforcement.

## 🛠️ Tech Stack

- **Language:** Go 1.22+
- **Framework:** Gin Web Framework
- **Database & ORM:** PostgreSQL (via Supabase) & GORM
- **Validation:** Go Playground Validator v10
- **Authentication:** JWT (JSON Web Tokens)
- **Configuration:** Viper & Godotenv
- **Documentation:** Swaggo (Swagger UI)
- **Third-Party Integrations:** Resend (Emails), Xendit (Payments), IsThereAnyDeal (Game Catalog Sync)

---

## 🚀 1. Project Initialization & Dependencies

To replicate this project from scratch, follow these exact initialization steps.

### Initialize Go Module

```bash
go mod init nexus-video-games
```

### Create Project Directory Structure (`mkdir` & `touch`)

```bash
# 1. Create all directories
mkdir -p cmd/api db/migrations docs internal/config internal/delivery/http/controller internal/delivery/http/middleware internal/delivery/http/route internal/entity internal/gateway/email internal/gateway/game internal/gateway/payment internal/model internal/pkg/jwt internal/repository/postgresql internal/usecase

# 2. Create entry point and root files
touch cmd/api/main.go .env .env.example README.md

# 3. Create database migrations
touch db/migrations/2026100901_create_table_users.up.sql db/migrations/2026100901_create_table_users.down.sql
touch db/migrations/2026100902_create_table_user_balances.up.sql db/migrations/2026100902_create_table_user_balances.down.sql
touch db/migrations/2026100903_create_table_xendit_topups.up.sql db/migrations/2026100903_create_table_xendit_topups.down.sql
touch db/migrations/2026100904_create_table_games_catalog.up.sql db/migrations/2026100904_create_table_games_catalog.down.sql
touch db/migrations/2026100905_create_table_games_owned.up.sql db/migrations/2026100905_create_table_games_owned.down.sql
touch db/migrations/2026100906_create_table_active_sessions.up.sql db/migrations/2026100906_create_table_active_sessions.down.sql

# 4. Create internal configurations
touch internal/config/app.go internal/config/gin.go internal/config/gorm.go internal/config/validator.go internal/config/viper.go

# 5. Create HTTP delivery layer (controllers, middleware, routes)
touch internal/delivery/http/controller/balance_controller.go internal/delivery/http/controller/dashboard_controller.go internal/delivery/http/controller/game_controller.go internal/delivery/http/controller/session_controller.go internal/delivery/http/controller/user_controller.go
touch internal/delivery/http/middleware/auth_middleware.go
touch internal/delivery/http/route/route.go

# 6. Create database entities
touch internal/entity/active_session_entity.go internal/entity/games_catalog_entity.go internal/entity/games_owned_entity.go internal/entity/user_balance_entity.go internal/entity/user_entity.go internal/entity/xendit_topup_entity.go

# 7. Create third-party gateways
touch internal/gateway/email/resend_gateway.go internal/gateway/game/itad_gateway.go internal/gateway/payment/xendit_gateway.go

# 8. Create data models
touch internal/model/balance_model.go internal/model/dashboard_model.go internal/model/game_model.go internal/model/model.go internal/model/session_model.go internal/model/user_model.go

# 9. Create security packages
touch internal/pkg/jwt/jwt.go

# 10. Create PostgreSQL repositories
touch internal/repository/postgresql/balance_repository.go internal/repository/postgresql/game_repository.go internal/repository/postgresql/session_repository.go internal/repository/postgresql/topup_repository.go internal/repository/postgresql/user_repository.go

# 11. Create business logic usecases
touch internal/usecase/balance_usecase.go internal/usecase/dashboard_usecase.go internal/usecase/game_usecase.go internal/usecase/session_usecase.go internal/usecase/user_usecase.go
```

### Install Dependencies

```bash
# Core framework & database
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Utilities & Config
go get github.com/google/uuid
go get github.com/spf13/viper
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10

# Security & Crypto
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5

# Third-Party API SDKs
go get github.com/resend/resend-go/v4

# Swagger Documentation
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
go install github.com/swaggo/swag/cmd/swag@latest

# Tidy module
go mod tidy
```

---

## 📂 2. Project Structure

```plaintext
├── cmd
│   └── api
│       └── main.go
├── db
│   └── migrations
│       ├── 2026100901_create_table_users.down.sql
│       ├── 2026100901_create_table_users.up.sql
│       ├── 2026100902_create_table_user_balances.down.sql
│       ├── 2026100902_create_table_user_balances.up.sql
│       ├── 2026100903_create_table_xendit_topups.down.sql
│       ├── 2026100903_create_table_xendit_topups.up.sql
│       ├── 2026100904_create_table_games_catalog.down.sql
│       ├── 2026100904_create_table_games_catalog.up.sql
│       ├── 2026100905_create_table_games_owned.down.sql
│       ├── 2026100905_create_table_games_owned.up.sql
│       ├── 2026100906_create_table_active_sessions.down.sql
│       └── 2026100906_create_table_active_sessions.up.sql
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── config
│   │   ├── app.go
│   │   ├── gin.go
│   │   ├── gorm.go
│   │   ├── validator.go
│   │   └── viper.go
│   ├── delivery
│   │   └── http
│   │       ├── controller
│   │       │   ├── balance_controller.go
│   │       │   ├── dashboard_controller.go
│   │       │   ├── game_controller.go
│   │       │   ├── session_controller.go
│   │       │   └── user_controller.go
│   │       ├── middleware
│   │       │   └── auth_middleware.go
│   │       └── route
│   │           └── route.go
│   ├── entity
│   │   ├── active_session_entity.go
│   │   ├── games_catalog_entity.go
│   │   ├── games_owned_entity.go
│   │   ├── user_balance_entity.go
│   │   ├── user_entity.go
│   │   └── xendit_topup_entity.go
│   ├── gateway
│   │   ├── email
│   │   │   └── resend_gateway.go
│   │   ├── game
│   │   │   └── itad_gateway.go
│   │   └── payment
│   │       └── xendit_gateway.go
│   ├── model
│   │   ├── balance_model.go
│   │   ├── dashboard_model.go
│   │   ├── game_model.go
│   │   ├── model.go
│   │   ├── session_model.go
│   │   └── user_model.go
│   ├── pkg
│   │   └── jwt
│   │       └── jwt.go
│   ├── repository
│   │   └── postgresql
│   │       ├── balance_repository.go
│   │       ├── game_repository.go
│   │       ├── session_repository.go
│   │       ├── topup_repository.go
│   │       └── user_repository.go
│   └── usecase
│       ├── balance_usecase.go
│       ├── dashboard_usecase.go
│       ├── game_usecase.go
│       ├── session_usecase.go
│       └── user_usecase.go
├── go.mod
├── go.sum
└── README.md
```

---

## 🏃 3. Running the Server

### 1. Generate Swagger Docs:

```bash
swag init --parseDependency --parseInternal -g cmd/api/main.go -o docs
```

### 2. Start the API Server:

```bash
go run cmd/api/main.go
```

### 3. Open Swagger UI in your browser: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 4. 🧪 Complete API Testing Flow

Use the Swagger UI to run through the entire application lifecycle step-by-step.

### Phase 1: Authentication & Onboarding

#### 1. Register Account

- Endpoint: `POST` `/api/v1/auth/register`
- Payload:

```JSON
{
  "email": "your.email@gmail.com",
  "password": "SecurePassword123!"
}
```

- _Result_: `201` `Created`. Check your real email inbox for the Resend verification link.

#### 2. Verify Email

- Click the "**Verify Account**" link in your email. It will route to `GET` `/api/v1/auth/verify?token=...` and verify your account.

#### 3. Login & Authorize

- Endpoint: `POST` `/api/v1/auth/login`
- Payload (same credentials):

```JSON
{
  "email": "your.email@gmail.com",
  "password": "SecurePassword123!"
}
```

- _Result_: `200` `OK`. Copy the `token` from the response. You will also receive a "**Recent Sign-In**" security alert via email.
- **Authorize Swagger**: Scroll to the top of Swagger, click _Authorize_, and enter `Bearer` `<YOUR_TOKEN>` (usually starts with `ey...`).

#### 4. Check Profile

- Endpoint: `GET` `/api/v1/users/me`
- _Result_: View your verified profile data.

### Phase 2: Catalog & Wallet

#### 5. Sync Game Catalog

- Endpoint: `GET` `/api/v1/games/sync`
- _Result_: Fetches games from IsThereAnyDeal API (or mock fallback) and saves them to the DB.

#### 6. Explore Games

- Endpoint: `GET` `/api/v1/games/explore`
- _Result_: Returns a list of games. Copy one `external_game_id` (e.g., `itad-01`) for purchasing later.

#### 7. Create Top-Up Invoice

- Endpoint: `POST` `/api/v1/games/explore`
- Payload:

```JSON
{
  "amount": 500000
}
```

- _Result_: Returns a Xendit checkout URL and an `invoice_id`. Copy the `invoice_id` (starts with `inv_...` or hex string).

#### 8. Simulate Xendit Webhook (Payment Success)

- Endpoint: `POST` `/api/v1/webhooks/xendit`
- Headers: Add `x-callback-token` if configured in `.env`.
- Payload:

```JSON
{
  "id": "YOUR_COPIED_INVOICE_ID",
  "external_id": "test_external_id",
  "user_id": "test_user_id",
  "status": "PAID",
  "paid_amount": 500000,
  "payer_email": "your.email@gmail.com",
  "payment_method": "QRIS"
}
```

- _Result_: `200` `OK`. You will receive a Top-Up Confirmed email receipt. Check your new balance using `GET` `/api/v1/balances`.

### Phase 3: Purchases & Gameplay

#### 9. Buy a Game

- Endpoint: `POST` `/api/v1/games/buy`
- Payload:

```JSON
{
  "external_game_id": "itad-01"
}
```

- _Result_: Deducts balance and adds the game to your library.

#### 10. View Personal Library

- Endpoint: `GET` `/api/v1/games/library`
- _Result_: Shows owned games and logged playtime.

#### 11. Launch Game (Start Session)

- Endpoint: `POST` `/api/v1/sessions/launch/{game_id}`
- Parameter: `itad-01`
- _Result_: Starts a 60-second active session lock. Hitting this again immediately will return a `409` `Conflict`, enforcing the "1 game at a time" rule.

#### 12. Check Active Session

- Endpoint: `GET` `/api/v1/sessions/active`
- _Result_: Shows the countdown timer until your current session expires.

### Phase 4: Analytics & Teardown

#### 13. View Dashboard Statistics

- Endpoint: `GET` `/api/v1/dashboard`
- _Result_: Displays current wallet balance, total owned games, total lifetime spend, and a list of games played this month.

#### 14. Logout

- Endpoint: `POST` `/api/v1/auth/logout`
- _Result_: Safely invalidates the client session.
