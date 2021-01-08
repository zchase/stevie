package auto_pulumi

import (
	"fmt"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/zchase/stevie/pkg/utils"
)

// buildReactApp builds a react application for a given directory.
func buildReactApp(dirPath string, endpoints []APIEndpoint) error {
	// Set the env variables for the react app.
	for _, endpoint := range endpoints {
		var err error
		endpoint.URL.ApplyString(func(url string) string {
			err = os.Setenv(endpoint.Name, url)
			return url
		})
		if err != nil {
			return err
		}
	}

	// Build the app.
	err := utils.RunCommandWithOutput("yarn", []string{"--cwd", "ui", "build"})
	return err
}

// CreateWebsiteFromDirectoryContents creates a website in S3 from a given path
// to a directory.
func CreateWebsiteFromDirectoryContents(ctx *pulumi.Context, endpoints []APIEndpoint, dirPath, environment string) error {
	// Build the react app.
	err := buildReactApp(dirPath, endpoints)
	if err != nil {
		return utils.NewErrorMessage("Error building react app", err)
	}

	// Create the website bucket.
	bucketName := fmt.Sprintf("%s-website-bucket", environment)
	bucket, err := s3.NewBucket(ctx, bucketName, &s3.BucketArgs{
		Website: s3.BucketWebsiteArgs{
			IndexDocument: pulumi.String("index.html"),
		},
	})
	if err != nil {
		return err
	}

	// Upload the contents of the directory to the bucket.
	siteDir := "ui/build"
	err = utils.RecursiveDirectoryRead(siteDir, "", func(filePath string) error {
		bucketObjectName := fmt.Sprintf("%s-%s-bucket-object", environment, filePath)
		objectFilePath := filepath.Join(siteDir, filePath)
		if _, err := s3.NewBucketObject(ctx, bucketObjectName, &s3.BucketObjectArgs{
			Bucket:      bucket.ID(),
			Key:         pulumi.String(filePath),
			Source:      pulumi.NewFileAsset(objectFilePath),
			ContentType: pulumi.String(mime.TypeByExtension(path.Ext(filePath))),
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil
	}

	// Create an S3 Bucket Policy to allow public read of all objects in bucket.
	bucketPolicyName := fmt.Sprintf("%s-website-bucket-policy", environment)
	if _, err = s3.NewBucketPolicy(ctx, bucketPolicyName, &s3.BucketPolicyArgs{
		Bucket: bucket.ID(),
		Policy: pulumi.Any(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":    "Allow",
					"Principal": "*",
					"Action": []interface{}{
						"s3:GetObject",
					},
					"Resource": []interface{}{
						pulumi.Sprintf("arn:aws:s3:::%s/*", bucket.ID()),
					},
				},
			},
		}),
	}); err != nil {
		return err
	}

	return nil
}
