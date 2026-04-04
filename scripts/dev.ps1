# check dependencies
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "go is required but not installed. See: https://golang.org/dl/"
    exit 1
}
if (-not (Get-Command ngrok -ErrorAction SilentlyContinue)) {
    Write-Host "ngrok is required but not installed. See: https://ngrok.com/download"
    exit 1
}

# safe cleanup
$ErrorActionPreference = "Stop"

$go    = (Get-Command go).Source
$ngrok = (Get-Command ngrok).Source

# define output binary in a safe folder
$binFolder = "C:\dev\webhook-tester-bin"
if (-not (Test-Path $binFolder)) { New-Item -ItemType Directory -Path $binFolder | Out-Null }
$exePath = Join-Path $binFolder "server.exe"

# build the server explicitly (avoids temp exe in %AppData%)
Write-Host "Building server..."
& $go build -o $exePath "server/main.go"

Write-Host "Starting server..."
$server = Start-Process -FilePath $exePath -NoNewWindow -PassThru

Start-Sleep -Seconds 2

Write-Host "Starting ngrok..."
$ngrokProc = Start-Process -FilePath $ngrok -ArgumentList "http 8080" -NoNewWindow -PassThru

# retry loop to get ngrok URL
$url = $null
for ($i = 0; $i -lt 10; $i++) {
    try {
        $response = Invoke-RestMethod http://127.0.0.1:4040/api/tunnels
        $url = $response.tunnels[0].public_url
        if ($url) { break }
    } catch {}
    Start-Sleep -Milliseconds 500
}

if (-not $url) {
    Write-Host "Failed to get ngrok URL. Make sure ngrok started correctly."
    Stop-Process $server.Id -Force -ErrorAction SilentlyContinue
    Stop-Process $ngrokProc.Id -Force -ErrorAction SilentlyContinue
    exit 1
}

Write-Host ""
Write-Host "Public URL:"
Write-Host $url
Write-Host ""
Write-Host "Webhook endpoint:"
Write-Host "$url/hooks/test"
Write-Host ""
Write-Host "Press Ctrl+C to stop."

# cleanup on exit
try {
    Wait-Process $ngrokProc.Id
} finally {
    Stop-Process $server.Id -Force -ErrorAction SilentlyContinue
    Stop-Process $ngrokProc.Id -Force -ErrorAction SilentlyContinue
}