---
layout: "yandex"
page_title: "Yandex: yandex_kubernetes_cluster"
sidebar_current: "docs-yandex-datasource-kubernetes-cluster"
description: |-
  Get information about a Yandex Kubernetes Cluster. For more information, see
  [the official documentation](https://cloud.yandex.com/docs/managed-kubernetes/concepts/#kubernetes-cluster).
---

# yandex\_kubernetes\_cluster

Get information about a Yandex Kubernetes Cluster.

## Example Usage

```hcl
data "yandex_kubernetes_cluster" "my_cluster" {
  cluster_id = "some_k8s_cluster_id"
}

output "cluster_external_v4_endpoint" {
  value = "${data.yandex_kubernetes_cluster.my_cluster.master.0.external_v4_endpoint}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Optional) ID of a specific Kubernetes cluster.

* `name` - (Optional) Name of a specific Kubernetes cluster.

~> **NOTE:** One of `cluster_id` or `name` should be specified.

## Attributes Reference

* `description` - A description of the Kubernetes cluster.
* `folder_id` - The ID of the folder that the Kubernetes cluster belongs to.
* `labels` - A set of key/value label pairs to assign to the Kubernetes cluster.
* `network_id` - The ID of the cluster network.

* `service_account_id` - Service account to be used for provisioning Compute Cloud and VPC resources
for Kubernetes cluster. Selected service account should have `edit` role on the folder where the Kubernetes
cluster will be located and on the folder where selected network resides.

* `node_service_account_id` - Service account to be used by the worker nodes of the Kubernetes cluster
to access Container Registry or to push node logs and metrics.

* `release_channel` - Cluster release channel.

* `master` - IP allocation policy of the Kubernetes cluster.

The structure is documented below.

* `created_at` - The Kubernetes cluster creation timestamp.
* `status` - Status of the Kubernetes cluster.
* `health` - Health of the Kubernetes cluster.

---

The `master` block supports:

* `zonal` - Information about cluster zonal master.

The structure is documented below.

* `regional` - Information about cluster zonal master.

The structure is documented below.

* `internal_v4_address` - An IPv4 internal network address that is assigned to the master.
* `external_v4_address` - An IPv4 external network address that is assigned to the master.
* `internal_v4_endpoint` - Internal endpoint that can be used to connect to the master from cloud networks. 
* `external_v4_endpoint` - External endpoint that can be used to access Kubernetes cluster API from the internet (outside of the cloud).
* `cluster_ca_certificate` - PEM-encoded public certificate that is the root of trust for the Kubernetes cluster.  

* `version_info` - Information about cluster version.

The structure is documented below.

---

The `zonal` block supports:

* `zone` - ID of the availability zone where the master resides. 
---

The `regional` block supports:

* `region` - ID of the availability region where the master resides. 
---

The `version_info` block supports:

* `current_version` - Current Kubernetes version, major.minor (e.g. 1.15).
* `new_revision_available` - True/false flag.
Newer revisions may include Kubernetes patches (e.g 1.15.1 -> 1.15.2) as well
as some internal component updates - new features or bug fixes in yandex-specific
components either on the master or nodes.

* `new_revision_summary` - Human readable description of the changes to be applied
when updating to the latest revision. Empty if new_revision_available is false.
* `version_deprecated` - True/false flag. The current version is on the deprecation schedule,
component (master or node group) should be upgraded.


---