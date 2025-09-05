package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/personalize"
	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
)

const PersonalizeSchemaResource = "PersonalizeSchema"

func init() {
	registry.Register(&registry.Registration{
		Name:     PersonalizeSchemaResource,
		Scope:    nuke.Account,
		Resource: &PersonalizeSchema{},
		Lister:   &PersonalizeSchemaLister{},
		DependsOn: []string{
			PersonalizeDatasetResource,
		},
	})
}

type PersonalizeSchemaLister struct{}

func (l *PersonalizeSchemaLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := personalize.New(opts.Session)

	params := &personalize.ListSchemasInput{}
	resp, err := svc.ListSchemas(params)
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0)
	for _, schema := range resp.Schemas {
		resources = append(resources, &PersonalizeSchema{
			svc: svc,
			id:  schema.SchemaArn,
		})
	}
	return resources, nil
}

type PersonalizeSchema struct {
	svc *personalize.Personalize
	id  *string
}

func (r *PersonalizeSchema) String() string {
	return *r.id
}

func (p *PersonalizeSchema) Remove(_ context.Context) error {
	params := &personalize.DeleteSchemaInput{
		SchemaArn: p.id,
	}
	_, err := p.svc.DeleteSchema(params)
	if err != nil {
		return err
	}

	return nil
}
