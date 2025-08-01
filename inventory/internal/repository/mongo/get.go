package mongo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoConverter "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid uuid.UUID) (*model.Part, error) {
	collection := r.db.Collection(partsCollection)

	filter := bson.M{"order_uuid": uuid.String()}

	var repoPart repoModel.RepositoryPart
	err := collection.FindOne(ctx, filter).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	// Преобразуем модель репозитория в основную модель
	part, err := repoConverter.ToModelPart(&repoPart)
	if err != nil {
		return nil, err
	}

	return part, nil
}
