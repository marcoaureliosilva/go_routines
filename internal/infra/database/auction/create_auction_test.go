package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCloseAuctionAutomatically(t *testing.T) {
	// Configura a variável de ambiente para a duração do leilão
	os.Setenv("AUCTION_INTERVAL", "1s")

	// Conexão com o MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Erro ao conectar no MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	repo := NewAuctionRepository(client.Database("auction_db"))
	auction := &auction_entity.Auction{Id: "test-auction", ProductName: "Produto Teste", Category: "Categoria Teste", Condition: auction_entity.ProductCondition_NEW}

	// Cria o leilão
	if err := repo.CreateAuction(context.Background(), auction); err != nil {
		t.Fatalf("Erro ao criar leilão: %v", err)
	}

	// Aguarda o leilão ser fechado
	time.Sleep(2 * time.Second)

	var closedAuction AuctionEntityMongo
	err = repo.Collection.FindOne(context.Background(), bson.M{"_id": auction.Id}).Decode(&closedAuction)
	if err != nil {
		t.Fatalf("Erro ao buscar leilão fechado: %v", err)
	}

	if closedAuction.Status != auction_entity.AuctionStatus_CLOSED {
		t.Fatalf("O leilão %s não foi fechado automaticamente.", auction.Id)
	}
}
