package schema

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/linkedin/goavro/v2"
)

type GlueClient struct {
	client     *glue.Glue
	schemaCache map[string]*goavro.Codec
	mutex      sync.RWMutex
}

func NewGlueClient() *GlueClient {
	// Configuración para LocalStack
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	// Configurar credenciales para LocalStack
	sess.Config.Credentials = credentials.NewStaticCredentials("test", "test", "")

	return &GlueClient{
		client:      glue.New(sess),
		schemaCache: make(map[string]*goavro.Codec),
	}
}

func (gc *GlueClient) GetSchema(registryName, schemaName string) (*goavro.Codec, error) {
	cacheKey := fmt.Sprintf("%s:%s", registryName, schemaName)
	
	// Verificar cache primero
	gc.mutex.RLock()
	if codec, exists := gc.schemaCache[cacheKey]; exists {
		gc.mutex.RUnlock()
		return codec, nil
	}
	gc.mutex.RUnlock()

	fmt.Println("📥 Obteniendo esquema desde Glue Schema Registry...")

	// Obtener esquema desde Glue
	input := &glue.GetSchemaVersionInput{
		SchemaId: &glue.SchemaId{
			RegistryName: aws.String(registryName),
			SchemaName:   aws.String(schemaName),
		},
		SchemaVersionNumber: &glue.SchemaVersionNumber{
			LatestVersion: aws.Bool(true),
		},
	}

	result, err := gc.client.GetSchemaVersion(input)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo esquema desde Glue: %w", err)
	}

	// Parsear el esquema Avro
	var schemaDefinition map[string]interface{}
	if err := json.Unmarshal([]byte(*result.SchemaDefinition), &schemaDefinition); err != nil {
		return nil, fmt.Errorf("error parseando esquema JSON: %w", err)
	}

	codec, err := goavro.NewCodec(*result.SchemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("error creando codec Avro: %w", err)
	}

	// Guardar en cache
	gc.mutex.Lock()
	gc.schemaCache[cacheKey] = codec
	gc.mutex.Unlock()

	fmt.Println("✅ Esquema cargado y cacheado")
	return codec, nil
}
