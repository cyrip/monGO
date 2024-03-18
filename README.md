# monGO by Cyrip

# start with mongodb backend
docker-compose -f app-with-mongo-backend.yml up
# must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

# start with elastic search backend
docker-compose -f app-with-elastic-backend.yml 
# must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

# start with mongodb sharded cluster
docker-compose -f app-with-sharded-mongo-backend.yml up
# must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

# start with elastic search cluster
sudo sysctl -w vm.max_map_count=262144 # it must be set on the docker host
docker-compose -f app-with-elastic-cluster-backend.yml
# must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

