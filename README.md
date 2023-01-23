## Overview

To make it easier to understand, let's take a simplified description of 
a cardiac health monitoring system.

### Diagram
Let's have a look at the diagram and move on to the practical part.

```mermaid
    C4Container
    title Container diagram for Wearable heart rate monitoring system

    Person(patient, Patient, "Patient with a heart health monitoring device")

    Container_Boundary(c1, "Wearable heart rate monitoring system") {
        Container(spa, "Single-Page App", "JavaScript, Angular", "Provides an interface to control the heart rate transmitter")
        Container(backend, "Backend", "Golang", "Contains all processing of data from the heart rate transmitter. Sends notifications in case of emergency")
        Container(push-sender, "Push Sender", "Golang", "Contains the logic for sending notifications")
    }
    
    System_Ext(email_system, "E-Mail System", "The external messaging system")
    System_Ext(voip_system, "VoIP System", "The external calling system")
    
    Rel(patient, spa, "Control the heart rate monitoring system")
    
    Rel(spa, backend, "Control the heart rate monitoring system")
    Rel(backend, push-sender, "Notify personal in case of emergency")
    Rel(push-sender, email_system, "Make email notification")
    Rel(push-sender, voip_system, "Make phone call")
```

### Usecase

The ```Wearable``` service monitors the patient's cardiac health
and notifies the staff in case of an emergency .

Let's pretend that we need to develop a ```Wearable``` service.

In that case, we probably want to disable sending messages to real
notification systems for the time of development and debugging.
To do this, we can run mock service instead of the ```push-sender``` service.

The ```grpc-wiremock``` will do this job just fine. Here is a step-by-step tutorial.

## Hands-on
All you need is ~~love~~ set of **openapi** or **proto** contracts.

### Contracts

#### Wearable

```protobuf
...

service WearableService {
  rpc BeatsPerMinute(BeatsPerMinuteRequest) 
      returns (stream BeatsPerMinuteResponse);
}
```

#### push-sender

```protobuf
...

service PushSender {
  rpc Notify(NotifyRequest) returns (NotifyResponse);
}
```

### Before run

To run ```grpc-wiremock```, you need to prepare three directories:

- path to the folder with the service-dependency contracts;
    ```
    deps
    └── services
        └── push-sender
            └── grpc
                └── push-sender.proto
    ```
- path to the directory for Wiremock config and mappings (can be empty);
  ```
  test/wiremock
  ```
- path to the directory for generated certificates. For example:
  ```
  /tmp/certs
  ```
  
### Run

Execute the folowing command in your console to start the ```grpc-wiremock```:
```bash
MOCKS_PATH="$(PWD)/test/wiremock"

CERTS_PATH="/tmp/certs"

CONTRACTS_PATH="deps"

docker run \
  -p 9000:9000 \
  -p 8000:8000 \
  -v ${MOCKS_PATH}:/home/mock \
  -v ${CERTS_PATH}:/etc/ssl/mock/share \
  -v ${CONTRACTS_PATH}:/proto \
  SberMarket-Tech/grpc-wiremock@latest
```

You can also run services using ```docker compose```. 
Check out [compose](build/docker-compose.yaml) file. 

```bash
make up
```

### After run

Let's see what happened after running ```grpc-wiremock```:
- you got generated certificates to test with secured connection:
  ```
  /tmp/certs/
  ├── mockCA.crt
  ├── mockCA.key
  └── mockCA.srl
  ```

- ```test/wiremock``` directory looks like this:
  ```
  test/wiremock
  ├── push-sender
  │   ├── __files
  │   └── mappings
  │       └── pushsender_notify_post_200.json
  └── services.json
  ```
  
- Wiremock config ```service.json```:
  ```json
  {
    "services" : [ {
      ...
      
      "name" : "push-sender",
      "port" : 8000,
    } ]
  }
  ```
- and generated mock for your ```Notify``` method from [push-sender.proto](deps/services/push-sender/grpc/push-sender.proto):
  ```json
  {
    ...
    
    "request" : {
      "urlPath" : "/PushSender/Notify",
      "method" : "POST"
    },
    "response" : {
      "status" : 200,
      "body" : "{\n  \"status\" : 425895108\n}",
      "headers" : {
        "Content-Type" : "application/json"
      }
    }
  }
  ```
  
### Wiremock APIs

As you can learn from the [docs](https://wiremock.org/studio/docs/stubbing), ```Wiremock Studio``` supports
multiple APIs simultaneously. 

This means that if you want to add more 
mocks besides ```push-sender```, ```grpc-wiremock``` will carefully create mocks 
for all dependencies in the ```deps``` directory for you.

Just to let you visualize it:

``` mermaid
  C4Container
  title Wiremock Studio APIs
    
  Container_Boundary(c1, "Wiremock") {
      Container(push-sender, "push-sender API", "", "Returns the mocs, listening on port 8000. API ports always start with 8000")
      Container(you-name-it, "You Name It API", "", "This could be your service. Port: 8001")
  }
```

## So what?

```Wiremock``` can make it easier to develop and test connected services.
By replacing the actual service API with a mock.

And with ```grpc-wiremock``` you got support for **proto** contracts.

You have the following:
- the mocks are generated automatically, but you can change them 
and ```grpc-wiremock``` will notice this and reload the mocks for you;
- you don't need to know which port Wiremock allocated for the ```push-sender```.

You can access the ```push-sender``` API like this:

#### from ```wearable``` container:
  - direct request to Wiremock API:
    ```bash
    curl -XPOST wearable-mock:8000/PushSender/Notify
    {
        "status" : 425895108
    }
    ```
  - request to reverse proxy server (need to specify name of dependency):
    ```bash
    curl -XPOST --header "Host: push-sender" wearable-mock/PushSender/Notify
    {
        "status" : 425895108
    }
    ```
  - of course, you can access gRPC handler:
    ```bash
    grpcurl \
        -d '{"uuid": "1234", "message": "foo"}' \
        --authority "push-sender" --plaintext \
        wearable-mock:3010 push_sender.PushSender/Notify 
    {
        "status": 425895108
    }
    ```

This is all well and good, but we absolutely need to 
get the status ```1126``` after sending the notification.
  
Let's change this value in [pushsender_notify_post_200.json](test/wiremock/push-sender/mappings/):
  ```json
  {
    ...
    
    "request" : {
      "urlPath" : "/PushSender/Notify",
      "method" : "POST"
    },
    "response" : {
      "status" : 200,
      "body" : "{\n  \"status\" : 1126\n}",
      "headers" : {
        "Content-Type" : "application/json"
      }
    }
  }
  ```

It's also possible to change the body of the mock 
using the GUI at http://localhost:9000.

How much fine-tuning of the mocks can be done, 
see [here](https://wiremock.org/studio/docs/advanced-stubbing).

Check the ```wearable-mock``` logs:
```bash
docker logs wearable-mock | grep wiremock.mappings
```
```
wiremock.mappings: Upload mappings for push-sender ... OK
```

As you may have guessed, the second 
request will give us our target status:

```bash
grpcurl \
    -d '{"uuid": "1234", "message": "foo"}' \
    --authority "push-sender" --plaintext \
    wearable-mock:3010 push_sender.PushSender/Notify 
{
    "status": 1126
}
```
