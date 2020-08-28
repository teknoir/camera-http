# Teknoir http camera App
The camera app collect images from cameras over http.

## Build
```bash
gcloud builds submit . --config=cloudbuild.yaml --timeout=3600
```

## Build locally
```bash
go build -o camera -a .
```

## Run locally
```bash
docker run -it -p 1883:1883 -p 9001:9001 eclipse-mosquitto
./camera -mqtt_broker_host=localhost
```


## Legacy build and publish docker images
```bash
docker build -t tekn0ir/camera-http:arm64v8 -f arm64v8.Dockerfile .
docker push tekn0ir/camera-http:arm64v8
```

## Run on device
```bash
sudo kubectl run camera -ti --rm --image tekn0ir/camera-http:arm64v8 --generator=run-pod/v1 --overrides='{"spec":{"containers":[{"image":"tekn0ir/camera-http:arm64v8","name":"camera","command":["/bin/bash"],"tty":true,"stdin":true,"imagePullPolicy":"Always","env":[{"name":"BASE_URL","value":"http://192.168.3.200/axis-cgi/jpg/image.cgi?resolution=800x600"}]}]}}'
```

