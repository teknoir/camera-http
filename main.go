package main

import (
	//"encoding/json"
	"encoding/base64"
	"fmt"
	"os/signal"
	"strconv"
	//"strings"

	"flag"
	"log"
	"os"
	"syscall"
	"time"
	"net/http"
	"io/ioutil"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	//"github.com/google/uuid"
	"github.com/bobziuchkovski/digest"
)

const (
	qosAtMostOnce byte = iota
	qosAtLeastOnce
	qosExactlyOnce
)

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var (
	mqttBroker = struct {
		host *string
		port *string
	}{
		flag.String("mqtt_broker_host", getEnv("MQTT_SERVICE_HOST", "mqtt.kube-system"), "MQTT Broker Host"),
		flag.String("mqtt_broker_port", getEnv("MQTT_SERVICE_PORT", "1883"), "MQTT Broker Port"),
	}
	defaultInterval, _ = strconv.Atoi(getEnv("UPDATE_INTERVAL", "1000")) // Milliseconds
	interval           = flag.Int("interval", defaultInterval, "Milliseconds between update")

	topic = flag.String("topic", getEnv("MQTT_OUT_0", "camera/images"), "The MQTT topic to publish images to")

	baseUrl = flag.String("base_url", getEnv("BASE_URL", "http://localhost/capture"), "The url to the camera image to capture")
	user = flag.String("user", getEnv("USERNAME", "root"), "The basic auth username")
	pass = flag.String("password", getEnv("PASSWORD", "teknoir"), "The basic auth password")
	deviceId = flag.String("device_id", getEnv("DEVICE_ID", "device-id-001"), "The device id of the camera")
)

type Payload struct {
	Image    string `json:"image"`
	DeviceID string `json:"device_id"`
	ImageId  string `json:"image_id"`
}

func captureImage(url string) (string, error) {
	//client := http.Client{
	//	Timeout: 5 * time.Second,
	//}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("[camera] Error preparing request", err.Error())
	}
	//req.SetBasicAuth(*user, *pass)
	//resp, err := client.Do(req)
	////resp, err := client.Get(url)
	t := digest.NewTransport(*user, *pass)
	resp, err := t.RoundTrip(req)
	if err != nil {
		log.Fatalln("[camera] Error getting image", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("[camera] Non-OK HTTP status", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[camera] Error reading image", err.Error())
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	//log.Println("[camera] Data: ", encoded)

	// return MIME base64 encoded jpeg
	return "data:image/jpeg;base64," + encoded, nil
}

func processCameras(mqttClient mqtt.Client, ticker *time.Ticker) {
	for t := range ticker.C {
		log.Println("[camera] at:", t)

		img, err := captureImage(*baseUrl)
		if err != nil {
			log.Println("[camera] Error capturing image", err.Error())
		} else {

			//id := uuid.New()
			//var payload = Payload{
			//	Image:    img,
			//	ImageId:  id.String(),
			//	DeviceID: *deviceId,
			//}

			//payloadStr, err := json.Marshal(payload)
			//if err != nil {
			//	log.Println("[camera] Error:", err)
			//}

			//token := mqttClient.Publish(*topic, qosAtLeastOnce, false, payloadStr)
			token := mqttClient.Publish(*topic, qosAtLeastOnce, false, img)
			if token.Wait() && token.Error() != nil {
				log.Println("[mqtt] Error:", token.Error())
			}
			//log.Println("[mqtt] Published image id: " + payload.ImageId + " from: " + payload.DeviceID)
			log.Println("[mqtt] Published image")
		}
	}
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	log.Println("[main] Entered")

	log.Println("[main] Flags")
	flag.Parse()

	// DEMUX
	mqttBroker := fmt.Sprintf("tcp://%v:%v", *mqttBroker.host, *mqttBroker.port)
	mqttOpts := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID("camera").SetCleanSession(true).SetAutoReconnect(true).SetMaxReconnectInterval(10 * time.Second)
	mqttClient := mqtt.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	} else {
		log.Println("[mqtt] Connected to server")
	}

	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)

	ticker := time.NewTicker(time.Duration(*interval) * time.Millisecond)
	defer ticker.Stop()
	go processCameras(mqttClient, ticker)

	signalHandler()
}

func signalHandler() {
	ch := make(chan os.Signal, 0)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-ch:
		log.Println("[main] signal received:", sig)
	}
	os.Exit(0)
}
