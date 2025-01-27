---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "AWS: aws_ami_copy"
description: |-
  Duplicates an existing Amazon Machine Image (AMI)
---

# Resource: aws_ami_copy

The "AMI copy" resource allows duplication of an Amazon Machine Image (AMI),
including cross-region copies.

If the source AMI has associated EBS snapshots, those will also be duplicated
along with the AMI.

This is useful for taking a single AMI provisioned in one region and making
it available in another for a multi-region deployment.

Copying an AMI can take several minutes. The creation of this resource will
block until the new AMI is available for use on new instances.

## Example Usage

```terraform
resource "aws_ami_copy" "example" {
  name              = "terraform-example"
  description       = "A copy of ami-xxxxxxxx"
  source_ami_id     = "ami-xxxxxxxx"
  source_ami_region = "us-west-1"

  tags = {
    Name = "HelloWorld"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A region-unique name for the AMI.
* `source_ami_id` - (Required) The id of the AMI to copy. This id must be valid in the region
  given by `source_ami_region`.
* `source_ami_region` - (Required) The region from which the AMI will be copied. This may be the
  same as the AWS provider region in order to create a copy within the same region.
* `destination_outpost_arn` - (Optional) The ARN of the Outpost to which to copy the AMI.
  Only specify this parameter when copying an AMI from an AWS Region to an Outpost. The AMI must be in the Region of the destination Outpost.  
* `encrypted` - (Optional) Specifies whether the destination snapshots of the copied image should be encrypted. Defaults to `false`
* `kms_key_id` - (Optional) The full ARN of the KMS Key to use when encrypting the snapshots of an image during a copy operation. If not specified, then the default AWS KMS Key will be used
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

This resource also exposes the full set of arguments from the [`aws_ami`](ami.html) resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 40 mins) Used when creating the AMI
* `update` - (Defaults to 40 mins) Used when updating the AMI
* `delete` - (Defaults to 90 mins) Used when deregistering the AMI

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the AMI.
* `id` - The ID of the created AMI.

This resource also exports a full set of attributes corresponding to the arguments of the
[`aws_ami`](/docs/providers/aws/r/ami.html) resource, allowing the properties of the created AMI to be used elsewhere in the
configuration.
