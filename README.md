# secure-cloudformation-resources

This is an effort to create a group of cloudformation resources which implement a security baseline as a parameter to their creation, which influences their default configuration to help keep up with the latest recommendations.

For more background on this problem take a look at my blog post [Why isn't my s3 bucket secure?](https://www.wolfe.id.au/2020/10/08/why-isnt-my-s3-bucket-secure/).

# Example

Keeping things simple for developers this is all that is required to create an S3 bucket with the baseline selected at creation.

```yaml
AWSTemplateFormatVersion: "2010-09-09"

Transform:
  - "YOUR_ACCOUNT_ID::SecurityTransform"
Resources:
  MyDataBucket:
    Type: Secure::S3::Bucket
    Properties: 
      Baseline: standards/aws-foundational-security-best-practices/v/1.0.0
```

After the magic of translation this becomes.

```json
{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Resources": {
    "MyDataBucket": {
      "Type": "AWS::S3::Bucket",
      "Properties": {
        "BucketEncryption": {
          "ServerSideEncryptionConfiguration": [
            {
              "ServerSideEncryptionByDefault": {
                "SSEAlgorithm": "AES256"
              }
            }
          ]
        },
        "PublicAccessBlockConfiguration": {
          "BlockPublicAcls": true,
          "BlockPublicPolicy": true,
          "IgnorePublicAcls": true,
          "RestrictPublicBuckets": true
        }
      }
    }
  }
}
```

So without any effort at all a developer has satisfied most of the security controls in [AWS Foundational Security Best Practices controls](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html) for creating s3 bucket resources.

[S3.1](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html#fsbp-s3-1) S3 Block Public Access setting should be enabled
[S3.2](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html#fsbp-s3-2) S3 buckets should prohibit public read access
[S3.3](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html#fsbp-s3-3) S3 buckets should prohibit public write access
[S3.4](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html#fsbp-s3-4) S3 buckets should have server-side encryption enabled

# AWS Links and Security Standards

* [AWS Security Hub User Guide](https://docs.aws.amazon.com/securityhub/latest/userguide/what-is-securityhub.html)
* [AWS Foundational Security Best Practices controls](https://docs.aws.amazon.com/securityhub/latest/userguide/securityhub-standards-fsbp-controls.html)
* [CIS Amazon Web Services Foundations 1.2.0 (PDF)](https://d1.awsstatic.com/whitepapers/compliance/AWS_CIS_Foundations_Benchmark.pdf)
