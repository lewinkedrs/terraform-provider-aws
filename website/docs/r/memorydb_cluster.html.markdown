---
subcategory: "MemoryDB for Redis"
layout: "aws"
page_title: "AWS: aws_memorydb_cluster"
description: |-
  Provides a MemoryDB Cluster.
---

# Resource: aws_memorydb_cluster

Provides a MemoryDB Cluster.

More information about MemoryDB can be found in the [Developer Guide](https://docs.aws.amazon.com/memorydb/latest/devguide/what-is-memorydb-for-redis.html).

## Example Usage

```terraform
resource "aws_memorydb_cluster" "example" {
  acl_name                 = "open-access"
  name                     = "my-cluster"
  node_type                = "db.t4g.small"
  num_shards               = 2
  security_group_ids       = [aws_security_group.example.id]
  snapshot_retention_limit = 7
  subnet_group_name        = aws_memorydb_subnet_group.example.id
}
```

## Argument Reference

The following arguments are required:

* `acl_name` - (Required) The name of the Access Control List to associate with the cluster.
* `node_type` - (Required) The compute and memory capacity of the nodes in the cluster. See AWS documentation on [supported node types](https://docs.aws.amazon.com/memorydb/latest/devguide/nodes.supportedtypes.html) as well as [vertical scaling](https://docs.aws.amazon.com/memorydb/latest/devguide/cluster-vertical-scaling.html).

The following arguments are optional:

* `auto_minor_version_upgrade` - (Optional, Forces new resource) When set to `true`, the cluster will automatically receive minor engine version upgrades after launch. Defaults to `true`.
* `description` - (Optional) Description for the cluster. Defaults to `"Managed by Terraform"`.
* `engine_version` - (Optional) Version number of the Redis engine to be used for the cluster. Downgrades are not supported.
* `final_snapshot_name` - (Optional) Name of the final cluster snapshot to be created when this resource is deleted. If omitted, no final snapshot will be made.
* `kms_key_arn` - (Optional, Forces new resource) ARN of the KMS key used to encrypt the cluster at rest.
* `maintenance_window` - (Optional) Specifies the weekly time range during which maintenance on the cluster is performed. Specify as a range in the format `ddd:hh24:mi-ddd:hh24:mi` (24H Clock UTC). The minimum maintenance window is a 60 minute period. Example: `sun:23:00-mon:01:30`.
* `name` - (Optional, Forces new resource) Name of the cluster. If omitted, Terraform will assign a random, unique name. Conflicts with `name_prefix`.
* `name_prefix` - (Optional, Forces new resource) Creates a unique name beginning with the specified prefix. Conflicts with `name`.
* `num_replicas_per_shard` - (Optional) The number of replicas to apply to each shard, up to a maximum of 5. Defaults to `1` (i.e. 2 nodes per shard).
* `num_shards` - (Optional) The number of shards in the cluster. Defaults to `1`.
* `parameter_group_name` - (Optional) The name of the parameter group associated with the cluster.
* `port` - (Optional, Forces new resource) The port number on which each of the nodes accepts connections. Defaults to `6379`.
* `security_group_ids` - (Optional) Set of VPC Security Group ID-s to associate with this cluster.
* `snapshot_arns` - (Optional, Forces new resource) List of ARN-s that uniquely identify RDB snapshot files stored in S3. The snapshot files will be used to populate the new cluster. Object names in the ARN-s cannot contain any commas.
* `snapshot_name` - (Optional, Forces new resource) The name of a snapshot from which to restore data into the new cluster.
* `snapshot_retention_limit` - (Optional) The number of days for which MemoryDB retains automatic snapshots before deleting them. When set to `0`, automatic backups are disabled. Defaults to `0`.
* `snapshot_window` - (Optional) The daily time range (in UTC) during which MemoryDB begins taking a daily snapshot of your shard. Example: `05:00-09:00`.
* `sns_topic_arn` - (Optional) ARN of the SNS topic to which cluster notifications are sent.
* `subnet_group_name` - (Optional, Forces new resource) The name of the subnet group to be used for the cluster. Defaults to a subnet group consisting of default VPC subnets.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `tls_enabled` - (Optional, Forces new resource) A flag to enable in-transit encryption on the cluster. When set to `false`, the `acl_name` must be `open-access`. Defaults to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Same as `name`.
* `arn` - The ARN of the cluster.
* `cluster_endpoint`
    * `address` - DNS hostname of the cluster configuration endpoint.
    * `port` - Port number that the cluster configuration endpoint is listening on.
* `engine_patch_version` - Patch version number of the Redis engine used by the cluster.
* `shards` - Set of shards in this cluster.
    * `name` - Name of this shard.
    * `num_nodes` - Number of individual nodes in this shard.
    * `slots` - Keyspace for this shard. Example: `0-16383`.
    * `nodes` - Set of nodes in this shard.
        * `availability_zone` - The Availability Zone in which the node resides.
        * `create_time` - The date and time when the node was created. Example: `2022-01-01T21:00:00Z`.
        * `name` - Name of this node.
        * `endpoint`
            * `address` - DNS hostname of the node.
            * `port` - Port number that this node is listening on.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Timeouts

`aws_memorydb_cluster` provides the following [timeout configuration options](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts):

- `create` - (Default `120 minutes`) Used when creating a cluster.
- `update` - (Default `120 minutes`) Used when updating a cluster.
- `delete` - (Default `120 minutes`) Used when deleting a cluster.

## Import

Use the `name` to import a cluster. For example:

```
$ terraform import aws_memorydb_cluster.example my-cluster
```
