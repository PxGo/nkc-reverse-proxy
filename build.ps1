echo "Building windows amd64 ..."

$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o build/app-windows-amd64.exe .

echo "Building windows arm64 ..."
$env:GOOS="windows"
$env:GOARCH="arm64"
go build -o build/app-windows-arm64.exe .

echo "Building linux amd64 ..."
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o build/app-linux-amd64 .

echo "Building linux arm64 ..."
$env:GOOS="linux"
$env:GOARCH="arm64"
go build -o build/app-linux-arm64 .


echo "Build successfully"
echo "Output: build/app-windows-amd64.exe"
echo "Output: build/app-windows-arm64.exe"
echo "Output: build/app-linux-amd64"
echo "Output: build/app-linux-arm64"