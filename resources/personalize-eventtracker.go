package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/service/personalize"
	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
)

const PersonalizeEventTrackerResource = "PersonalizeEventTracker"

func init() {
	registry.Register(&registry.Registration{
		Name:     PersonalizeEventTrackerResource,
		Scope:    nuke.Account,
		Resource: &PersonalizeEventTracker{},
		Lister:   &PersonalizeEventTrackerLister{},
	})
}

type PersonalizeEventTrackerLister struct{}

func (l *PersonalizeEventTrackerLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := personalize.New(opts.Session)

	params := &personalize.ListEventTrackersInput{}
	resp, err := svc.ListEventTrackers(params)
	if err != nil {
		return nil, err
	}

	resources := make([]resource.Resource, 0)
	for _, tracker := range resp.EventTrackers {
		resources = append(resources, &PersonalizeEventTracker{
			svc: svc,
			id:  tracker.EventTrackerArn,
		})
	}
	return resources, nil
}

type PersonalizeEventTracker struct {
	svc *personalize.Personalize
	id  *string
}

func (r *PersonalizeEventTracker) String() string {
	return *r.id
}

func (p *PersonalizeEventTracker) Remove(_ context.Context) error {
	params := &personalize.DeleteEventTrackerInput{
		EventTrackerArn: p.id,
	}
	_, err := p.svc.DeleteEventTracker(params)
	if err != nil {
		return err
	}

	return nil
}
