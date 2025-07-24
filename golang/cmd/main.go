package main

import (
	"context"
	"fmt"
	"kafka-glue-poc/internal/consumer"
	"kafka-glue-poc/internal/producer"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	topic   = "users.signedup"
	groupID = "email-service-group"
)

func main() {
	// Crear contexto para manejo de señales
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Manejo de señales para cierre limpio
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Iniciar consumidor
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		fmt.Println("🚀 Iniciando consumidor...")
		
		c, err := consumer.NewConsumer(topic, groupID)
		if err != nil {
			log.Fatalf("Error creando consumidor: %v", err)
		}
		defer c.Close()

		if err := c.Start(ctx); err != nil {
			log.Printf("Error en consumidor: %v", err)
		}
	}()

	// Esperar un poco antes de iniciar el productor
	time.Sleep(2 * time.Second)

	// Iniciar productor
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		fmt.Println("🚀 Iniciando productor...")
		
		p, err := producer.NewProducer(topic)
		if err != nil {
			log.Fatalf("Error creando productor: %v", err)
		}
		defer p.Close()

		if err := p.Start(); err != nil {
			log.Printf("Error en productor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-sigChan
	fmt.Println("\n🛑 Deteniendo aplicación...")
	
	cancel()
	wg.Wait()
	
	fmt.Println("✅ Aplicación detenida correctamente")
}
