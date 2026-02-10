//go:build integration
// +build integration

package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/user/gin-microservice-boilerplate/models"
)

func setupMongoContainer(t *testing.T) *mongo.Database {
	t.Helper()
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:7")
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Disconnect(ctx) })

	return client.Database("testdb")
}

func TestUserRepoMongo_CreateAndGetByID(t *testing.T) {
	db := setupMongoContainer(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	id, err := repo.Create(ctx, &models.User{
		Name:  "Alice",
		Email: "alice@example.com",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	user, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.NotZero(t, user.CreatedAt)
}

func TestUserRepoMongo_GetAll(t *testing.T) {
	db := setupMongoContainer(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	_, err := repo.Create(ctx, &models.User{Name: "Alice", Email: "alice@example.com"})
	require.NoError(t, err)
	_, err = repo.Create(ctx, &models.User{Name: "Bob", Email: "bob@example.com"})
	require.NoError(t, err)

	users, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepoMongo_Update(t *testing.T) {
	db := setupMongoContainer(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	id, err := repo.Create(ctx, &models.User{Name: "Alice", Email: "alice@example.com"})
	require.NoError(t, err)

	err = repo.Update(ctx, &models.User{
		ID:    id,
		Name:  "Alice Updated",
		Email: "alice.updated@example.com",
	})
	require.NoError(t, err)

	user, err := repo.GetByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Alice Updated", user.Name)
	assert.Equal(t, "alice.updated@example.com", user.Email)
}

func TestUserRepoMongo_Delete(t *testing.T) {
	db := setupMongoContainer(t)
	repo := NewUserRepo(db)
	ctx := context.Background()

	id, err := repo.Create(ctx, &models.User{Name: "Alice", Email: "alice@example.com"})
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, id)
	assert.Error(t, err)
}
