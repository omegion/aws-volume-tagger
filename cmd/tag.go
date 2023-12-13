package cmd

import (
	"context"
	"fmt"
	awsUtils "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/omegion/aws-volume-tagger/pkg/aws"
	"strings"

	"github.com/spf13/cobra"
)

func TagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag the volumes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err := cmd.Flags().GetBool("dry-run")

			cfg, err := config.LoadDefaultConfig(context.TODO())
			if err != nil {
				return err
			}

			svc := aws.EC2{Client: ec2.NewFromConfig(cfg)}
			volumes, err := svc.DescribeVolumes(cmd.Context(), []types.Filter{
				{
					Name:   awsUtils.String("tag-key"),
					Values: []string{"kubernetes.io/created-for/pvc/name"},
				},
			})

			if err != nil {
				return err
			}

			taggedVolumes := 0
			for _, volume := range volumes {
				csiCreated := false
				pvName := ""
				kubernetesCluster := ""

				for _, tag := range volume.Tags {
					switch *tag.Key {
					case "CSIVolumeName":
						csiCreated = true
					case "kubernetes.io/created-for/pv/name":
						pvName = *tag.Value
					}

					if strings.Contains(*tag.Key, "kubernetes.io/cluster/") {
						kubernetesCluster, _ = strings.CutPrefix(*tag.Key, "kubernetes.io/cluster/")
					}
				}

				if !csiCreated {
					taggedVolumes++
					printVolumeInfo(kubernetesCluster, pvName, *volume.VolumeId)

					if !dryRun {
						err = svc.CreateTags(cmd.Context(), *volume.VolumeId, []types.Tag{
							{
								Key:   awsUtils.String("CSIVolumeName"),
								Value: awsUtils.String(pvName),
							},
							{
								Key:   awsUtils.String("ebs.csi.aws.com/cluster"),
								Value: awsUtils.String("true"),
							},
						})
						if err != nil {
							return err
						}
					}
				}
			}

			fmt.Printf("---\n%d volumes tagged successfully\n", taggedVolumes)

			return nil
		},
	}

	cmd.Flags().BoolP("dry-run", "d", false, "Dry run")

	return cmd
}

func printVolumeInfo(kubernetesCluster, pvName, volumeId string) {
	fmt.Println("---")
	fmt.Println("Cluster Name:\t", kubernetesCluster)
	fmt.Println("PVC Name:\t", pvName)
	fmt.Println("Volume ID:\t", volumeId)
}
