package operations

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (o RealOperator) retrieveLatestSavepoint(dir string) (string, error) {
	if strings.HasSuffix(dir, "/") {
		dir = strings.TrimSuffix(dir, "/")
	}

	u, err := url.Parse(dir)
	if err != nil {
		return "", err
	}

	output, err := o.S3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(u.Host),
		Prefix: aws.String(u.Path),
	})
	if err != nil {
		return "", fmt.Errorf("failed to load objects from s3 %q: %v", dir, err)
	}

	if len(output.Contents) == 0 {
		return "", errors.New("No savepoints present in directory: " + dir)
	}

	var newestFile string
	var newestTime int64
	for _, o := range output.Contents {
		if o.LastModified == nil {
			continue
		}
		currTime := o.LastModified.Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = filepath.Join(dir, aws.StringValue(o.Key))
		}
	}

	return newestFile, nil
}
