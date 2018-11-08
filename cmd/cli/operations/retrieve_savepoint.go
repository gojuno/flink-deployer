package operations

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/spf13/afero"
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

type Filesystem interface {
	ReadDir(dir string) ([]os.FileInfo, error)
}

type LocalFs struct {
	client afero.Fs
}

func (fs *LocalFs) ReadDir(dir string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs.client, dir)
}
