package subscriber

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mail-server/internal/model"
	"time"
)

func (r *repository) AddSubscriber(subscriber model.Subscriber) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res bson.M
	filter := bson.D{
		{
			Key:   "email",
			Value: subscriber.Email,
		},
	}

	err := r.storage.Collection("subscribers").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = r.storage.Collection("subscribers").InsertOne(ctx, subscriber)
			if err != nil {
				return false, "", errors.New("[AddSubscriber]: cannot registered subscriber")
			}
			return true, fmt.Sprintf("[AddSubscriber]: new subscriber added"), nil
		}
		return false, "", errors.New("[AddSubscriber]: cannot query database")
	}
	return true, "", nil
}

func (r *repository) FindSubscribers() ([]primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res []bson.M
	cursor, err := r.storage.Collection("subscribers").Find(ctx, bson.D{})
	if err != nil {
		return []bson.M{}, errors.New("[FindSubscribers]: cannot query database")
	}

	if err = cursor.All(ctx, &res); err != nil {
		return []bson.M{}, errors.New("[FindSubscribers]: cannot get all mail")
	}
	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			log.Fatal("[FindSubscribers]: cannot close cursor")
		}
	}()

	if err = cursor.Err(); err != nil {
		return []bson.M{}, errors.New("[FindSubscribers]: cursor error")
	}

	return res, nil
}
