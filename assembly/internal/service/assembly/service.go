package assembly

import (
	def "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
)

var _ def.AssemblyService = (*service)(nil)

type service struct {
	assemblyProducerService def.ProducerService
}

func NewService(assemblyProducerService def.ProducerService) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
	}
}
