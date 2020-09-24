#!/bin/bash
APP=tinyUrl
BasePath=$(dirname $(readlink -f $0))
mkdir -p $BasePath/bin/

CfgPath=$BasePath/etc/cfg.yml
DbPath=$BasePath/url.db

cd $BasePath

rm -rf $BasePath/bin/$APP

go build -o $BasePath/bin/$APP $BasePath/cmd/main.go

$BasePath/bin/$APP -c=$CfgPath -db=$DbPath