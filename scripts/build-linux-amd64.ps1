$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o ../build/app-linux-amd64 ../.
echo "Build successfully"
echo "Output: build/app-linux-amd64"