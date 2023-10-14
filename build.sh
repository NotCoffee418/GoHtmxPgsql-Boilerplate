#!/bin/bash
# Array of OS and architectures
os_archs=("windows/amd64" "linux/amd64" "darwin/amd64")
app_name="GoWebsite-Boilerplate"

# Build releases
mkdir -p bin/release
for os_arch in "${os_archs[@]}"; do
  # Split the string into OS and ARCH
  IFS="/" read -r -a split <<< "$os_arch"
  os=${split[0]}
  arch=${split[1]}
  
  # Determine file extension
  ext=""
  if [ "$os" == "windows" ]; then
    ext=".exe"
  fi

  echo "Building for $os/$arch"

  # Build the Go binary
  env GOOS=$os GOARCH=$arch go build -o "bin/release/${app_name}_${os}_${arch}${ext}"
done
echo "Build complete"
