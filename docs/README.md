# Worksquare Housing Listings API

A comprehensive RESTful API built with Go Fiber for managing housing listings with JWT authentication, pagination, filtering, and rate limiting.

## 🚀 Features

- **RESTful API** with clean architecture
- **JWT Authentication** with access and refresh tokens
- **Rate Limiting** (100 requests per hour per IP)
- **Pagination** with configurable page size
- **Advanced Filtering** by location, property type, price range, bedrooms, bathrooms
- **Request Logging** middleware
- **Swagger Documentation** with OpenAPI 3.0
- **Docker Support** with multi-stage builds
- **Unit & Integration Tests**
- **Standardized Error Handling**

## 📋 Requirements

- Go 1.22+
- Docker & Docker Compose (optional)
- Make (optional, for development commands)

## 🛠️ Setup Instructions

### 1. Clone and Setup

```bash
# Clone the project
git clone git@github.com:conquext/worksquare-backend-task.git

# Navigate to project directory
cd worksquare-backend-task

# Install dependencies
make deps
# or
go mod download
go mod tidy
```

### 2. Environment Configuration

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your configurations
nano .env
```

### 3. Add Listings Data

Copy the provided `listings.json` file to the `data/` directory:

```bash
cp /path/to/your/listings.json data/listings.json
```

### 4. Run the Application

#### Using Go directly:

```bash
# Development mode
make dev
# or
go run main.go

# Production build
make build && make run
```

#### Using Docker:

```bash
# Docker Compose (recommended)
make docker-compose

# Or build and run manually
make docker-build
make docker-run
```

## 📚 API Documentation

### Base URL

```
http://localhost:3000/api/v1
```

### Swagger Documentation

```
http://localhost:3000/swagger/
```

### Authentication

#### Demo Credentials

```json
{
  "email": "demo@worksquare.com",
  "password": "demo123456"
}
```

#### Get Demo Credentials

```http
GET /api/v1/demo/credentials
```

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "demo@worksquare.com",
  "password": "demo123456"
}
```

#### Register

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "securepassword"
}
```

#### Refresh Token

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token_here"
}
```

### Listings Endpoints

#### Get All Listings (Paginated)

```http
GET /api/v1/listings?page=1&limit=10&location=Lagos&property_type=House&min_price=1000000&max_price=5000000
```

#### Get Listing by ID

```http
GET /api/v1/listings/1
```

#### Search Listings

```http
GET /api/v1/listings/search?q=Lagos&page=1&limit=10
```

#### Get Filter Metadata

```http
GET /api/v1/listings/filters
```

#### Get Listing Statistics (Protected)

```http
GET /api/v1/listings/stats
Authorization: Bearer <your_jwt_token>
```

### Query Parameters

#### Pagination

- `page` (int): Page number (default: 1, min: 1)
- `limit` (int): Items per page (default: 10, min: 1, max: 100)

#### Filtering

- `location` (string): Filter by location (partial match)
- `property_type` (string): Filter by property type (exact match)
- `city` (string): Filter by city (partial match)
- `min_price` (int): Minimum price filter
- `max_price` (int): Maximum price filter
- `min_bedrooms` (int): Minimum number of bedrooms
- `max_bedrooms` (int): Maximum number of bedrooms
- `min_bathrooms` (int): Minimum number of bathrooms
- `max_bathrooms` (int): Maximum number of bathrooms

## 🏗️ Architecture & Design

### Project Structure

```
worksquare-backend-task/
├── cmd/server/           # Application entry points
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── controllers/     # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models and structures
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic layer
│   └── utils/           # Utility functions
├── pkg/                 # Public packages
│   ├── jwt/            # JWT utilities
│   ├── logger/         # Logging utilities
│   ├── pagination/     # Pagination helpers
│   └── response/       # HTTP response helpers
├── api/                # API layer
│   └── routes/         # Route definitions
├── tests/              # Test files
├── data/               # Data files
```

### Architecture Layers

1. **HTTP Layer** (`api/routes`, `internal/controllers`)

   - Route definitions and HTTP request handling
   - Request validation and response formatting

2. **Business Logic Layer** (`internal/services`)

   - Core business logic and rules
   - Data transformation and validation

3. **Data Access Layer** (`internal/repositories`)

   - Data retrieval and manipulation
   - File I/O operations for JSON data

4. **Models Layer** (`internal/models`)
   - Data structures and DTOs
   - Request/response models

### Authentication & Security

- **JWT Tokens**: Stateless authentication with access and refresh tokens
- **Password Hashing**: Bcrypt with salt for secure password storage
- **Rate Limiting**: IP-based rate limiting to prevent abuse
- **CORS**: Configurable cross-origin resource sharing
- **Security Headers**: Helmet middleware for security headers
- **Input Validation**: Comprehensive request validation

### Error Handling Strategy

- **Standardized Responses**: Consistent API response format
- **Error Codes**: HTTP status codes with detailed error messages
- **Validation Errors**: Field-level validation error reporting
- **Logging**: Comprehensive error logging for debugging

### Data Model

#### Listing Model

```go
type Listing struct {
    ID         int      `json:"id"`
    Title      string   `json:"title"`
    Price      string   `json:"price"`
    Bedrooms   int      `json:"bedrooms"`
    Bathrooms  int      `json:"bathrooms"`
    Location   string   `json:"location"`
    Status     []string `json:"status"`
    Image      string   `json:"image"`
}
```

#### Filter Options

- Location-based filtering (partial match)
- Property type filtering (House, Flat, Terrace, etc.)
- Price range filtering with numeric conversion
- Bedroom/bathroom count filtering
- City-specific filtering

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Unit tests only
go test ./tests/unit/...

# Integration tests only
go test ./tests/integration/...
```

### Test Structure

- **Unit Tests**: Service layer testing with mocks
- **Integration Tests**: End-to-end API testing
- **Coverage Reports**: HTML coverage reports generated

## 📦 Deployment

### Docker Deployment

```bash
# Build and run with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Considerations

- Set strong JWT secrets in production
- Configure appropriate rate limits
- Use HTTPS in production
- Set up proper logging and monitoring
- Consider using a real database for user storage

## 🛠️ Development Tools

### Available Make Commands

```bash
make dev              # Start development server
make build            # Build binary
make run              # Build and run
make test             # Run tests
make test-coverage    # Run tests with coverage
make docker-build     # Build Docker image
make docker-run       # Run Docker container
make docker-compose   # Start with Docker Compose
make clean            # Clean build artifacts
make fmt              # Format code
make lint             # Run linter
```

### Code Quality

- **Go Modules**: Dependency management
- **golangci-lint**: Code linting and static analysis
- **gofmt**: Code formatting
- **Structured Logging**: JSON-formatted logs
- **Error Wrapping**: Detailed error context

## 📊 API Response Format

### Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Validation failed",
    "details": "Invalid request parameters"
  },
  "data": null
}
```

### Validation Error Response

```json
{
  "success": false,
  "error": {
    "code": 422,
    "message": "Validation failed"
  },
  "data": {
    "message": "Validation failed",
    "errors": [
      {
        "field": "email",
        "message": "Please provide a valid email address",
        "value": "invalid-email"
      }
    ]
  }
}
```

## 🔧 Configuration

### Environment Variables

- `NODE_ENV`: Environment (development/production)
- `PORT`: Server port (default: 3000)
- `JWT_SECRET`: JWT signing secret
- `JWT_EXPIRES_IN`: Access token expiry (default: 24h)
- `JWT_REFRESH_EXPIRES_IN`: Refresh token expiry (default: 7d)
- `RATE_LIMIT_MAX_REQUESTS`: Rate limit per window (default: 100)
- `RATE_LIMIT_WINDOW_MS`: Rate limit window (default: 1h)
- `LOG_LEVEL`: Logging level (debug/info/warn/error)

## 📈 Performance Considerations

- **In-Memory Data**: JSON file loaded into memory for fast access
- **Efficient Filtering**: Optimized filtering algorithms
- **Pagination**: Memory-efficient pagination implementation
- **Rate Limiting**: Prevents API abuse and ensures fair usage
- **Structured Logging**: Efficient logging with minimal performance impact

## 🚀 Future Enhancements

1. **Database Integration**: Replace JSON file with PostgreSQL/MongoDB
2. **Caching**: Redis caching for frequently accessed data
3. **Search Engine**: Elasticsearch for advanced search capabilities
4. **File Upload**: Image upload functionality for listings
5. **Real-time Updates**: WebSocket support for real-time listing updates
6. **Admin Panel**: Administrative interface for managing listings
7. **Analytics**: API usage analytics and reporting
8. **Geolocation**: Geographic search and mapping features

## 🐛 Troubleshooting

### Common Issues

1. **Port Already in Use**

   ```bash
   # Change port in .env file or kill existing process
   lsof -ti:3000 | xargs kill -9
   ```

2. **Missing listings.json**

   ```bash
   # Ensure listings.json is in data/ directory
   ls -la data/listings.json
   ```

3. **JWT Token Expired**
   ```bash
   # Use refresh token endpoint to get new access token
   curl -X POST http://localhost:3000/api/v1/auth/refresh \
   -H "Content-Type: application/json" \
   -d '{"refresh_token":"your_refresh_token"}'
   ```

## 📞 Support

For questions or issues:

1. Check the Swagger documentation at `/swagger/`
2. Review the logs for error details
3. Ensure all environment variables are properly set
4. Verify the listings.json file is present and valid

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤖 Note on AI Collaboration

## Some parts of this project — including documentation and test scaffolding — were written with the assistance of AI tools to improve productivity and ensure consistency. All outputs have been reviewed and validated by the author.
