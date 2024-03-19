#!/bin/bash

docker run -v "$(pwd):/app" --network host openjdk:23-slim-bullseye java -jar /app/stress-1.0.0.jar -v
