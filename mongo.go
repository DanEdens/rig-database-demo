package main

import (
	"context"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/exp/slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoConfig = dbConfig{
	dbname:   os.Getenv("DATABASE_NAME"),
	username: os.Getenv("DATABASE_CLIENT_ID"),
	password: os.Getenv("DATABASE_CLIENT_SECRET"),
}

func newMongoRepository(ctx context.Context, uri string, dbname string) (Repository, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %q", err)
	}

	return &mongoRepository{
		db: mongoClient.Database(dbname),
	}, nil
}

type mongoRepository struct {
	db *mongo.Database
}

func (m *mongoRepository) Setup(ctx context.Context) error {
	collections, err := m.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}
	if slices.Contains(collections, _imagesTableName) {
		return nil
	}
	return m.db.CreateCollection(ctx, _imagesTableName)
}

func (m *mongoRepository) AddImage(ctx context.Context, url string) (image, error) {
	collection := m.db.Collection(_imagesTableName)
	resp, err := collection.InsertOne(ctx, bson.D{
		{Key: "url", Value: url},
		{Key: "elo", Value: _baseELO},
	})
	if err != nil {
		return image{}, err
	}
	oid := resp.InsertedID.(primitive.ObjectID)
	return image{
		Id:  oid.String(),
		Url: url,
		Elo: _baseELO,
	}, nil
}

func (m *mongoRepository) GetPair(ctx context.Context) ([2]image, error) {
	var empty [2]image

	collection := m.db.Collection(_imagesTableName)
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{
		{
			{
				Key: "$sample",
				Value: bson.D{{
					Key:   "size",
					Value: 2,
				}},
			},
		},
	})
	if err != nil {
		return empty, fmt.Errorf("failed to sample 2: %q", err)
	}

	var mongoImages [2]mongoImage
	if err := cursor.All(ctx, &mongoImages); err != nil {
		return empty, err
	}

	return [2]image{
		mongoImages[0].toImage(),
		mongoImages[1].toImage(),
	}, nil
}

func (m *mongoRepository) Vote(ctx context.Context, winnerUUID, loserUUID string) error {
	winnerID, err := primitive.ObjectIDFromHex(winnerUUID)
	if err != nil {
		return err
	}
	loserID, err := primitive.ObjectIDFromHex(loserUUID)
	if err != nil {
		return err
	}

	collection := m.db.Collection(_imagesTableName)
	var winnerImg image
	if err := collection.FindOne(ctx, bson.M{"_id": winnerID}).Decode(&winnerImg); err != nil {
		return err
	}

	var loserImg image
	if err := collection.FindOne(ctx, bson.M{"_id": loserID}).Decode(&loserImg); err != nil {
		return err
	}

	// TODO Proper ELO calculation
	winnerImg.Elo += 100
	loserImg.Elo -= 100
	if _, err := collection.UpdateByID(ctx, winnerID, bson.M{"$set": bson.M{"elo": winnerImg.Elo}}); err != nil {
		return err
	}
	if _, err := collection.UpdateByID(ctx, loserID, bson.M{"$set": bson.M{"elo": loserImg.Elo}}); err != nil {
		return err
	}

	return nil
}

func (m *mongoRepository) ListImages(ctx context.Context) ([]image, error) {
	collection := m.db.Collection(_imagesTableName)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var mongoImages []mongoImage
	if err := cur.All(ctx, &mongoImages); err != nil {
		return nil, err
	}

	var images []image
	for _, img := range mongoImages {
		images = append(images, img.toImage())
	}
	return images, nil
}

func (m *mongoRepository) Clean(ctx context.Context) error {
	collection := m.db.Collection(_imagesTableName)
	_, err := collection.DeleteMany(ctx, bson.D{})
	return err
}

type mongoImage struct {
	Id  primitive.ObjectID `json:"_id" bson:"_id"`
	Url string             `json:"url" bson:"url"`
	Elo float64            `json:"elo" bson:"elo"`
}

func (m mongoImage) toImage() image {
	return image{
		Id:  m.Id.String(),
		Url: m.Url,
		Elo: m.Elo,
	}
}
