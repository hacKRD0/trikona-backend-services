# User Management Service

A microservice for handling user authentication, registration, and profile management. Built with Go and following clean architecture principles.

## Features

- User registration and authentication
- Email verification
- Password reset functionality
- Profile management
- Role-based access control
- JWT-based authentication
- LinkedIn OAuth integration

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL database
- Mailjet account (for email services)
- LinkedIn OAuth credentials

## Configuration

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=user_management

# JWT Configuration
JWT_SECRET=your_jwt_secret

# Email Configuration (Mailjet)
MAILJET_API_KEY=your_mailjet_key
MAILJET_SECRET_KEY=your_mailjet_secret
MAILJET_FROM_EMAIL=your_email
MAILJET_FROM_NAME=your_name

# Frontend Configuration
FRONTEND_URL=http://localhost:3000

# LinkedIn OAuth
LINKEDIN_CLIENT_ID=your_client_id
LINKEDIN_CLIENT_SECRET=your_client_secret
LINKEDIN_REDIRECT_URI=http://localhost:8080/auth/linkedin/callback
```

## Local Development

1. Clone the repository:
```bash
git clone https://github.com/your-username/trikona_go.git
cd trikona_go
```

2. Install dependencies:
```bash
go mod download
```

3. Start the service:
```bash
go run cmd/user-management-service/main.go
```

The service will be available at `http://localhost:8080`.

## API Endpoints

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - User login
- `POST /auth/verify` - Verify email
- `POST /auth/reset-password` - Request password reset
- `POST /auth/reset-password/confirm` - Confirm password reset
- `GET /auth/linkedin` - Get LinkedIn OAuth URL

### User Management
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile
- `DELETE /users/profile` - Delete user account

## Docker Deployment

1. Build the Docker image:
```bash
docker build -t user-management-service -f cmd/user-management-service/Dockerfile .
```

2. Run the container:
```bash
docker run -d \
  -p 8080:8080 \
  --env-file .env \
  user-management-service
```

## Architecture

The service follows clean architecture principles with the following structure:

```
.
├── cmd/
│   └── user-management-service/
│       └── main.go
├── internal/
│   └── user-management-service/
│       ├── domain/
│       ├── repository/
│       └── usecase/
├── pkg/
│   ├── auth/
│   ├── errors/
│   ├── logger/
│   └── validation/
└── configs/
```

- `cmd/`: Application entry points
- `internal/`: Service-specific code
  - `domain/`: Business entities and interfaces
  - `repository/`: Data access layer
  - `usecase/`: Business logic
- `pkg/`: Shared packages
  - `auth/`: Authentication utilities
  - `errors/`: Custom error types
  - `logger/`: Logging utilities
  - `validation/`: Input validation

## Logging

The service uses structured logging with the following levels:
- INFO: General operational information
- ERROR: Error conditions that need attention
- DEBUG: Detailed information for debugging

## Error Handling

The service uses custom error types for consistent error responses:
- ValidationError (400)
- AuthenticationError (401)
- AuthorizationError (403)
- NotFoundError (404)
- ConflictError (409)
- InternalError (500)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
