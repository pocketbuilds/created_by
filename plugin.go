package created_by

import (
	"strings"

	"github.com/PocketBuilds/xpb"
	"github.com/pocketbase/pocketbase/core"
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
		// TODO: validation update on xpb
		// Prevalidate (interface for filling defaults programatically)
		// validation.Validate
		collectionName, fieldName, _ := strings.Cut(raw, ".")
		app.OnRecordCreateRequest(collectionName).
			BindFunc(p.setCreatedByField(fieldName))
	}
	return nil
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
