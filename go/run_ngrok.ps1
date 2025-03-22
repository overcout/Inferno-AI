
Start-Process -NoNewWindow -FilePath "ngrok" -ArgumentList "http 8080"

Start-Sleep -Seconds 2

$response = Invoke-RestMethod -Uri http://localhost:4040/api/tunnels
$httpsUrl = ($response.tunnels | Where-Object { $_.proto -eq "https" }).public_url

(Get-Content .env) -replace '^OAUTH_PUBLIC_URL=.*', "OAUTH_PUBLIC_URL=$httpsUrl" |
  ForEach-Object { $_ -replace '^OAUTH_REDIRECT_URL=.*', "OAUTH_REDIRECT_URL=$httpsUrl/oauth/callback" } |
  Set-Content .env

Write-Host "✅ Updated .env:"
Write-Host "OAUTH_PUBLIC_URL=$httpsUrl"
Write-Host "OAUTH_REDIRECT_URL=$httpsUrl/oauth/callback"
Write-Host ""
