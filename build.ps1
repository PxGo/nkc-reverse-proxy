param(
    [string]$version = ""
)

# Check if version number is specified
if ($version -eq "") {
    Write-Host "Version number is not specified and will be read from the version.txt"
    # Read the version number from the version.txt file in the project root directory
    $versionFilePath = Join-Path $PSScriptRoot "version.txt"
    if (Test-Path $versionFilePath) {
        $version = Get-Content $versionFilePath -Raw
        Write-Host "The version number is: $version"
    } else {
        Write-Host "Cannot find the version.txt file, please specify the version number parameter."
        return
    }
}

$dir = "./build/$version"

$linuxArm64Dir = $dir + "/nrp-linux-arm64-" + $version
$linuxAmd64Dir = $dir + "/nrp-linux-amd64-" + $version
$windowsArm64Dir = $dir + "/nrp-windows-arm64-" + $version
$windowsAmd64Dir = $dir + "/nrp-windows-amd64-" + $version

$linuxArm64Path = $linuxArm64Dir + "/nrp-linux-arm64-" + $version
$linuxAmd64Path = $linuxAmd64Dir + "/nrp-linux-amd64-" + $version
$windowsArm64Path = $windowsArm64Dir + "/nrp-windows-arm64-" + $version + ".exe"
$windowsAmd64Path = $windowsAmd64Dir + "/nrp-windows-amd64-" + $version + ".exe"

Write-Host $linuxArm64Path
$env:GOOS="linux"
$env:GOARCH="arm64"
go build -o $linuxArm64Path .

Write-Host $linuxAmd64Path
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o $linuxAmd64Path .

Write-Host $windowsArm64Path
$env:GOOS="windows"
$env:GOARCH="arm64"
go build -o $windowsArm64Path .

Write-Host $windowsAmd64Path
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o $windowsAmd64Path .

cp ./config.yaml $linuxArm64Dir
cp ./config.yaml $linuxAmd64Dir
cp ./config.yaml $windowsAmd64Dir
cp ./config.yaml $windowsArm64Dir

Write-Host "Done"