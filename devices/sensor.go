package devices

const (
	// MemoryDevice free memory
	MemoryDevice = "memory"
	// CPUDevice idle cpu
	CPUDevice = "cpu"
	// HeartBeatDevice ...
	HeartBeatDevice = "heartbeat"
	// TemperatureDevice ...
	TemperatureDevice = "temperature"
	// TemperatureDHT11Device ...
	TemperatureDHT11Device = "temperature_dht11"
	// HumidityDevice ...
	HumidityDevice = "humidity"
	// GPSDevice ...
	GPSDevice = "gps"
)

// Device ...
type Device interface {
	Start()
}
