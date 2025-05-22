# Trikona Backend Services

A comprehensive backend system built with Go, following clean architecture principles. The system consists of multiple microservices for user management, directory services, and more.

## Features

### User Management Service
- User registration and authentication
- Email verification
- Password reset functionality
- Profile management
- Role-based access control
- JWT-based authentication
- LinkedIn OAuth integration

### Directory Service
#### Student Management
- CRUD operations for student profiles
- Advanced filtering and pagination
- Skills management
- Education history tracking

#### Professional Directory
- Professional profile management
- Advanced search with multiple filters
- Skills and experience tracking
- Paginated results

#### Corporate Management
- Corporate profile management
- Office locations
- Industry classification

#### College Management
- College profile management
- Branch and department tracking
- Student enrollment

#### Master Data Management
- Countries and states
- Industries and sectors
- Skills and services catalogs

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- PostgreSQL database
- Mailjet account (for email services)
- LinkedIn OAuth credentials
- Git

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

### Directory Service

#### Students
- `GET /api/v1/directory/students` - List all students (with filtering)
- `GET /api/v1/directory/students/:id` - Get a specific student by ID
- `POST /api/v1/directory/students` - Create a new student
- `PUT /api/v1/directory/students/:id` - Update an existing student
- `DELETE /api/v1/directory/students/:id` - Delete a student

#### Professionals
- `GET /api/v1/directory/professionals` - List all professionals (with filtering)
- `GET /api/v1/directory/professionals/:id` - Get a specific professional by ID
- `POST /api/v1/directory/professionals` - Create a new professional
- `PUT /api/v1/directory/professionals/:id` - Update a professional
- `DELETE /api/v1/directory/professionals/:id` - Delete a professional

#### Corporates
- `GET /api/v1/directory/corporates` - List all corporates (with filtering)
- `GET /api/v1/directory/corporates/:id` - Get a specific corporate by ID
- `POST /api/v1/directory/corporates` - Create a new corporate
- `PUT /api/v1/directory/corporates/:id` - Update a corporate
- `DELETE /api/v1/directory/corporates/:id` - Delete a corporate

#### Colleges
- `GET /api/v1/directory/colleges` - List all colleges (with filtering)
- `GET /api/v1/directory/colleges/:id` - Get a specific college by ID
- `POST /api/v1/directory/colleges` - Create a new college
- `PUT /api/v1/directory/colleges/:id` - Update a college
- `DELETE /api/v1/directory/colleges/:id` - Delete a college

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

The system follows clean architecture principles with the following structure:

```
.
├── api/
│   └── directory-service/
│       ├── college/
│       ├── corporate/
│       ├── professional/
│       └── student/
├── cmd/
│   ├── directory-service/
│   │   └── main.go
│   └── seed/
├── internal/
│   ├── directory-service/
│   │   ├── domain/
│   │   ├── repository/
│   │   └── usecase/
│   └── user-management-service/
│       ├── domain/
│       ├── repository/
│       └── usecase/
├── pkg/
│   ├── auth/
│   ├── config/
│   ├── database/
│   ├── errors/
│   ├── logger/
│   ├── middleware/
│   ├── pagination/
│   ├── utils/
│   └── validation/
├── migrations/
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
3. Commit your changes with descriptive commit messages
4. Push to the branch
5. Create a Pull Request with a clear description of changes

## Related Repositories

- [Trikona Frontend](https://github.com/hacKRD0/trikona-frontend)

## Support

For support or feature requests, please open an issue in the [GitHub repository](https://github.com/hacKRD0/trikona-backend-services/issues).
