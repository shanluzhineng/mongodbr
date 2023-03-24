package mongodbr

import "go.mongodb.org/mongo-driver/mongo/options"

type FindOneOption func(*options.FindOneOptions)
