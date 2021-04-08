# Envoy XDS Server

This is a sample repo which demonstrates how to spin up an xDS Server for Envoy Proxy. 
This repo is used in Steve Sloka's fantastic video: [Build your own Envoy Control Plane](https://www.youtube.com/watch?v=qAuq4cKEG_E&t=837s).

The use of the code makes the most sense after watching the video.

## Config
The management server expects a config at ~config/config.yaml~.  This will be the file that drives the envoy configuration.  A simple example of a config would be:

``` yaml
# config/config.yaml

name: test_config
spec:
  listeners:
  - name: listener_0
    address: 0.0.0.0
    port: 9000
    routes:
    - name: echoroute
      prefix: /
      clusters:
      - echo
  clusters:
  - name: echo
    endpoints:
    - address: 0.0.0.0
      port: 9101
    - address: 0.0.0.0
      port: 9102
```

In this example, we'd have a listener for port `9000` that directs incoming requests to our echo cluster, which is load balanced across two endpoints.  The sample apps below would be the service at these endpoints.

## Sample Apps

Run some sample apps in docker to give some endpoints to route to:
```
docker run -d --rm --name=echo9100 -p 9100:8080 stevesloka/echo-server echo-server --echotext=Sample-Endpoint!
docker run -d --rm --name=echo9101 -p 9101:8080 stevesloka/echo-server echo-server --echotext=Sample-Endpoint!
docker run -d --rm --name=echo9102 -p 9102:8080 stevesloka/echo-server echo-server --echotext=Sample-Endpoint!
docker run -d --rm --name=echo9103 -p 9103:8080 stevesloka/echo-server echo-server --echotext=Sample-Endpoint!
docker run -d --rm --name=echo9104 -p 9104:8080 stevesloka/echo-server echo-server --echotext=Sample-Endpoint!
```

## Stop All Sample Apps

Stop all the sample endpoints created in the previous step:
```
docker stop echo9100
docker stop echo9101
docker stop echo9102
docker stop echo9103
docker stop echo9104
```

