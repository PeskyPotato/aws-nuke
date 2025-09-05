package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/personalize"
	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
)

const PersonalizeDatasetGroupResource = "PersonalizeDatasetGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     PersonalizeDatasetGroupResource,
		Scope:    nuke.Account,
		Resource: &PersonalizeDatasetGroup{},
		Lister:   &PersonalizeDatasetGroupLister{},
		DependsOn: []string{
			PersonalizeDatasetResource,
		},
	})
}

type PersonalizeDatasetGroupLister struct{}

func (l *PersonalizeDatasetGroupLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := personalize.New(opts.Session)

	params := &personalize.ListDatasetGroupsInput{}
	resp, err := svc.ListDatasetGroups(params)
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0)
	for _, group := range resp.DatasetGroups {
		resources = append(resources, &PersonalizeDatasetGroup{
			svc: svc,
			id:  group.DatasetGroupArn,
		})
	}
	return resources, nil
}

type PersonalizeDatasetGroup struct {
	svc *personalize.Personalize
	id  *string
}

func (r *PersonalizeDatasetGroup) String() string {
	return *r.id
}

func (p *PersonalizeDatasetGroup) Remove(_ context.Context) error {
	params := &personalize.DeleteDatasetGroupInput{
		DatasetGroupArn: p.id,
	}
	_, err := p.svc.DeleteDatasetGroup(params)
	if err != nil {
		return err
	}

	return nil
}
