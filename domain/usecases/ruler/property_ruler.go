package ruler

import (
	"lahaus/config"
	"lahaus/domain/model"
	"lahaus/domain/usecases/ruler/internal"
)

type PropertyRules struct {
	rulers []internal.PropertyRulerFunc
}

func NewPropertyRulerUseCase(config *config.Config) *PropertyRules {
	pv := &PropertyRules{
		rulers: []internal.PropertyRulerFunc{
			internal.NewLocationRuler(),
			internal.NewPropertyTypeRuler(config),
			internal.NewPriceRuler(config),
		},
	}
	return pv
}

func (pv *PropertyRules) Execute(property *model.Property) {
	for _, rule := range pv.rulers {
		if err := rule(property); err != nil {
			property.Status = model.INVALID
			break
		}
	}
}
