# Railway Deployment Guide

## Environment Variables to Set in Railway

Go to your Railway project → Variables tab and add:

### Required Database Variables
```
DB_HOST=your-railway-postgres-host
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-db-password
DB_NAME=certikioskdb
```

### Application Variables
```
PORT=8080
SECRET_KEY=your-secret-key
```

### CORS Configuration (Optional)
```
ALLOWED_ORIGINS=https://certikiosk-production.up.railway.app,http://localhost:3000
```
If not set, defaults to: `http://localhost:3000,http://192.168.0.70:3000,https://certikiosk-production.up.railway.app:3000,https://certikiosk-production.up.railway.app`

### Email Configuration (Optional)
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_MAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

### Google Drive API (Optional)
```
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=your-redirect-url
GOOGLE_TOKEN_FILE=token.json
```

## Deploy to Railway

1. **Connect your GitHub repository** to Railway
2. **Set environment variables** as listed above
3. **Deploy**: Railway will auto-detect the Dockerfile and build

## Health Check

Test the deployment:
```bash
curl https://certikiosk-api-production.up.railway.app/health
```

Expected response:
```json
{"status":"ok","service":"certikiosk-api"}
```

## Test CORS

```bash
curl -i -X OPTIONS "https://certikiosk-api-production.up.railway.app/api/auth/login" \
  -H "Origin: https://certikiosk-production.up.railway.app" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type, Authorization"
```

Look for `Access-Control-Allow-Origin: https://certikiosk-production.up.railway.app` in the response headers.

## Troubleshooting

### CORS Issues
- Check Railway logs for: `[info] CORS allowed origins: ...`
- Verify the origin in the log matches your frontend domain exactly
- Ensure no trailing slashes or port mismatches

### Database Connection Issues
- Check Railway logs for: `[info] connecting to database: ...`
- Verify all DB_* environment variables are set correctly
- Railway Postgres uses internal networking; use the internal connection string

### Build Issues
- Ensure `go.mod` and `go.sum` are committed
- Check Railway build logs for errors
- Verify all dependencies are available
