$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o ../build/app-windows-amd64.exe ../.
echo "Build successfully"
echo "Output: build/app-windows-amd64.exe"