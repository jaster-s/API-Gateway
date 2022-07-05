package lib

import (
	"fmt"
	"github.com/go-bongo/bongo"
	"os"
)

func DBConnection() *bongo.Connection {
	_logger := new(Logger)
	config := &bongo.Config{
		ConnectionString: getConnectionString(),
	}

	connection, err := bongo.Connect(config)

	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot connect to database: %v", err), true)
	}

	return connection
}

func getConnectionString() string {
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoDatabase := os.Getenv("MONGO_DATABASE")

	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", mongoUsername, mongoPassword, mongoHost, mongoPort, mongoDatabase)
}