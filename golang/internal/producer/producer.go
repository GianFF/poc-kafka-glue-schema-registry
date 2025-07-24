package producer

import (
	"fmt"
	"kafka-glue-poc/internal/kafka"
	"kafka-glue-poc/internal/schema"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer    sarama.SyncProducer
	glueClient  *schema.GlueClient
	topic       string
}

type UserSignedUp struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email,omitempty"`
	Timestamp string `json:"timestamp"`
}

func NewProducer(topic string) (*Producer, error) {
	config := kafka.NewKafkaConfig()
	producer, err := sarama.NewSyncProducer(kafka.GetBrokers(), config)
	if err != nil {
		return nil, fmt.Errorf("error creando productor: %w", err)
	}

	return &Producer{
		producer:   producer,
		glueClient: schema.NewGlueClient(),
		topic:      topic,
	}, nil
}

func (p *Producer) Start() error {
	fmt.Printf("🚀 Productor iniciado - enviando mensajes al topic \"%s\" cada 2 segundos...\n", p.topic)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		user := p.generateRandomUser()
		
		if err := p.sendMessage(user); err != nil {
			fmt.Printf("❌ Mensaje rechazado - no cumple esquema:\n")
			fmt.Printf("   Usuario: %+v\n", user)
			fmt.Printf("   Error: %s\n", err.Error())
		}
	}

	return nil
}

func (p *Producer) sendMessage(user UserSignedUp) error {
	// 1. Obtener esquema desde Glue
	codec, err := p.glueClient.GetSchema("user-events", "UserSignedUp")
	if err != nil {
		return fmt.Errorf("error obteniendo esquema: %w", err)
	}

	// 2. VALIDACIÓN DEL ESQUEMA: validar antes del envío
	userMap := map[string]interface{}{
		"user_id":   user.UserID,
		"email":     user.Email,
		"timestamp": user.Timestamp,
	}

	// Validar y serializar con Avro
	avroData, err := codec.BinaryFromNative(nil, userMap)
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	// 3. Enviar mensaje a Kafka
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(avroData),
	}

	_, _, err = p.producer.SendMessage(message)
	if err != nil {
		fmt.Printf("💥 Error enviando mensaje a Kafka: %s\n", err.Error())
		return err
	}

	// 4. Log de éxito
	fmt.Printf("🚀 Mensaje enviado:\n")
	fmt.Printf("   Usuario: %s, Email: %s\n", user.UserID, user.Email)

	return nil
}

func (p *Producer) generateRandomUser() UserSignedUp {
	id := rand.Intn(100000)
	userID := fmt.Sprintf("user%d", id)
	timestamp := time.Now().Format(time.RFC3339)
	
	var email string
	if id%2 == 0 {
		// error de validación
	} else {
		email = fmt.Sprintf("user%d@example.com", id)
	}

	return UserSignedUp{
		UserID:    userID,
		Email:     email,
		Timestamp: timestamp,
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
