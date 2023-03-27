go build -o "Kurumi.exe" -tags static -ldflags="-s -w -H=windowsgui" main.go
upx "Kurumi.exe"