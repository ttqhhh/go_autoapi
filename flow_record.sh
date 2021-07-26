#!/bin/bash

function record_unit() {
  # timeout 300 ./gor --input-raw :8003 --input-raw-track-response --output-stdout --output-file=rankingtest.gor --http-allow-url /recommend/ranking

  cd /data/ranking-data/zy_ps_10/tmp

  service_name=`printenv | grep 'K8S_POD_NAME=online' | awk -F '-' '{print $2}'`

  # The default port is 8003, don't edit.
  port=8003
  # The default url is "/recommend/ranking", don't   edit.
  url="/recommend/ranking"
  # record 5min
  timeout 300 gor --input-raw :${port} --input-raw-track-response --output-stdout --output-file=${service_name}.gor --http-allow-url ${url}

  file_name=${service_name}_0.gor
  echo "The service name is "${service_name}
  echo "The Flow file name in request is "${file_name}

  if [ -f ${file_name} ]; then
      curl -F file=@$file_name http://172.16.2.86:8080/flowreplay/collect_flow_file
      # Delete the flow file
      rm  -rf ${file_name}
  else
    echo "The ${file_name} is not exist!!!"
  fi
}
# The script need a param, be used in the loop. eg: ./gor_record 720
while :
do
    period=$1
    record_unit
    sleep ${period}h
done