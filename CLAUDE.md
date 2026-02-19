# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Databasus is a self-hosted database backup management tool supporting PostgreSQL, MySQL, MariaDB, and MongoDB. It provides scheduled backups, multiple storage destinations (S3, Azure, Google Drive, etc.), and notifications (Email, Telegram, Slack, Discord, etc.).

**Architecture:**
- **Backend:** Go 1.24.4 with Gin framework, GORM ORM, PostgreSQL 17 (internal DB)
- **Frontend:** React 19.1 with TypeScript, Vite 6.3, Ant Design UI, Tailwind CSS 4.1
- **Deployment:** Docker-first with multi-stage builds, Helm charts for Kubernetes

## Common Commands

### Backend Development
```bash
# Run backend (from backend/ directory)
make run

# Run tests
make test

# Run linting
make lint

# Database migrations
make migration-create name=MIGRATION_NAME
make migration-up
make migration-down

# Generate Swagger docs
make swagger
# Access at http://localhost:4005/api/v1/docs/swagger/index.html#/
```

**Important:** Use `dev-db` from `backend/docker-compose.yml` for development, not `databasus-db` from root.

### Frontend Development
```bash
# From frontend/ directory
npm run dev          # Start dev server
npm run build        # Production build
npm run lint         # Check linting
npm run format       # Format with Prettier
npm run test         # Run tests (CI mode)
npm run test:watch   # Watch mode
```

### Docker Build
```bash
# Build full image (frontend + backend)
docker build -t databasus:latest .

# Run container
docker run -d --name databasus -p 4005:4005 -v ./databasus-data:/databasus-data databasus:latest
```

## Code Architecture

### Backend Structure (Clean Architecture)

Each feature follows this structure under `backend/internal/features/`:
```
feature/
├── controller.go      # HTTP handlers (Gin)
├── service.go         # Business logic
├── repository.go      # Data access (GORM)
├── model.go           # Database models
├── dto.go             # Data transfer objects
├── di.go              # Dependency injection (singletons)
└── *_test.go          # Unit tests
```

**Key Features:**
- `databases/` - Multi-database support (PostgreSQL, MySQL, MariaDB, MongoDB)
- `backups/` - Backup orchestration with scheduling and encryption
- `restores/` - Restore operations
- `storages/` - 8 storage backends (S3, Azure, Google Drive, FTP, SFTP, NAS, local, rclone)
- `notifiers/` - Notification channels (Email, Telegram, Slack, Discord, webhook, MS Teams)
- `users/` - Authentication, authorization, user management
- `workspaces/` - Multi-tenancy with role-based access control

### Frontend Structure
```
frontend/src/
├── entity/           # API clients & data models (11 entities)
├── features/         # Feature-specific UI components
├── pages/            # Route-level components
└── shared/           # Reusable utilities
    ├── api/          # HTTP client with auth
    ├── ui/           # Shared UI components
    ├── theme/        # Dark/light theme
    ├── hooks/        # Custom React hooks
    └── lib/          # Utility functions
```

## Code Style Rules (Backend)

### Method Organization
**Always place private methods at the bottom of ALL Go files.**

Order:
1. Type definitions and constants
2. Public methods/functions (uppercase)
3. Private methods/functions (lowercase)

Example:
```go
type UserService struct {}

// Public methods first
func (s *UserService) CreateUser(user *User) error { ... }

// Private methods at bottom
func (s *UserService) validateUser(user *User) error { ... }
```

### Controllers
- Name all route handlers using `.WhatWeDo` pattern (not "handlers")
- Use `*gin.Context` for all routes
- Combine all routes in single controller
- Document with Swagger annotations (see example in `.cursor/rules/controllers-rule.mdc`)

### Dependency Injection
Use **implicit field declaration style** in `di.go` files:

```go
var orderController = &OrderController{
    orderService,                    // ✅ Correct
    bot_users.GetBotUserService(),   // ✅ Correct
}
```

NOT:
```go
var orderController = &OrderController{
    orderService: orderService,      // ❌ Wrong
    botUserService: bot_users.GetBotUserService(),  // ❌ Wrong
}
```

This prevents forgetting to update DI when adding dependencies.

### CRUD Pattern
Follow the complete CRUD example in `.cursor/rules/crud.mdc`, which includes:
- Controller with Swagger docs
- Service with business logic
- Repository with GORM queries
- DTOs for request/response
- Models with GORM tags
- Unit tests for controller and service
- DI setup with singleton getters

## Code Style Rules (Frontend)

### React Component Structure
Follow this exact order:

```tsx
interface Props {
   someValue: SomeValue;
}

const someHelperFunction = () => {
    ...
}

export const ReactComponent = ({ someValue }: Props): JSX.Element => {
    // 1. States first
    const [someState, setSomeState] = useState<...>(...)

    // 2. Functions
    const loadSomeData = async () => {
        ...
    }

    // 3. Hooks
    useEffect(() => {
       loadSomeData();
    });

    // 4. Calculated values
    const calculatedValue = someValue.calculate();

    return <div> ... </div>
}
```

See `.cursor/rules/react-component-structure.mdc` for details.

## Key Technical Details

### Database Support
Uses **native database client binaries** (not language drivers):
- PostgreSQL: pg_dump/pg_restore (versions 12-18)
- MySQL: mysqldump (5.7, 8.0, 8.4, 9)
- MariaDB: mariadb-dump (10.6, 12.1)
- MongoDB: mongodump/mongorestore (v4-v8)

Binaries are embedded in Docker image at `/assets/`.

### Storage Architecture
All storages implement common interface:
- `UploadBackup()` - Upload encrypted backup
- `DownloadBackup()` - Download for restore
- `DeleteBackup()` - Cleanup old backups
- `ListBackups()` - Enumerate stored backups

### Security
- **Encryption:** AES-256-GCM for backup files
- **Field Encryption:** Sensitive credentials encrypted in database
- **Read-Only DB Users:** Backups use read-only accounts
- **Workspace Isolation:** Multi-tenancy with RBAC
- **Audit Logging:** All actions tracked

### Backup Flow
1. Scheduler triggers backup (robfig/cron v3)
2. Connect with read-only database user
3. Run native dump tool
4. Compress (4-8x ratio)
5. Encrypt with AES-256-GCM
6. Upload to storage(s)
7. Send notifications
8. Cleanup old backups based on retention policy

### API Design
- RESTful with `/api/v1` prefix
- JWT authentication middleware
- Swagger auto-generated via swaggo/swag
- CORS and GZIP middleware

### Testing
- **Backend:** Standard Go testing + testify, run with `make test` (uses `-count=1` to disable caching)
- **Frontend:** Vitest with coverage
- **CI/CD:** GitHub Actions with lint, test, and release jobs

## Configuration

### Backend
- Environment-based config via `.env` files
- Examples: `.env.development.example`, `.env.production.example`
- Config loaded with `cleanenv` library
- Key settings: Database DSN, paths to DB client binaries, test ports

### Frontend
- Vite env variables (prefixed with `VITE_APP_`)
- Example: `VITE_APP_VERSION`

## Migration from Postgresus

Databus was renamed from Postgresus. They use different:
- Data folders
- Internal database naming
- Docker image names

See README.md "Migration guide" section for details.

## Additional Notes

- Use Context7 for up-to-date documentation of libraries (per user instructions)
- Generate minimal code only; remove all unused variables (per user instructions)
- All features use singleton pattern with `Get*Service()`/`Get*Controller()` getters
- Backend uses PostgreSQL 17 as internal database
- Temp folders must have permissions 0700 for PostgreSQL .pgpass ownership requirements
- MySQL/MariaDB: PROCESS permission is not mandatory for backups

See `.cursor/rules/` directory for complete project rules including:
- Comments style
- Controller patterns
- Time handling
- Test patterns
- Refactoring guidelines
- Migration rules
