package auction

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
	EndsAt      int64                           `bson:"ends_at"` // Adiciona o campo EndsAt
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

var mu sync.Mutex

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	duration := getAuctionDuration()
	endsAt := time.Now().Add(duration).Unix()

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auction_entity.AuctionStatus_OPEN,
		Timestamp:   auctionEntity.Timestamp.Unix(),
		EndsAt:      endsAt,
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	// Inicia a goroutine para fechar o leilão automaticamente
	go closeAuctionAfterDuration(ctx, auctionEntityMongo)

	return nil
}

// fecha o leilão se o tempo se esgotar.
func closeAuctionAfterDuration(ctx context.Context, auction *AuctionEntityMongo) {
	time.Sleep(time.Unix(auction.EndsAt, 0).Sub(time.Now()))

	mu.Lock()
	defer mu.Unlock()

	// Atualiza o status do leilão
	auction.Status = auction_entity.AuctionStatus_CLOSED
	_, err := auction.Collection.UpdateOne(ctx, bson.M{"_id": auction.Id}, bson.M{"$set": bson.M{"status": auction.Status}})
	if err != nil {
		logger.Error("Error closing auction", err)
		return
	}
	fmt.Printf("Leilão %s fechado automaticamente.\n", auction.Id)
}

// Função para obter a duração do leilão a partir de variáveis de ambiente.
func getAuctionDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("AUCTION_INTERVAL"))
	if err != nil {
		return 10 * time.Minute
	}
	return duration
}
