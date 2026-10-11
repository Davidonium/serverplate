package server

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/davidonium/serverplate/internal/domain"
	"github.com/davidonium/serverplate/internal/templates"
)

func generateHandler(generator *domain.Generator) appHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		lengthEnabled := r.FormValue("length_enabled")
		lengthMode := r.FormValue("length_mode")
		lengthValue, _ := strconv.Atoi(r.FormValue("length_value"))
		componentType := r.URL.Query().Get("component")

		res, err := generator.Generate(r.Context(), domain.GenerateOptions{
			LengthEnabled: lengthEnabled == "on",
			LengthMode:    domain.LengthMode(lengthMode),
			LengthValue:   lengthValue,
		})
		if err != nil {
			return err
		}

		var c templ.Component
		switch componentType {
		case "bucket-input":
			c = templates.BucketNameInput(res.Name)
		default:
			c = templates.GeneratePartial(templates.GenerateViewModel{
				Name: res.Name,
			})
		}

		return component(w, r, http.StatusOK, c)
	}
}
