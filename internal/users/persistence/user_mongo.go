package persistence

import (
	"github.com/mohamedveron/go_app_template/internal/pkg/datastore"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserCollection = "users"

type UserMongoPersistence struct {
	collection *mongo.Collection
}

func NewUserMongoPersistence(mongodbCli *datastore.MongoDB) *UserMongoPersistence {
	return &UserMongoPersistence{
		collection: mongodbCli.Database.Collection(UserCollection),
	}
}
