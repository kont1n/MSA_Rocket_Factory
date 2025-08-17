package assembly

import (
	def "github.com/kont1n/MSA_Rocket_Factory/assembly/internal/service"
)

var _ def.AssemblyService = (*service)(nil)

type service struct {
	assemblyProducerService def.AssemblyProducerService
}

func NewService(assemblyProducerService def.AssemblyProducerService) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
	}
}
