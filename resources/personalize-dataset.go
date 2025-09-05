package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/personalize"
	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
)

const PersonalizeDatasetResource = "PersonalizeDataset"

func init() {
	registry.Register(&registry.Registration{
		Name:     PersonalizeDatasetResource,
		Scope:    nuke.Account,
		Resource: &PersonalizeDataset{},
		Lister:   &PersonalizeDatasetLister{},
	})
}

type PersonalizeDatasetLister struct{}

func (l *PersonalizeDatasetLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := personalize.New(opts.Session)

	params := &personalize.ListDatasetsInput{}
	resp, err := svc.ListDatasets(params)
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0)
	for _, group := range resp.Datasets {
		resources = append(resources, &PersonalizeDataset{
			svc: svc,
			id:  group.DatasetArn,
		})
	}
	return resources, nil
}

type PersonalizeDataset struct {
	svc *personalize.Personalize
	id  *string
}

func (r *PersonalizeDataset) String() string {
	return *r.id
}

func (p *PersonalizeDataset) Remove(_ context.Context) error {
	params := &personalize.DeleteDatasetInput{
		DatasetArn: p.id,
	}
	_, err := p.svc.DeleteDataset(params)
	if err != nil {
		return err
	}

	return nil
}
