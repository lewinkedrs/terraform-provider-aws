---
subcategory: "Transit Gateway"
layout: "aws"
page_title: "AWS: aws_ec2_transit_gateway"
description: |-
  Manages an EC2 Transit Gateway
---

# Resource: aws_ec2_transit_gateway

Manages an EC2 Transit Gateway.

## Example Usage

```terraform
resource "aws_ec2_transit_gateway" "example" {
  description = "example"
}
```

## Argument Reference

The following arguments are supported:

* `amazon_side_asn` - (Optional) Private Autonomous System Number (ASN) for the Amazon side of a BGP session. The range is `64512` to `65534` for 16-bit ASNs and `4200000000` to `4294967294` for 32-bit ASNs. Default value: `64512`.
* `auto_accept_shared_attachments` - (Optional) Whether resource attachment requests are automatically accepted. Valid values: `disable`, `enable`. Default value: `disable`.
* `default_route_table_association` - (Optional) Whether resource attachments are automatically associated with the default association route table. Valid values: `disable`, `enable`. Default value: `enable`.
* `default_route_table_propagation` - (Optional) Whether resource attachments automatically propagate routes to the default propagation route table. Valid values: `disable`, `enable`. Default value: `enable`.
* `description` - (Optional) Description of the EC2 Transit Gateway.
* `dns_support` - (Optional) Whether DNS support is enabled. Valid values: `disable`, `enable`. Default value: `enable`.
* `multicast_support` - (Optional) Whether Multicast support is enabled. Required to use `ec2_transit_gateway_multicast_domain`. Valid values: `disable`, `enable`. Default value: `disable`.
* `tags` - (Optional) Key-value tags for the EC2 Transit Gateway. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `transit_gateway_cidr_blocks` - (Optional) One or more IPv4 or IPv6 CIDR blocks for the transit gateway. Must be a size /24 CIDR block or larger for IPv4, or a size /64 CIDR block or larger for IPv6.
* `vpn_ecmp_support` - (Optional) Whether VPN Equal Cost Multipath Protocol support is enabled. Valid values: `disable`, `enable`. Default value: `enable`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - EC2 Transit Gateway Amazon Resource Name (ARN)
* `association_default_route_table_id` - Identifier of the default association route table
* `id` - EC2 Transit Gateway identifier
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).
* `owner_id` - Identifier of the AWS account that owns the EC2 Transit Gateway
* `propagation_default_route_table_id` - Identifier of the default propagation route table

## Timeouts

`aws_ec2_transit_gateway` provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for EC2 Transit Gateway creation
- `update` - (Default `10 minutes`) Used for EC2 Transit Gateway updates
- `delete` - (Default `10 minutes`) Used for EC2 Transit Gateway deletion

## Import

`aws_ec2_transit_gateway` can be imported by using the EC2 Transit Gateway identifier, e.g.,

```
$ terraform import aws_ec2_transit_gateway.example tgw-12345678
```
