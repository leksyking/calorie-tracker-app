package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/leksyking/calorie-tracker-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	validate                          = validator.New()
	entryCollection *mongo.Collection = OpenCollection(Client, "calories")
)

func AddEntry(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	entry.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprintln("entry item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func GetEntries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err := cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntriesByIngredient(c *gin.Context) {
	ingredient := c.Params.ByName("ingredients")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	if err := cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	validationErr := validate.Struct(entry)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	filter := bson.M{"_id": docID}
	update := bson.M{"dish": entry.Dish, "calories": entry.Calories, "ingredients": entry.Ingredients, "fat": entry.Fat}
	result, err := entryCollection.ReplaceOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, result.ModifiedCount)
}
func UpdateIngredient(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := c.BindJSON(&ingredient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID}, bson.D{{Key: "$set", Value: bson.D{{Key: "ingredients", Value: ingredient.Ingredients}}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func DeleteEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, result.DeletedCount)
}
