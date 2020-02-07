#!/bin/bash
a=`ps -ef|grep tinyUrl|grep -v grep|wc -l`
if [[ $a == 0 ]]
then
  cd /www/wwwroot/aid.im
  nohup ./tinyUrl &
fi
