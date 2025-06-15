package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"rpi-data-logger/config"
	"rpi-data-logger/database"
	"rpi-data-logger/handler"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	slog.Info("Starting environment data logger...")

	db, err := database.New(cfg.DBPath)
	if err != nil {
		slog.Error("Database setup failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	mqttBrokerURL := fmt.Sprintf("tcp://%s:%s", cfg.MQTTBroker, cfg.MQTTPort)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBrokerURL)
	opts.SetClientID("rpi-data-logger-go")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	opts.SetDefaultPublishHandler(handler.NewMessageHandler(db))

	opts.OnConnect = func(c mqtt.Client) {
		slog.Info("Connected to MQTT broker")
		topic := "env/sensor/+/data"
		if token := c.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			slog.Error("Failed to subscribe to topic", "topic", topic, "error", token.Error())
		} else {
			slog.Info("Successfully subscribed", "topic", topic)
		}
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		slog.Warn("MQTT connection lost", "error", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("Failed to connect to MQTT broker", "error", token.Error())
		os.Exit(1)
	}

	slog.Info("Application is running. Press Ctrl+C to exit.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down...")
	client.Disconnect(250)
	slog.Info("Application exited.")
}
