# Puna - Badminton Court Reservation System

A modern web application for managing badminton court reservations with an elegant admin dashboard. Built with Go, featuring a responsive Bulma-based interface and comprehensive reservation management system.

![Puna Dashboard](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)

## 🏸 Features

### Public Features
- **Court Availability Search**: Real-time availability checking
- **Reservation System**: Easy court booking with date/time selection
- **Court Information**: Detailed court specifications and amenities
- **User Registration & Authentication**: Secure user accounts
- **Responsive Design**: Mobile-friendly interface

### Admin Dashboard
- **Modern Admin Interface**: Professional Bulma-based dashboard
- **Reservation Management**: View and manage all bookings
- **User Management**: Administer user accounts
- **Analytics Dashboard**: Performance metrics and charts
- **Court Management**: Add, edit, and manage court information

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/puna.git
   cd puna
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   # Create PostgreSQL database
   createdb puna_db
   
   # Run migrations (if any)
   # TODO: Add migration commands
   ```

4. **Configure environment**
   ```bash
   # Copy example config
   cp .env.example .env
   
   # Edit configuration
   nano .env
   ```

5. **Run the application**
   ```bash
   go run cmd/web/main.go
   ```

6. **Access the application**
   - Main site: http://localhost:8080
   - Admin dashboard: http://localhost:8080/dashboard

## 🛠️ Development

### Project Structure
```
puna/
├── cmd/web/           # Application entry point
├── internal/          # Internal application code
│   ├── config/       # Configuration management
│   ├── handlers/     # HTTP handlers
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   └── render/       # Template rendering
├── templates/        # HTML templates
├── static/          # Static assets (CSS, JS, images)
└── frontend/        # Frontend assets
```

### Building for Production

1. **Build the application**
   ```bash
   go build -o bin/puna cmd/web/main.go
   ```

2. **Run the binary**
   ```bash
   ./bin/puna
   ```

### Development Commands

```bash
# Run with hot reload (requires air)
air

# Run tests
go test ./...

# Run with race detection
go run -race cmd/web/main.go

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o bin/puna-linux cmd/web/main.go
GOOS=windows GOARCH=amd64 go build -o bin/puna-windows.exe cmd/web/main.go
GOOS=darwin GOARCH=amd64 go build -o bin/puna-macos cmd/web/main.go
```

## 📊 Admin Dashboard

The admin dashboard provides comprehensive management capabilities:

- **Dashboard Overview**: Key metrics and performance indicators
- **Reservation Management**: View, edit, and cancel bookings
- **User Administration**: Manage user accounts and permissions
- **Court Configuration**: Add and configure court details
- **Analytics**: Performance charts and reporting

### Dashboard Features
- Real-time data visualization with Chart.js
- Responsive design for all devices
- Material Design Icons for intuitive navigation
- Modal dialogs for quick actions
- Sortable and searchable data tables

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=puna_db

# Server
PORT=8080
ENV=development

# Security
SECRET_KEY=your_secret_key
```

### Database Setup
```sql
-- Create database
CREATE DATABASE puna_db;

-- Create user (optional)
CREATE USER puna_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE puna_db TO puna_user;
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers -v

# Run benchmarks
go test -bench=. ./...
```

## 📦 Deployment

### Docker Deployment
```bash
# Build Docker image
docker build -t puna .

# Run container
docker run -p 8080:8080 puna
```

### Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=puna_db
      - POSTGRES_USER=puna_user
      - POSTGRES_PASSWORD=your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go coding standards
- Write tests for new features
- Update documentation as needed
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Third-Party Licenses
This project uses several third-party libraries. See [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md) for complete license information.

## 🙏 Acknowledgments

- [Admin One Bulma Dashboard](https://github.com/vikdiesel/admin-one-bulma-dashboard) - Admin template
- [Bulma CSS Framework](https://bulma.io/) - CSS framework
- [Chart.js](https://www.chartjs.org/) - Charting library
- [Material Design Icons](https://materialdesignicons.com/) - Icon font

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/puna/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/puna/discussions)
- **Email**: support@yourdomain.com

## 🔄 Changelog

### v1.0.0 (2024-01-XX)
- Initial release
- Basic reservation system
- Admin dashboard
- User authentication
- Court management

---

**Made with ❤️ using Go and Bulma**