#### @copyright CrowdStrike, Inc. Copyright (c) 2024-.  All rights reserved.

# Crowdstrike Interview Take-Home Assignment
This is a take home task based off of real services we've built in the data engineering team, meant to be a realistic assessment
on how you might handle a similar task were you to be assigned one here.

## Sensor Data Streaming
## Background:
You have just been hired at Crowdstrike as a data engineer on the data science team (congratulations!), and you have been
tasked with creating an application that continuously reads a down sampled stream of data coming from the sensors. The data
being sent contains the byte contents of a file, the sha256 hash of the file, the platform the sensor is running on, the customer
id, and a timestamp. You must take this data and call the static file classifier API that classifies the byte content as
either malicious or benign. This data must be recorded in our static file classification table in our Cassandra database,
this data is used downstream by other processes to quickly look up classification information on hashes.

### TODO:
  1. Create a persistently running go application in the `kafka_event_file` directory we've provided that:
     1. Reads a message off of the Kafka topic.
     2. Deserializes the message.
     3. Makes an API request.
     4. Parses the response.
     5. Creates a record in the Cassandra database.
     6. Commits the offset.
  2. Fill out `Dockerfile` for building your application's docker image.
  3. Fill out the TODO section of `Makefile` that contains the docker build command to build the application.
  4. Fill out the TODO section in the docker
  5. Add any dependencies you might need to the `go.mod` file and run `go mod tidy && go mod vendor` to pull dependencies
  6. Once you think you have something working, run the command: `make run` to orchestrate everything
  7. Optional: Create a README in the directory that you put  with a short write-up if you think it will help us understand
    your thought process.

### API Details:
API is hosted on `http://api:8080` in docker, can be accessed outside of docker at `localhost:8080`, endpoint is `/classify`
sample API request and response:
```
curl -v -X "POST"  "localhost:8080/classify" -d '{"data": "M2ZkYSAzMnJhc2QyM3JhIGRmYTMyIHNkYWdkc2Fncw==", "platform_type": "PLATFORM_OSX"}'

{"classification":"benign","score":0.12135594786656093}
```


### Cassandra Table Schema:
```
CREATE TABLE IF NOT EXISTS cs.classification_results (
    sha256 ASCII,
    maliciousness_score FLOAT,
    classification ASCII,
    ts TIMESTAMP,
    PRIMARY KEY ((sha256), ts)
);
```

the hostname within docker is `cassandra:9094`

### Kafka Details:
1. The kafka topic is `cs.sensor_events`
2. Kafka is running on this host/port in docker: `kafka:9092`

### Requirements:
1. Must be written in golang.
2. We must be able to build and run this application locally.
3. The data written to the Cassandra table passes our verification.

### Tips:
* Feel free to use anything included in this zip file as examples to work off of.
* Use any and all internet resources to help you, such as:
  * github
  * go packages documentation
  * stack overflow
  * chatgpt
  * etc
* Some helpful links to get you started:
  * https://github.com/IBM/sarama
  * https://github.com/gocql/gocql
  * https://pkg.go.dev/net/http
* That being said, don't just blindly copy-paste. Try to understand everything you're writing because you'll likely have
to explain how and why you did things in a follow-up.
* Show off! Write code that you would consider professional level, really put your best foot forward.
* This may sound cliché, but try to have fun with it. If you see this as a boring chore, you won't enjoy having to do similar
tasks all day here full time. 