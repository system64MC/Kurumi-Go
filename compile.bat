go build -o "Kurumi.exe" -trimpath -tags static -ldflags="-s -w -H=windowsgui" main.go
upx "Kurumi.exe"