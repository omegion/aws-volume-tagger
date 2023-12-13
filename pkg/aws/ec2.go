package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2 struct {
	Client *ec2.Client
}

func (e EC2) DescribeVolumes(ctx context.Context, filters []types.Filter) ([]types.Volume, error) {
	nextToken := aws.String("")
	var volumes []types.Volume

	for nextToken != nil {
		result, err := e.Client.DescribeVolumes(
			ctx,
			&ec2.DescribeVolumesInput{
				Filters:    filters,
				MaxResults: aws.Int32(500),
				NextToken:  nextToken,
			},
		)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, result.Volumes...)

		nextToken = result.NextToken
	}

	return volumes, nil
}

func (e EC2) CreateTags(ctx context.Context, volumeID string, tags []types.Tag) error {
	_, err := e.Client.CreateTags(
		ctx,
		&ec2.CreateTagsInput{
			Resources: []string{volumeID},
			Tags:      tags,
		},
	)

	return err
}
