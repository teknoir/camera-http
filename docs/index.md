---
id: index
title: User Documentation
# prettier-ignore
description: The camera app collect images from cameras over http/api.
---

# HTTP camera App
The camera app collect images from cameras over http/api.

![Overview diagram](./assets/diagram.svg){ width="800" }

## Features

The most simple app to ingest images into a message stream with a fixed interval. The app calls a http(s) API that 
provide a binary JPEG image response.

## Settings

| Var                  | Description                              | Default                  |
|----------------------|------------------------------------------|--------------------------|
| `MQTT_SERVICE_HOST`  | MQTT Broker Host                         | mqtt.kube-system         |
| `MQTT_SERVICE_PORT`  | MQTT Broker Port                         | 1883                     |
| `UPDATE_INTERVAL`    | Milliseconds between update              | 1000                     |
| `MQTT_OUT_0`         | The MQTT topic to publish images to      | camera/images            |
| `BASE_URL`           | The url to the camera image to capture   | http://localhost/capture |
| `AUTH_TYPE`          | The auth type can be `digest` or `basic` | digest                   |
| `USERNAME`           | The auth username                        | root                     |
| `PASSWORD`           | The auth password      Delete resource   | teknoir                  |

## Advanced

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: camera
spec:
  replicas: 1
  selector:
    matchLabels:
      app: camera
  template:
    metadata:
      labels:
        app: camera
    spec:
      containers:
        - name: camera
          image: gcr.io/teknoir/camera-http:latest-master
          imagePullPolicy: IfNotPresent
          env:
            - name: MQTT_SERVICE_HOST
              value: "mqtt.kube-system"
            - name: MQTT_SERVICE_PORT
              value: "1883"
            - name: MQTT_OUT_0
              value: "camera/images"
            - name: UPDATE_INTERVAL
              value: "5000"
            - name: BASE_URL
              value: "http://192.168.0.186/snap.jpg?JpegSize=L"
            - name: AUTH_TYPE
              value: "digest"
            - name: USERNAME
              value: "username"
            - name: PASSWORD
              value: "password"
```
