# monGO by Cyrip

# to start Darth-Veda
./start.sh

# some details

## start with mongodb backend

docker-compose -f app-with-mongo-backend.yml up
### must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

### 3/1
All done in 1.01 seconds total 
Last datapoint is:
  time=511.55, expectedRequests=59606, scheduledRequests=59606, startedRequests=59606, successfulRequests=59606, validResponses=59606, failedRequests=0, requestsPerSecond=0, latency90th=52
  100.00% expected requests scheduled
  100.00% scheduled requests started
  100.00% started requests finished successful
  100.00% finished successful requests valid
Latency buckets are:
  <5ms,29269
  <10ms,8483
  <25ms,10030
  <50ms,5786
  <100ms,4659
  <250ms,1367
  <500ms,12
  <1000ms,0
  <10000ms,0

# thse solutions are not perfect - either :)

## start with elastic search cluster

sudo sysctl -w vm.max_map_count=262144 # it must be set on the docker host
docker-compose -f app-with-elastic-cluster-backend.yml
### must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080

All done in 1.04 seconds total 
Last datapoint is:
  time=555.44, expectedRequests=59606, scheduledRequests=59602, startedRequests=59602, successfulRequests=59602, validResponses=59596, failedRequests=0, requestsPerSecond=0, latency90th=61
  99.99% expected requests scheduled
  100.00% scheduled requests started
  100.00% started requests finished successful
  99.99% finished successful requests valid
Latency buckets are:
  <5ms,18510
  <10ms,10754
  <25ms,14294
  <50ms,8736
  <100ms,3612
  <250ms,2572
  <500ms,1107
  <1000ms,17
  <10000ms,0

## start with elastic search backend

docker-compose -f app-with-elastic-backend.yml 
### must wait until the backend and the app start
### this solution is not 100% perfect
app_1   | [GIN-debug] Listening and serving HTTP on :8080

## start with mongodb sharded cluster, it is not finished
docker-compose -f app-with-sharded-mongo-backend.yml up
### must wait until the backend and the app start
app_1   | [GIN-debug] Listening and serving HTTP on :8080
