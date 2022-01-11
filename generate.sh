#!/bin/sh
git pull
docker rm -f devtipmebackend
docker rmi devtipmebackend
docker build -t devtipmebackend .
docker run -p 5000:5000 --name devtipmebackend -d devtipmebackend