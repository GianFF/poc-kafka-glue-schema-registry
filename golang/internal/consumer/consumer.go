package consumer

import (
	"context"
	"fmt"
	"kafka-glue-poc/internal/kafka"
	"kafka-glue-poc/internal/schema"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	glueClient    *schema.GlueClient
	topic         string
	groupID       string
}

type ConsumerGroupHandler struct {
	glueClient *schema.GlueClient
}

func NewConsumer(topic, groupID string) (*Consumer, error) {
	config := kafka.NewKafkaConfig()
	
	consumerGroup, err := sarama.NewConsumerGroup(kafka.GetBrokers(), groupID, config)
	if err != nil {
		return nil, fmt.Errorf("error creando consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		glueClient:    schema.NewGlueClient(),
		topic:         topic,
		groupID:       groupID,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	fmt.Printf("🔍 Consumidor esperando mensajes en el topic \"%s\"...\n", c.topic)

	handler := &ConsumerGroupHandler{
		glueClient: c.glueClient,
	}

	topics := []string{c.topic}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.consumerGroup.Consume(ctx, topics, handler); err != nil {
				log.Printf("Error en consumer group: %v", err)
				return err
			}
		}
	}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.processMessage(message)
		session.MarkMessage(message, "")
	}
	return nil
}

func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) {
	messageInfo := fmt.Sprintf("Topic: %s, Partition: %d, Offset: %d", 
		message.Topic, message.Partition, message.Offset)

	// 1. Obtener esquema desde Glue
	codec, err := h.glueClient.GetSchema("user-events", "UserSignedUp")
	if err != nil {
		fmt.Printf("❌ Error procesando mensaje:\n")
		fmt.Printf("   %s\n", messageInfo)
		fmt.Printf("   💥 Error técnico: %s\n", err.Error())
		return
	}

	// 2. Validar y deserializar el mensaje
	native, _, err := codec.NativeFromBinary(message.Value)
	if err != nil {
		fmt.Printf("❌ Error procesando mensaje:\n")
		fmt.Printf("   %s\n", messageInfo)
		
		if strings.Contains(err.Error(), "invalid") {
			fmt.Printf("   🚫 Mensaje no cumple con el esquema: %s\n", err.Error())
		} else {
			fmt.Printf("   💥 Error técnico: %s\n", err.Error())
		}
		return
	}

	// 3. Log de éxito
	fmt.Printf("📨 Mensaje recibido:\n")
	fmt.Printf("   %s\n", messageInfo)
	fmt.Printf("   Datos: %+v\n", native)
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}
