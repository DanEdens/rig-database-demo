package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/nuntiodev/nuntio-go-api/api/v1/database"
	nuntio "github.com/nuntiodev/nuntio-go-sdk"
)

const (
	_baseELO         = 1000.0
	_imagesTableName = "images"
)

var _apiKey = os.Getenv("NUNTIO_PROJECT_ID")

type Repository interface {
	Setup(ctx context.Context) error
	AddImage(ctx context.Context, url string) (image, error)
	GetPair(ctx context.Context) ([2]image, error)
	Vote(ctx context.Context, winnerUUID, loserUUID string) error
	ListImages(ctx context.Context) ([]image, error)
	Clean(ctx context.Context) error
}

type dbConfig struct {
	dbname   string
	username string
	password string
}

var client nuntio.Client
var repo Repository

func setupRepository(ctx context.Context) error {
	dbResponse, err := client.Database().GetByName(ctx, connect.NewRequest(&database.GetByNameRequest{
		Name: mongoConfig.dbname,
	}))
	if err != nil {
		return fmt.Errorf("no database found: %q", err)
	}

	endpointResponse, err := client.Database().GetEndpoint(ctx, connect.NewRequest(&database.GetEndpointRequest{
		DatabaseId:   dbResponse.Msg.Database.Id,
		ClientID:     mongoConfig.username,
		ClientSecret: mongoConfig.password,
	}))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get endpoint: %q", err))
	}
	fmt.Printf("uri: %s\n", endpointResponse.Msg.Endpoint)

	repo, err = newMongoRepository(ctx, endpointResponse.Msg.Endpoint, endpointResponse.Msg.DatabaseName)
	if err != nil {
		return err
	}

	if err := repo.Setup(ctx); err != nil {
		return fmt.Errorf("failed to setup repo: %q", err)
	}

	return nil
}

func addImage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	imgURL := r.URL.Query().Get("imgurl")
	img, err := repo.AddImage(ctx, imgURL)
	if err != nil {
		return err
	}
	w.Write([]byte(img.Id))
	return nil
}

type image struct {
	Id  string  `json:"id"`
	Url string  `json:"url"`
	Elo float64 `json:"elo"`
}

func listImages(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	images, err := repo.ListImages(ctx)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(&images)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)

	return nil
}

func clean(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return repo.Clean(ctx)
}

func pair(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	images, err := repo.GetPair(ctx)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(&images)
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}

func vote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	winner := r.URL.Query().Get("winner")
	loser := r.URL.Query().Get("loser")
	return repo.Vote(ctx, winner, loser)
}

func main() {
	client = nuntio.NewClient()
	ctx := context.Background()
	if err := setupRepository(ctx); err != nil {
		log.Fatal(err)
	}

	runServer()
}

func runServer() error {
	http.HandleFunc("/addImage", requestWrapper(addImage))
	http.HandleFunc("/listImages", requestWrapper(listImages))
	http.HandleFunc("/clean", requestWrapper(clean))
	http.HandleFunc("/pair", requestWrapper(pair))
	http.HandleFunc("/vote", requestWrapper(vote))

	err := http.ListenAndServe(":3333", nil)
	return err
}

func requestWrapper(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			fmt.Printf("error: %s\n", err.Error())
			// TODO Extract a proper status code from the err
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
