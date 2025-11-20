# Emergency CORS Fix for Railway

If Railway deployment is stuck, set this environment variable temporarily:

## Option 1: Allow Specific Origin (Recommended)
In Railway Variables tab:
```
ALLOWED_ORIGINS=https://certikiosk-production.up.railway.app
```

## Option 2: Allow All Origins (TEMPORARY ONLY - Security Risk)
```
ALLOWED_ORIGINS=*
```
⚠️ Only use this for testing, then immediately change to Option 1

## Verify After Setting
1. Restart/redeploy the Railway service
2. Test with:
```bash
curl -i -X OPTIONS "https://certikiosk-api-production.up.railway.app/api/auth/login" \
  -H "Origin: https://certikiosk-production.up.railway.app" \
  -H "Access-Control-Request-Method: POST"
```

3. Look for these headers in response:
```
Access-Control-Allow-Origin: https://certikiosk-production.up.railway.app
Access-Control-Allow-Credentials: true
```

## Why This Happens
- Railway may cache old builds
- Auto-deploy might be disabled
- Build might have failed silently
- Environment variables weren't reloaded

## Check Railway Logs
Look for this line in logs:
```
[info] CORS allowed origins: https://certikiosk-production.up.railway.app
```

If you see old origins or the line is missing, the new code isn't running yet.
