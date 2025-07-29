# URL-Shortener

Turn long URLs into short ones! 

## 🚀 What it does

- Shorten long URLs into tiny links
- User accounts with login/signup
- See all your shortened links

## 🛠️ Built with

- **Go**
- **Fiber** 
- **PostgreSQL**
- **JWT** - For user authentication

## 🏃‍♂️ Running it locally

1. **Clone it**
   ```bash
   git clone https://github.com/electrofull/URL-Shortener.git
   cd URL-Shortener/src
   ```

2. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL**
   - Create a database called `URL-Shortener`
   - Update the `.env` file with your database credentials

4. **Environment setup**
   ```env
   PORT=8000
   DATABASE_URL="postgres://username:password@localhost:5432/URLShortener"
   JWT_SECRET="your-secret-key"
   BASE_URL="http://localhost:8000"
   ```

5. **Run it**
   ```bash
   go run main.go
   ```



## 📡 API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message |
| POST | `/register` | Create new user account |
| POST | `/login` | User login |
| GET | `/:shortcode` | Redirect to original URL |

### Protected Endpoints (Require JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/shorten` | Create shortened URL |
| GET | `/all` | Get all user's URLs |

## 🔧 API Usage Examples

### Register a new user
```bash
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=johndoe&password=secretpassword"
```

### Login
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=johndoe&password=secretpassword"
```

### Create a short URL
```bash
curl -X POST http://localhost:8000/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"url": "https://www.google.com"}'
```

### Get all your URLs
```bash
curl -X GET http://localhost:8000/all \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Use a short URL
```bash
# Visit in browser or curl
curl -L http://localhost:8000/abc123
```

## 🔐 Security Features

- **Password Hashing**: Uses bcrypt with default cost
- **JWT Authentication**: Secure token-based authentication
- **Input Validation**: URL format validation
- **SQL Injection Protection**: Parameterized queries