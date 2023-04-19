$version = $args[0]

$dir = "./build/v$version"

$linuxArm64Dir = $dir + "/nkc-reverse-proxy-linux-arm64-v" + $version
$linuxAmd64Dir = $dir + "/nkc-reverse-proxy-linux-amd64-v" + $version
$windowsArm64Dir = $dir + "/nkc-reverse-proxy-windows-arm64-v" + $version
$windowsAmd64Dir = $dir + "/nkc-reverse-proxy-windows-amd64-v" + $version


$linuxArm64Path = $linuxArm64Dir + "/nkc-reverse-proxy"
$linuxAmd64Path = $linuxAmd64Dir + "/nkc-reverse-proxy"
$windowsArm64Path = $windowsArm64Dir + "/nkc-reverse-proxy.exe"
$windowsAmd64Path = $windowsAmd64Dir + "/nkc-reverse-proxy.exe"

echo $linuxArm64Path
$env:GOOS="linux"
$env:GOARCH="arm64"
go build -o $linuxArm64Path .

echo $linuxAmd64Path
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o $linuxAmd64Path .

echo $windowsArm64Path
$env:GOOS="windows"
$env:GOARCH="arm64"
go build -o $windowsArm64Path .

echo $windowsAmd64Path
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o $windowsAmd64Path .

cp ./configs.template.yaml $linuxArm64Dir
cp ./configs.template.yaml $linuxAmd64Dir
cp ./configs.template.yaml $windowsAmd64Dir
cp ./configs.template.yaml $windowsArm64Dir

echo "Done"