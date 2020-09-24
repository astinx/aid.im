#!/bin/bash
APP=tinyUrl
BasePath=$(dirname $(readlink -f $0))
CfgPath=$BasePath/etc/cfg.yml
DbPath=$BasePath/url.db
cd $BasePath && \
$BasePath/bin/$APP -c=$CfgPath -db=$DbPath