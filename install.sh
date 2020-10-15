#!/bin/bash

echo "installing quire cli..."

OS=""
VERSION="0.2"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="darwin"
elif [[ "$OSTYPE" == "cygwin" ]]; then
        OS="windows"
elif [[ "$OSTYPE" == "msys" ]]; then
        OS="windows"
elif [[ "$OSTYPE" == "win32" ]]; then
        OS="windows"
fi

echo " > running the installation for quire cli v${VERSION} on OS ${OS}"

echo " > fetching the files"
curl -L "https://github.com/AienTech/quire-cli/releases/download/v${VERSION}/quire-cli-v${VERSION}-${OS}-amd64.tar.gz" > quire-cli.tar.gz
echo " > untaring the files"
tar -zxvf quire-cli.tar.gz
echo " > moving the file to /usr/local/bin"
mv quire-cli /usr/local/bin/quire
rm quire-cli.tar.gz

echo "installation was successful"
echo "you can use the command 'quire' to use the quire cli"
