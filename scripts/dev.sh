#!/usr/bin/env bash

set -e

# check dependencies
command -v jq >/dev/null 2>&1 || { echo "jq is required but not installed. Run: sudo apt install jq / brew install jq"; exit 1; }
command -v ngrok >/dev/null 2>&1 || { echo "ngrok is required but not installed. See: https://ngrok.com/download"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "go is required but not installed. See: https://golang.org/dl/"; exit 1; }

echo "Starting server..."
go run server/main.go &
SERVER_PID=$!

# cleanup on exit (Ctrl+C or crash)
trap "echo 'Shutting down...'; kill $SERVER_PID $NGROK_PID 2>/dev/null; exit" EXIT INT TERM

echo "Starting ngrok..."
ngrok http 8080 > /tmp/ngrok.log &
NGROK_PID=$!

# retry loop to get ngrok URL
URL=""
for i in $(seq 1 10); do
  URL=$(curl -s http://127.0.0.1:4040/api/tunnels | jq -r '.tunnels[0].public_url')
  if [ -n "$URL" ] && [ "$URL" != "null" ]; then
    break
  fi
  sleep 0.5
done

if [ -z "$URL" ] || [ "$URL" = "null" ]; then
  echo "Failed to get ngrok URL. Check /tmp/ngrok.log for details."
  exit 1
fi

echo ""
echo "Public URL:"
echo "$URL"
echo ""
echo "Webhook endpoint:"
echo "$URL/hooks/test"
echo ""
echo "Press Ctrl+C to stop."

wait $SERVER_PID