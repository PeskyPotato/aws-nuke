package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeResourceGatewayResource = "VPCLatticeResourceGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeResourceGatewayResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeResourceGateway{},
		Lister:   &VPCLatticeResourceGatewayLister{},
	})
}

type VPCLatticeResourceGatewayLister struct{}

func (l *VPCLatticeResourceGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := vpclattice.NewFromConfig(*opts.Config)
	var resources []resource.Resource

	params := &vpclattice.ListResourceGatewaysInput{}

	for {
		resp, err := svc.ListResourceGateways(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, gateway := range resp.Items {
			tags, err := svc.ListTagsForResource(ctx, &vpclattice.ListTagsForResourceInput{
				ResourceArn: gateway.Arn,
			})

			if err != nil {
				continue
			}
			resources = append(resources, &VPCLatticeResourceGateway{
				svc:  svc,
				ID:   gateway.Id,
				ARN:  gateway.Arn,
				Tags: tags.Tags,
			})

		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type VPCLatticeResourceGateway struct {
	svc  *vpclattice.Client
	ID   *string
	ARN  *string
	Tags map[string]string
}

func (r *VPCLatticeResourceGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteResourceGateway(ctx, &vpclattice.DeleteResourceGatewayInput{
		ResourceGatewayIdentifier: r.ID,
	})
	return err
}

func (r *VPCLatticeResourceGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeResourceGateway) String() string {
	return *r.ID
}
