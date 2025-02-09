package mongodb

import (
	"context"
	"time"

	"github.com/subscription-tracker/user/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	_, err := r.collection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}
	userData := &domain.User{
		Email:              user.Email,
		Name:               user.Name,
		Password:           user.Password,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          time.Now(),
		SpotifyCredentials: user.SpotifyCredentials,
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": userData},
	)
	return err
}

func (r *UserRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
