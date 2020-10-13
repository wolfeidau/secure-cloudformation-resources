package cf

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/s3"
)

type Event struct {
	Region                  string                   `json:"region"`
	AccountID               string                   `json:"accountId"`
	Fragment                *cloudformation.Template `json:"fragment"`
	Params                  map[string]interface{}   `json:"params"`
	RequestID               string                   `json:"requestId"`
	TemplateParameterValues map[string]interface{}   `json:"templateParameterValues"`
	TransformID             string                   `json:"transformId"`
}

type Result struct {
	RequestID string                   `json:"requestId"`
	Status    string                   `json:"status"`
	Fragment  *cloudformation.Template `json:"fragment"`
}

type Transform struct {
}

func (t *Transform) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	req := new(Event)

	err := json.Unmarshal(payload, req)
	if err != nil {
		return nil, err
	}

	buckets := req.Fragment.GetAllS3BucketResources()

	for name, bucket := range buckets {
		log.Ctx(ctx).Info().Str("name", name).Msg("bucket")

		if bucket.PublicAccessBlockConfiguration == nil {

			log.Ctx(ctx).Info().Str("name", name).Msg("update public access configuration")

			pbc := new(s3.Bucket_PublicAccessBlockConfiguration)
			pbc.BlockPublicAcls = true
			pbc.IgnorePublicAcls = true
			pbc.RestrictPublicBuckets = true
			pbc.BlockPublicPolicy = true

			bucket.PublicAccessBlockConfiguration = pbc
		}

		if bucket.BucketEncryption == nil {

			log.Ctx(ctx).Info().Str("name", name).Msg("update encryption configuration")

			be := new(s3.Bucket_BucketEncryption)

			rules := []s3.Bucket_ServerSideEncryptionRule{}
			rules = append(rules, s3.Bucket_ServerSideEncryptionRule{
				ServerSideEncryptionByDefault: &s3.Bucket_ServerSideEncryptionByDefault{
					SSEAlgorithm: "AES256",
				},
			})

			bucket.BucketEncryption = be
			bucket.BucketEncryption.ServerSideEncryptionConfiguration = rules
		}
	}

	res := new(Result)

	res.Fragment = req.Fragment
	res.RequestID = req.RequestID
	res.Status = "success"

	return json.Marshal(res)
}
