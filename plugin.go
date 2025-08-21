package created_by

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbuilds/xpb"
)

type Plugin struct {
	// Single relation fields to auth collections to
	//   automatically set on record create in the
	//   format "collection.field_name".
	Fields []string `json:"fields"`
}

func init() {
	xpb.Register(&Plugin{})
}

// Name implements xpb.Plugin.
func (p *Plugin) Name() string {
	return "created_by"
}

// This variable will automatically be set at build time by xpb
var version string

// Version implements xpb.Plugin.
func (p *Plugin) Version() string {
	return version
}

// Description implements xpb.Plugin.
func (p *Plugin) Description() string {
	return "Allows for easily configured created_by fields."
}

// Init implements xpb.Plugin.
func (p *Plugin) Init(app core.App) error {
	for _, raw := range p.Fields {
		collectionName, fieldName, _ := strings.Cut(raw, ".")
		app.OnRecordCreateRequest(collectionName).
			BindFunc(p.setCreatedByField(fieldName))
	}
	return nil
}

// copied from core.collectionNameRegex
var collectionNameRegex = regexp.MustCompile(`^\w+$`)

// Validate implements validation.Validatable.
func (p *Plugin) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Fields,
			validation.Each(
				validation.By(func(value any) error {
					str, err := validation.EnsureString(value)
					if err != nil {
						return err
					}
					if strings.Count(str, ".") != 1 {
						return errors.New("must contain a single period separator")
					}
					collectionName, fieldName, _ := strings.Cut(str, ".")
					if err := validation.Validate(collectionName,
						validation.Match(collectionNameRegex),
					); err != nil {
						return err
					}
					if err := validation.Validate(fieldName,
						validation.By(core.DefaultFieldNameValidationRule),
					); err != nil {
						return err
					}
					return nil
				}),
			),
		),
	)
}

func (p *Plugin) setCreatedByField(fieldName string) func(e *core.RecordRequestEvent) error {
	return func(e *core.RecordRequestEvent) error {

		if e.Auth == nil {
			return e.Next()
		}

		if existingValue := e.Record.GetString(fieldName); existingValue != "" {
			return e.Next()
		}

		field, ok := e.Record.Collection().Fields.GetByName(fieldName).(*core.RelationField)
		if !ok {
			// TODO: log warning
			return e.Next()
		}

		if field.IsMultiple() {
			// TODO: log warning
			return e.Next()
		}

		if e.Auth.Collection().Id != field.CollectionId {
			return e.Next()
		}

		e.Record.Set(fieldName, e.Auth.Id)

		return e.Next()
	}
}
