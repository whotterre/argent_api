# Argent Wallet API

The Argent Wallet API is a wallet API built for the penultimate stage of the HNG13 backend internship. It allows users to deposit money using Paystack, manage wallet balances, view transaction history, and transfer funds to other users.

### Prerequisites
- Go 1.21+
- PostgreSQL 14
- Paystack Account (for API keys)
- Google Cloud Console Project (for OAuth credentials)

### Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/argent_api.git
   cd argent_api
   ```

2. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   PORT=8080
   DB_URL=postgres://user:password@localhost:5432/argent_db?sslmode=disable
   
   # Google Auth
   GOOGLE_CLIENT_ID=your_google_client_id
   GOOGLE_CLIENT_SECRET=your_google_client_secret
   GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
   JWT_SECRET=your_jwt_secret

   # Paystack
   PAYSTACK_SECRET_KEY=your_paystack_secret_key
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run Database Migrations**
   (Ensure your PostgreSQL database is running and accessible)
   ```bash
   # If using a migration tool like golang-migrate or internal migration script
   go run main.go migrate
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

## Functional Requirements & Endpoints

### 1. Google Authentication (JWT)
- **GET /auth/google**: Triggers Google sign-in.
- **GET /auth/google/callback**: Logs in the user, creates user if not existing, returns a JWT token.

### 2. API Key Management
**Rules:**
- Max 5 active keys per user.
- Expiry accepts: `1H`, `1D`, `1M`, `1Y`. Backend converts to `expires_at`.
- Permissions must be explicitly assigned.

#### a. Create API Key
- **POST /keys/create**
- Request:
  ```json
  {
    "name": "wallet-service",
    "permissions": ["deposit", "transfer", "read"],
    "expiry": "1D"
  }
  ```
- Response:
  ```json
  {
    "api_key": "sk_live_xxxxx",
    "expires_at": "2025-01-01T12:00:00Z"
  }
  ```

#### b. Rollover Expired API Key
- **POST /keys/rollover**
- Purpose: Create a new API key using the same permissions as an expired key.
- Request:
  ```json
  {
    "expired_key_id": "FGH2485K6KK79GKG9GKGK",
    "expiry": "1M"
  }
  ```
- **Rules:**
  - The expired key must truly be expired.
  - The new key must reuse the same permissions.
  - Expiry must again be converted to a new `expires_at` value.

### 3. Wallet Deposit (Paystack)
- **POST /wallet/deposit**
- Auth: JWT or API Key with `deposit` permission.
- Request: `{ "amount": 5000 }`
- Response:
  ```json
  {
    "reference": "...",
    "authorization_url": "https://paystack.co/checkout/..."
  }
  ```

### 4. Paystack Webhook (Mandatory)
- **POST /wallet/paystack/webhook**
- Purpose: Receive transaction updates from Paystack. Credit wallet only after webhook confirms success.
- Security: Validate Paystack signature.
- Actions: Verify signature, find transaction by reference, update transaction status and wallet balance.

### 5. Verify Deposit Status
- **GET /wallet/deposit/{reference}/status**
- Response: `{ "reference": "...", "status": "success|failed|pending", "amount": 5000 }`
- *Warning: This endpoint must not credit wallets. Only the webhook is allowed to credit wallets.*

### 6. Get Wallet Balance
- **GET /wallet/balance**
- Auth: JWT or API key with `read` permission.
- Response: `{ "balance": 15000 }`

### 7. Wallet Transfer
- **POST /wallet/transfer**
- Auth: JWT or API key with `transfer` permission.
- Request:
  ```json
  {
    "wallet_number": "4566678954356",
    "amount": 3000
  }
  ```
- Response:
  ```json
  {
    "status": "success",
    "message": "Transfer completed"
  }
  ```

### 8. Transaction History
- **GET /wallet/transactions**
- Auth: JWT or API key with `read` permission.
- Response:
  ```json
  [
    { "type": "deposit", "amount": 5000, "status": "success" },
    { "type": "transfer", "amount": 3000, "status": "success" }
  ]
  ```

## Access Rules & Security

### Access Rules
- **Authorization: Bearer <token>**: Treat as user (can perform all actions).
- **x-api-key: <key>**: Treat as service.
- API keys must have valid permissions and not be expired/revoked.

### Security Considerations
- Do not expose secret keys.
- Validate Paystack webhooks.
- Do not allow transfers with insufficient balance.
- Do not allow API keys without correct permissions.
- Expired API keys must be rejected automatically.

### Error Handling & Idempotency
- Paystack reference must be unique.
- Webhooks must be idempotent (no double-credit).
- Transfers must be atomic (no partial deductions).
- Return clear errors for: insufficient balance, invalid/expired API key, missing permissions.
