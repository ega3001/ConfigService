package inputStructs

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

// Resolver validators
type LPRList struct {
	List []LPRListElement
}

type LPRListElement struct {
	Plate string `json:"plate"`
	Brand string `json:"brand,omitempty"`
	Model string `json:"model,omitempty"`
	Color string `json:"color,omitempty"`
}

type LPRListPatchedElement struct {
	Plate string `json:"plate,omitempty"`
	Brand string `json:"brand,omitempty"`
	Model string `json:"model,omitempty"`
	Color string `json:"color,omitempty"`
}

// Resolvers
func (ml MListElemsPatch) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	validator := huma.NewModelValidator()
	switch ml.ModuleName {
	case "LPR":
		for _, elem := range ml.Body {
			var bodyMap map[string]any
			if err := mapstructure.Decode(elem, &bodyMap); err != nil {
				return []error{err}
			}

			if errs := validator.Validate(
				reflect.TypeOf(LPRListElement{}),
				bodyMap,
			); errs != nil {
				return errs
			}
		}
	default:
		return nil
	}
	return nil
}

func (ml MListElemPatch) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	validator := huma.NewModelValidator()
	switch ml.ModuleName {
	case "LPR":
		var bodyMap map[string]any
		if err := mapstructure.Decode(ml.Body, &bodyMap); err != nil {
			return []error{err}
		}

		if errs := validator.Validate(
			reflect.TypeOf(LPRListPatchedElement{}),
			bodyMap,
		); errs != nil {
			return errs
		}
	default:
		return nil
	}
	return nil
}
