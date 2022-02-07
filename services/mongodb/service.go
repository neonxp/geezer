package mongodb

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/neonxp/geezer"
	"github.com/neonxp/geezer/render"
)

type Service[T any] struct {
	collection *mongo.Collection
}

func New[T any](collection *mongo.Collection) *Service[T] {
	return &Service[T]{collection: collection}
}

func (s Service[T]) Find(ctx context.Context, params geezer.Params) (render.Renderer, error) {
	var model []*T
	where := bson.D{}
	for k, v := range params.Query {
		where = append(where, bson.E{Key: k, Value: v[0]})
	}
	cursor, err := s.collection.Find(ctx, where)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		r := new(T)
		if err := cursor.Decode(r); err != nil {
			return nil, err
		}
		model = append(model, r)
	}

	return render.JSON(model), nil
}

func (s Service[T]) Get(ctx context.Context, id string, params geezer.Params) (render.Renderer, error) {
	var model T
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	sr := s.collection.FindOne(ctx, bson.M{"_id": oid})
	if err := sr.Err(); err != nil {
		return nil, err
	}
	if err := sr.Decode(&model); err != nil {
		return nil, err
	}

	return render.JSON(model), nil
}

func (s Service[T]) Create(ctx context.Context, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	var model T
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}
	ir, err := s.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return render.JSON(InsertResult{
		ID:   ir.InsertedID,
		Item: model,
	}), nil
}

func (s Service[T]) Update(ctx context.Context, id string, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service[T]) Patch(ctx context.Context, id string, data geezer.Data, params geezer.Params) (render.Renderer, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service[T]) Remove(ctx context.Context, id string, params geezer.Params) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (s Service[T]) Setup(app geezer.AppKernel, path string) error {
	return nil
}
