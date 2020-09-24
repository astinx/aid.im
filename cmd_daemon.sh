#!/bin/bash
APP=tinyUrl
BasePath=$(dirname $(readlink -f $0))
CfgPath=$BasePath/etc/cfg.yml
DbPath=$BasePath/url.db
cd $BasePath

a=`ps -ef|grep tinyUrl|grep -v grep|wc -l`
if [[ $a == 0 ]]
then
  nohup $BasePath/bin/$APP -c=$CfgPath -db=$DbPath &
fi