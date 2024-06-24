package handler

import (
    "context"
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "path/filepath"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
    Name         string `bson:"name" json:"name"`
    PhoneNumber  string `bson:"phone_number" json:"phone_number"`
    Address      string `bson:"address" json:"address"`
    DeliveryTime string `bson:"delivery_time" json:"delivery_time"`
}

var client *mongo.Client

func init() {
    // Connect to MongoDB
    clientOptions := options.Client().ApplyURI("mongodb+srv://prachhhi:oprybBJBWko7zbjE@cluster0.r487mib.mongodb.net/?retryWrites=true&w=majority")
    var err error
    client, err = mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = client.Disconnect(context.Background()); err != nil {
            log.Fatal(err)
        }
    }()
}

// HandleFormSubmission is the exported handler function for form submission
func HandleFormSubmission(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var order Order
        order.Name = r.FormValue("name")
        order.PhoneNumber = r.FormValue("phone_number")
        order.Address = r.FormValue("address")
        order.DeliveryTime = r.FormValue("delivery_time")

        collection := client.Database("order_locator").Collection("orders")
        _, err := collection.InsertOne(context.Background(), order)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Respond with success status
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(order)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
