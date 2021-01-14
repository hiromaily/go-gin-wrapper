#!/bin/bash
# $1 retry count
# $2 command
# e.g. ./scripts/retry.sh 5 echo test1 test2

count=${1:-5}
if [ $# lt = 2 ]; then
    exit 1
fi
shift
cmd="$@"

counter=1
while [ $counter -le $count ]
do
  $cmd
  if [ $? -eq 0 ]; then
    exit 0
  fi
  ((counter++))
done
exit 1