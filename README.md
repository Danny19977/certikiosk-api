# CertiKiosk API

A robust Go-based REST API for managing sales data, user authentication, and comprehensive dashboard analytics.

## Features

- 🔐 **Authentication & Authorization**: JWT-based authentication with role-based access control
- 👥 **User Management**: Complete CRUD operations for users with country and province associations
- 📊 **Sales Tracking**: Multi-level sales data management (Product, Province, Year, Month, Week)
- 📈 **Advanced Dashboards**:
  - Global Overview Dashboard
  - Provincial Analysis Dashboard
  - Daily Operations Monitor
  - Historical Trends Analysis
- 🔔 **Notifications System**: User notifications with status tracking
- 📝 **Activity Logging**: Comprehensive user activity tracking
- 🌍 **Geographic Organization**: Country and Province-based data hierarchy

## Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) - Express-inspired web framework
- **Database**: PostgreSQL with [GORM](https://gorm.io/) ORM
- **Authentication**: JWT with bcrypt password hashing
- **Validation**: go-playground/validator

## Project Structure

```
certikiosk/
├── controller/
│   ├── auth/              # Authentication & password management
│   ├── country/           # Country management
│   ├── dashboard/         # Analytics dashboards
│   ├── month/             # Monthly targets
│   ├── Notification/      # User notifications
│   ├── product/           # Product management
│   ├── province/          # Province management
│   ├── sale/              # Sales transactions
│   ├── user/              # User management
│   ├── userlog/           # Activity logs
│   ├── week/              # Weekly targets
│   └── year/              # Yearly targets
├── database/              # Database connection
├── middlewares/           # Authentication middleware
├── models/                # Data models
├── routes/                # API routes
├── utils/                 # Helper functions
└── examples/              # Example payloads
```

## Getting Started

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- SMTP server for email notifications (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Danny19977/certikiosk.git
cd certikiosk
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory:
```env
SECRET_KEY=your_secret_key_here

# Server Configuration
PORT=8000

# Database Configuration
DB_HOST=127.0.0.1
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=certikioskdb
DB_PORT=5432

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_MAIL=your-email@gmail.com
SMTP_PASSWORD=your_app_password
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8000`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/user` - Get authenticated user
- `PUT /api/auth/change-password` - Change password
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset/:token` - Reset password
- `POST /api/auth/logout` - Logout

### Users
- `GET /api/users/all` - Get all users
- `GET /api/users/all/paginate` - Get paginated users
- `GET /api/users/get/:uuid` - Get user by UUID
- `POST /api/users/create` - Create new user
- `PUT /api/users/update/:uuid` - Update user
- `DELETE /api/users/delete/:uuid` - Delete user

### Countries
- `GET /api/countries/all` - Get all countries
- `GET /api/countries/get/:uuid` - Get country by UUID
- `POST /api/countries/create` - Create country
- `PUT /api/countries/update/:uuid` - Update country
- `DELETE /api/countries/delete/:uuid` - Delete country

### Provinces
- `GET /api/provinces/all` - Get all provinces
- `GET /api/provinces/all/country/:country_uuid` - Get provinces by country
- `GET /api/provinces/get/:uuid` - Get province by UUID
- `POST /api/provinces/create` - Create province
- `PUT /api/provinces/update/:uuid` - Update province
- `DELETE /api/provinces/delete/:uuid` - Delete province

### Products
- `GET /api/products/all` - Get all products
- `GET /api/products/get/:uuid` - Get product by UUID
- `POST /api/products/create` - Create product
- `PUT /api/products/update/:uuid` - Update product
- `DELETE /api/products/delete/:uuid` - Delete product

### Sales
- `GET /api/sales/all` - Get all sales
- `GET /api/sales/all/province/:province_uuid` - Get sales by province
- `GET /api/sales/get/:uuid` - Get sale by UUID
- `POST /api/sales/create` - Create sale
- `PUT /api/sales/update/:uuid` - Update sale
- `DELETE /api/sales/delete/:uuid` - Delete sale

### Dashboards
- `GET /api/dashboard/global-overview?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Global overview dashboard
- `GET /api/dashboard/provincial-analysis?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Provincial analysis
- `GET /api/dashboard/daily-monitor?date=YYYY-MM-DD` - Daily operations monitor
- `GET /api/dashboard/historical-trends?years=2025,2024&view_by=monthly` - Historical trends

### Years, Months & Weeks
- `GET /api/years/all` - Get all years
- `GET /api/months/all` - Get all months
- `GET /api/weeks/all` - Get all weeks

## Database Schema

The application uses the following main entities:
- **User**: User accounts with role-based permissions
- **Country**: Geographic countries
- **Province**: Geographic provinces/regions
- **Product**: Products being sold
- **Sale**: Sales transactions
- **Year**: Yearly sales targets
- **Month**: Monthly sales targets
- **Week**: Weekly sales targets
- **UserLogs**: Activity logging
- **Notification**: User notifications
- **PasswordReset**: Password reset tokens

## Dashboard Features

### Global Overview Dashboard
- Total sales and period-over-period comparison
- Provincial sales performance
- Sales trends with dynamic granularity (daily/weekly/monthly)
- Sales heatmap visualization

### Provincial Analysis Dashboard
- Multi-province comparison time series
- Contribution analysis (stacked area chart)
- Intra-day pattern analysis
- Target achievement tracking

### Daily Monitor Dashboard
- Real-time daily sales tracking
- Time slot analysis (8am, 12pm, 3pm, 8pm)
- Province entry status monitoring
- Cumulative sales charts

### Historical Trends Dashboard
- Year-over-year comparisons
- Cumulative yearly sales tracking
- Growth heatmaps
- Target achievement analysis

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o certikiosk
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SECRET_KEY` | JWT secret key | - |
| `PORT` | Server port | 8000 |
| `DB_HOST` | Database host | 127.0.0.1 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | certikioskdb |
| `DB_PORT` | Database port | 5432 |
| `SMTP_HOST` | SMTP server host | - |
| `SMTP_PORT` | SMTP server port | - |
| `SMTP_MAIL` | SMTP email address | - |
| `SMTP_PASSWORD` | SMTP password | - |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Author

Danny19977

## Support

For support, please open an issue in the GitHub repository.
