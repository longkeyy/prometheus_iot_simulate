package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	temperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_temperature_celsius",
			Help: "Simulated temperature of devices in Celsius.",
		},
		[]string{"province", "city", "district", "site", "device_id"},
	)
	humidityGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_humidity_percent",
			Help: "Simulated humidity of devices in percentage.",
		},
		[]string{"province", "city", "district", "site", "device_id"},
	)
)

func init() {
	// Register the custom metrics with Prometheus's default registry.
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(humidityGauge)
}

func simulateDeviceData(province, city, district, site, deviceID string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Simulate temperature between 15 and 35 degrees Celsius.
		temperature := 15 + rand.Float64()*20
		temperatureGauge.With(prometheus.Labels{
			"province":  province,
			"city":      city,
			"district":  district,
			"site":      site,
			"device_id": deviceID,
		}).Set(temperature)

		// Simulate humidity between 30% and 70%.
		humidity := 30 + rand.Float64()*40
		humidityGauge.With(prometheus.Labels{
			"province":  province,
			"city":      city,
			"district":  district,
			"site":      site,
			"device_id": deviceID,
		}).Set(humidity)

		time.Sleep(2 * time.Second)
	}
}

func main() {
	var wg sync.WaitGroup

	// Simulated data for various locations and devices.
	locations := []map[string]string{
		{"province": "Guangdong", "city": "Shenzhen", "district": "Nanshan", "site": "Site1", "device_id": "DeviceA"},
		{"province": "Guangdong", "city": "Shenzhen", "district": "Nanshan", "site": "Site1", "device_id": "DeviceB"},
		{"province": "Beijing", "city": "Beijing", "district": "Haidian", "site": "Site2", "device_id": "DeviceC"},
		{"province": "Shanghai", "city": "Shanghai", "district": "Pudong", "site": "Site3", "device_id": "DeviceD"},
		{"province": "Shanghai", "city": "Shanghai", "district": "Pudong", "site": "Site3", "device_id": "DeviceE"},
		// Add more locations and devices as needed
	}

	// Start a goroutine for each device to simulate data concurrently.
	for _, loc := range locations {
		wg.Add(1)
		go simulateDeviceData(loc["province"], loc["city"], loc["district"], loc["site"], loc["device_id"], &wg)
	}

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Serving metrics on http://localhost:10000/metrics")
	go func() {
		defer wg.Done()
		wg.Add(1)
		http.ListenAndServe(":10000", nil)
	}()

	// Wait for all goroutines to complete (in this case, they run indefinitely).
	wg.Wait()
}
