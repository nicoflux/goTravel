package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectToMongoDB() (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://user:qwe123@cluster0.szzuy.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}

func closeMongoDBConnection(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		fmt.Println("Error al desconectar de MongoDB:", err)
	}
}

func insertData(client *mongo.Client) error {
	collection := client.Database("goTravel").Collection("reservations")
	data := bson.M{
		"campo1": "campo1",
		"campo2": 1,
		"campo3": "campo3",
	}

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	fmt.Println("Documento insertado con éxito")
	return nil
}

func main() {
	client, err := connectToMongoDB()
	if err != nil {
		fmt.Println("Error al conectar a MongoDB:", err)
		return
	}
	defer closeMongoDBConnection(client)

	if err := insertData(client); err != nil {
		fmt.Println("Error al insertar datos en MongoDB:", err)
	}
}
