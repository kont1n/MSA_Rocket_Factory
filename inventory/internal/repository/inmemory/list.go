package inmemory

import (
	"context"

	"github.com/samber/lo"

	"github.com/kont1n/MSA_Rocket_Factory/inventory/internal/model"
	repoConverter "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/inventory/internal/repository/model"
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

	parts := make([]model.Part, 0, len(partsFiltered))
	for _, repoPart := range partsFiltered {
		part, err := repoConverter.RepoToModel(repoPart)
		if err != nil {
			return nil, err
		}
		parts = append(parts, lo.FromPtr(part))
	}

	return &parts, nil
}

func filtration(filter *model.Filter, parts []*repoModel.RepositoryPart) (result []*repoModel.RepositoryPart) {
	// Создаем мап для фильтрации
	uuidSet := make(map[string]bool)
	for _, uuid := range filter.Uuids {
		uuidSet[uuid.String()] = true
	}

	nameSet := make(map[string]bool)
	for _, name := range filter.Names {
		nameSet[name] = true
	}

	categorySet := make(map[int]bool)
	for _, category := range filter.Categories {
		categorySet[int(category)] = true
	}

	manufacturerCountrySet := make(map[string]bool)
	for _, manufacturerCountry := range filter.ManufacturerCountries {
		manufacturerCountrySet[manufacturerCountry] = true
	}

	tagSet := make(map[string]bool)
	for _, tag := range filter.Tags {
		tagSet[tag] = true
	}

	// Фильтруем детали
	for _, part := range parts {
		if len(uuidSet) > 0 {
			if _, ok := uuidSet[part.OrderUuid]; !ok {
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
			// Проверяем, есть ли хотя бы один тег из фильтра в тегах детали
			tagFound := false
			for _, partTag := range part.Tags {
				if _, ok := tagSet[partTag]; ok {
					tagFound = true
					break
				}
			}
			if !tagFound {
				continue
			}
		}

		result = append(result, part)
	}

	return result
}
