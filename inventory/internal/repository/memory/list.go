package memory

import (
	"context"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"

	repoConverter "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
	inventoryV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/inventory/v1"
)

func (r *repository) ListParts(ctx context.Context, filter *model.Filter) (*[]model.Part, error) {
	partsFiltered := make([]*repoModel.RepositoryPart, 0)
	r.mu.RLock()
	for _, part := range r.data {
		partsFiltered = append(partsFiltered, part)
	}
	r.mu.RUnlock()

	if filter != nil {
		partsFiltered = filtration(filter, partsFiltered)
	}

	parts := make([]model.Part, 0, len(r.data))
	for _, repoPart := range r.data {
		part := *repoConverter.RepoToModel(repoPart)
		parts = append(parts, part)
	}

	return &parts, nil
}

func filtration(filter *model.Filter, parts []*repoModel.RepositoryPart) (result []*repoModel.RepositoryPart) {
	// Создаем мап для фильтрации
	uuidSet := make(map[string]bool)
	for _, uuid := range filter.GetPartUuid() {
		uuidSet[uuid] = true
	}

	nameSet := make(map[string]bool)
	for _, name := range filter.GetPartName() {
		nameSet[name] = true
	}

	categorySet := make(map[inventoryV1.Category]bool)
	for _, category := range filter.GetCategory() {
		categorySet[category] = true
	}

	manufacturerCountrySet := make(map[string]bool)
	for _, manufacturerCountry := range filter.GetManufacturerCountry() {
		manufacturerCountrySet[manufacturerCountry] = true
	}

	tagSet := make(map[string]bool)
	for _, tag := range filter.GetTags() {
		tagSet[tag] = true
	}

	// Фильтруем детали
	for _, part := range parts {
		if len(uuidSet) > 0 {
			if _, ok := uuidSet[part.PartUuid]; !ok {
				continue
			}
		}

		if len(nameSet) > 0 {
			if _, ok := nameSet[part.Name]; !ok {
				continue
			}
		}

		if len(categorySet) > 0 {
			if _, ok := categorySet[part.Category]; !ok {
				continue
			}
		}

		if len(manufacturerCountrySet) > 0 {
			if _, ok := manufacturerCountrySet[part.Manufacturer.Country]; !ok {
				continue
			}
		}

		if len(tagSet) > 0 {
			if _, ok := tagSet[part.Tags[0]]; !ok {
				continue
			}
		}

		result = append(result, part)
	}

	return result
}
