#/bin/bash

for i in $(docker ps -a|grep -v CONTAINER| awk '{ print $1;}');do docker rm $i;done
for i in $(docker volume ls|grep -v VOLUME| awk '{ print $2;}');do docker volume rm $i;done
for i in $(docker images|grep darth-| awk '{ print $1;}');do docker rmi $i;done
