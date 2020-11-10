@REM go get github.com/akavel/rsrc
rsrc -manifest inet.manifest -o app.syso -ico icon/icon.ico


go build -o build/inet.exe -ldflags "-H=windowsgui" .