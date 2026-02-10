package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/user/gin-microservice-boilerplate/internal/repository"
	"github.com/user/gin-microservice-boilerplate/models"
)

const collectionName = "users"

type userRepoMongo struct {
	collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) repository.UserRepository {
	return &userRepoMongo{
		collection: db.Collection(collectionName),
	}
}

func (r *userRepoMongo) Create(ctx context.Context, user *models.User) (string, error) {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	oid := result.InsertedID.(primitive.ObjectID)
	return oid.Hex(), nil
}

func (r *userRepoMongo) GetByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}

	user.ID = oid.Hex()
	return &user, nil
}

func (r *userRepoMongo) GetAll(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, cursor.Err()
}

func (r *userRepoMongo) Update(ctx context.Context, user *models.User) error {
	oid, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	return err
}

func (r *userRepoMongo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
