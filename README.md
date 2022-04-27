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

Or if you have exposed the mqtt roker on a device:
```bash
./camera -mqtt_broker_host=jetsonnano-b00.local -mqtt_broker_port=31883 -base_url='http://192.168.1.137/axis-cgi/jpg/image.cgi?resolution=800x600' -auth_type=digest -user=root -password=Teknoir1
./camera -mqtt_broker_host=jetsonnano-b00.local -mqtt_broker_port=31883 -base_url='http://192.168.1.164/snap.jpg?JpegSize=M' -auth_type=digest -user=service -password='Teknoir1!'
```

## Legacy build and publish docker images
```bash
docker build -t tekn0ir/camera-http:arm64v8 -f arm64v8.Dockerfile .
docker push tekn0ir/camera-http:arm64v8
```

## Run on device
```bash
sudo kubectl run camera -ti --rm --image tekn0ir/camera-http:arm64v8 --generator=run-pod/v1 --overrides='{"spec":{"containers":[{"image":"tekn0ir/camera-http:arm64v8","name":"camera","command":["/bin/bash"],"tty":true,"stdin":true,"imagePullPolicy":"Always","env":[{"name":"BASE_URL","value":"http://192.168.3.101/axis-cgi/jpg/image.cgi?resolution=800x600"},{"name":"USERNAME","value":"root"},{"name":"PASSWORD","value":"teknoir"}]}]}}'
```

### Axis

#### Find camera
```
sudo apt-get install avahi-utils
avahi-browse -ltr _rtsp._tcp # Lists all RTSP cameras on the network
```
The RTSP URL wont work with this app.
This app only takes a still image URL.

#### Create URL
To create a URL for an Axis camera, use the IP and set a RESOLUTION(ex. 800x600) small enough GCP IoT Core max size is 250Kb.
```
http://<IP>/axis-cgi/jpg/image.cgi?resolution=800x600
```

### Bosch

#### Create URL
To create a URL for a Bosch camera, use the IP and choose the most appropriate size M (small enough GCP IoT Core max size is 250Kb).

[Documentation says](http://resource.boschsecurity.com/documents/Configuration_Note_enUS_1552286731.pdf): 

| T-shirt Size     | Size in documentation       | 5000 HD actual size |
| :---             | :----:                      |                ---: |
| S (small)        | 176 × 144/120 pixels (QCIF) | 256 × 144 pixels    |
| M (medium)       | 352 × 288/240 pixels (CIF)  | 512 × 528 pixels    |
| L (large)        | 704 × 288/240 pixels (2CIF) | 1280 × 720 pixels   |
| XL (extra large) | 704 × 576/480 pixels (4CIF) | 1920 × 1080 pixels  |
___I have added the atual sizes I got from my Bosch Dinion 5000 HD (NBN-50022-C)___

The URL looks like this with T-Shirt Size medium:
```
http://192.168.1.164/snap.jpg?JpegSize=M
```
