package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDocument struct {
	Token     string    `bson:"token"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}

type MongoTokenRepository struct {
	collection *mongo.Collection
}

func NewMongoTokenRepository(db *mongo.Database) *MongoTokenRepository {
	return &MongoTokenRepository{
		collection: db.Collection("login_tokens"),
	}
}

func (r *MongoTokenRepository) Store(token string, userID string, expiry time.Time) error {
	_, err := r.collection.InsertOne(context.Background(), TokenDocument{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiry,
	})
	return err
}

func (r *MongoTokenRepository) Verify(token string) (string, error) {
	var doc TokenDocument
	err := r.collection.FindOne(context.Background(), bson.M{
		"token":      token,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&doc)

	if err != nil {
		return "", err
	}

	return doc.UserID, nil
}

func (r *MongoTokenRepository) Delete(token string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"token": token})
	return err
}
