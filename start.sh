#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
echo $basepath
cd $basepath

mkdir -p $basepath/logs
nohup $basepath/tg.robot > $basepath/logs/daemon.log 2>&1 & 
