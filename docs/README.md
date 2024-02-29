Sensor Event Streaming using Apache Kafka
=========================================

JIRA Tracking
https://gorairahul2022.atlassian.net/browse/ENGGASSIGN-1

# Background

This is a kafka streaming application which reads data from the sensor stream and updates
the persistent database for other applications to use that data. 

This had 4 major components - 

1. Event Generator
2. Event Pipeline
3. API
4. Database


# Design

![kafka](https://github.com/rahulgorai22/kafka-streaming-project/assets/42316294/9348df29-0c9f-476b-9afe-b6238baf74fd)


# Running Application

```shell
# If you have an old docker images of event_generator, kafka_event_pipeline and api then delete the cached docker files.
rm docker-build-api docker-build-event-generator docker-build-kafka-event-pipeline
```


```shell
# To download and vendor the packages run 
go mod tidy && go mod vendor
```
To orchestrate everything execute ```make run```
```shell
# Build the docker image of kafka_event_pipeline
# Run the kafka producer and consumer
# Consumer consumes the messages, makes an API Call and updates the database.
make run

[+] Building 35.5s (12/12) FINISHED
 => [internal] load build definition from Dockerfile.kafka_event_pipeline                             0.0s
 => => transferring dockerfile: 143B                                                                  0.0s
 => [internal] load .dockerignore                                                                     0.0s
 => => transferring context: 2B                                                                       0.0s
 => [internal] load metadata for docker.io/library/golang:1.21.7-alpine3.19                           3.2s
 => [internal] load build context                                                                     0.2s
 => => transferring context: 120.34kB                                                                 0.2s
 => CACHED [builder 1/6] FROM docker.io/library/golang:1.21.7-alpine3.19@sha256:0ff68fa7b2177e8d68b4  0.0s
 => [builder 2/6] RUN go install golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616         14.2s
 => [builder 3/6] RUN mkdir /build                                                                    0.2s
 => [builder 4/6] ADD . /build/                                                                       0.3s
 => [builder 5/6] WORKDIR /build                                                                      0.0s
 => [builder 6/6] RUN export GOFLAGS=-mod=vendor     && go list ./... | grep -v vendor | xargs go v  17.2s
 => [deploy 1/2] COPY --from=builder /build/kafka-event-pipeline /                                    0.0s
 => exporting to image                                                                                0.1s
 => => exporting layers                                                                               0.0s
 => => writing image sha256:33755e0bd0b092a4c45904da36784f6a4dac2ba98706f643b850c5c7a253dba4          0.0s
 => => naming to docker.io/library/cs-streaming-take-home-task-kafka-event-pipeline                   0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
WARN[0000] Found orphan containers ([kafka_streaming_project-kafka_event_pipeline-1]) for this project. If you removed or renamed this service in your compose file, you can run this command with the --remove-orphans flag to clean it up.
[+] Running 7/7
 ⠿ Container kafka_streaming_project-cassandra-1        Created                                       0.0s
 ⠿ Container kafka_streaming_project-kafka-1            Created                                       0.0s
 ⠿ Container kafka_streaming_project-kafkasetup-1       Created                                       0.0s
 ⠿ Container kafka_streaming_project-api-1              Created                                       0.0s
 ⠿ Container kafka_streaming_project-cassandrasetup-1   Created                                       0.0s
 ⠿ Container kafka_streaming_project-pipeline-1         Recreated                                     0.1s
 ⠿ Container kafka_streaming_project-event_generator-1  Created  
 
 kafka_streaming_project-pipeline-1         | {"level":"info","ts":1709197642.774904,"caller":"kafka_event_pipeline/kafka_event_pipeline.go:54","msg":"Sarama consumer up and running !! {0x4000210300}"}
 kafka_streaming_project-pipeline-1         | {"level":"info","ts":1709197642.9582465,"caller":"kafka_event_pipeline/kafka_event_pipeline.go:105","msg":"Customer [45860326449]: Successfully updated database"}
kafka_streaming_project-pipeline-1         | {"level":"debug","ts":1709197642.9583123,"caller":"kafka_event_pipeline/kafka_event_pipeline.go:88","msg":"Message claimed: timestamp = 2024-02-29 09:07:22.783 +0000 UTC, topic = cs.sensor_events\n"}
kafka_streaming_project-pipeline-1         | {"level":"info","ts":1709197642.964055,"caller":"kafka_event_pipeline/kafka_event_pipeline.go:105","msg":"Customer [982126636067]: Successfully updated database"}
kafka_streaming_project-event_generator-1  | 2024/02/29 09:07:25 produced kafka message for sha256: 619d74b16cab1a3da214142527c81e03f05bb5350f6c608952e1684482b6148a, partition: 0 offset: 3261
kafka_streaming_project-pipeline-1         | {"level":"debug","ts":1709197645.416367,"caller":"kafka_event_pipeline/kafka_event_pipeline.go:88","msg":"Message claimed: timestamp = 2024-02-29 09:07:25.391 +0000 UTC, topic = cs.sensor_events\n"}
```

```shell
# List all the images created for the cs streaming project
docker images | grep cs
cs-streaming-take-home-task-kafka-event-pipeline                 latest              89668efa7db9   15 hours ago    17.3MB
cs-streaming-take-home-task-api                                  latest              2ade9d2ece25   46 hours ago    6.73MB
cs-streaming-take-home-task-event-generator                      latest              8ec95d1b721f   46 hours ago    13.7MB

# Check all the running containers
docker ps | grep cs
CONTAINER ID   IMAGE                                                     COMMAND                  CREATED        STATUS              PORTS                                                       NAMES
950663eddda9   cs-streaming-take-home-task-kafka-event-pipeline:latest   "/kafka-event-pipeli…"   15 hours ago   Up About a minute                                                               kafka_streaming_project-pipeline-1
ef12e14f2331   cs-streaming-take-home-task-event-generator:latest        "/event-generator"       46 hours ago   Up About a minute                                                               kafka_streaming_project-event_generator-1
c5c6eeb77c56   cs-streaming-take-home-task-api:latest                    "/api"                   46 hours ago   Up About a minute   0.0.0.0:8080->8080/tcp                                      kafka_streaming_project-api-1

# Other kafka and cassandra containers
f888fd36eb25   cassandra:latest                                          "docker-entrypoint.s…"   2 days ago     Up About a minute   7000-7001/tcp, 7199/tcp, 9160/tcp, 0.0.0.0:9042->9042/tcp   kafka_streaming_project-cassandra-1
d17012ca8a99   bitnami/kafka:latest                                      "/opt/bitnami/script…"   2 days ago     Up About a minute   0.0.0.0:9092->9092/tcp                                      kafka_streaming_project-kafka-1

# Enter the cassandra database and check the data
docker exec -it f888fd36eb25 bash
root@f888fd36eb25:/# cqlsh localhost -u cassandra -p cassandra
Warning: Using a password on the command line interface can be insecure.
Recommendation: use the credentials file to securely provide the password.
Connected to Test Cluster at localhost:9042
[cqlsh 6.1.0 | Cassandra 4.1.4 | CQL spec 3.4.6 | Native protocol v5]
Use HELP for help.
cassandra@cqlsh> SELECT * FROM cs.classification_results LIMIT 2;

@ Row 1
---------------------+------------------------------------------------------------------
 sha256              | 80b8f4425b082584909700d6c5dc98cd93da62e51c0734dba4afdfe882c617ac
 ts                  | 2024-02-27 13:50:05.000000+0000
 classification      | benign
 maliciousness_score | -1.5435e-10

@ Row 2
---------------------+------------------------------------------------------------------
 sha256              | 80b8f4425b082584909700d6c5dc98cd93da62e51c0734dba4afdfe882c617ac
 ts                  | 2024-02-27 18:10:30.000000+0000
 classification      | benign
 maliciousness_score | -1.5435e-10
```

