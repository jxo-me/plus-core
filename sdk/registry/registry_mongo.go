package registry

import "go.mongodb.org/mongo-driver/mongo"

type MongoRegistry struct {
	registry[*mongo.Client]
}
