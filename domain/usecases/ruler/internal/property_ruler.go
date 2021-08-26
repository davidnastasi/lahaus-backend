package internal

import "lahaus/domain/model"

type PropertyRulerFunc func(property *model.Property) error

type PropertyRulerFuncs []PropertyRulerFunc
