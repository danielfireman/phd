clear

PORT=8000

set -x

mvn clean package && \
docker build -t restserver . && \
docker run -d --name restserver -p $PORT:$PORT restserver
