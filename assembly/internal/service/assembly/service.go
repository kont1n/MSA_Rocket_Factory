package assembly

import (
	def "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
)

var _ def.AssemblyService = (*service)(nil)

type service struct {
	assemblyProducerService def.ProducerService
	metrics                 *assemblyMetrics
}

func NewService(assemblyProducerService def.ProducerService) *service {
	// Инициализируем метрики
	metrics, err := newAssemblyMetrics()
	if err != nil {
		// В случае ошибки создаем nil метрики - сервис должен работать без них
		metrics = nil
	}

	return &service{
		assemblyProducerService: assemblyProducerService,
		metrics:                 metrics,
	}
}
