package aws

import (
	"context"

	"github.com/puppetlabs/wash/activity"
	"github.com/puppetlabs/wash/plugin"

	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3Client "github.com/aws/aws-sdk-go/service/s3"
)

// s3Dir represents the resources/s3 directory
type s3Dir struct {
	plugin.EntryBase
	session *session.Session
	client  *s3Client.S3
}

func s3DirBase(forInstance bool) *s3Dir {
	s3Dir := &s3Dir{
		EntryBase: plugin.NewEntryBase(),
	}
	s3Dir.SetName("s3").IsSingleton()
	return s3Dir
}

func newS3Dir(session *session.Session) *s3Dir {
	s3Dir := s3DirBase(true)
	s3Dir.session = session
	s3Dir.client = s3Client.New(session)
	return s3Dir
}

func (s *s3Dir) ChildSchemas() []*plugin.EntrySchema {
	return plugin.ChildSchemas(s3BucketBase(false))
}

// List lists the buckets.
func (s *s3Dir) List(ctx context.Context) ([]plugin.Entry, error) {
	resp, err := s.client.ListBucketsWithContext(ctx, nil)
	if err != nil {
		return nil, err
	}

	activity.Record(ctx, "Listing %v S3 buckets", len(resp.Buckets))

	buckets := make([]plugin.Entry, len(resp.Buckets))
	for i, bucket := range resp.Buckets {
		buckets[i] = newS3Bucket(
			awsSDK.StringValue(bucket.Name),
			awsSDK.TimeValue(bucket.CreationDate),
			s.session,
		)
	}

	return buckets, nil
}
