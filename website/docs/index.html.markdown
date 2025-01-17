---
layout: "yandex"
page_title: "Provider: Yandex.Cloud"
sidebar_current: "docs-yandex-index"
description: |-
  The Yandex.Cloud provider is used to interact with Yandex.Cloud services.
  The provider needs to be configured with the proper credentials before it can be used.
---

# Yandex.Cloud Provider

The Yandex.Cloud provider is used to interact with
[Yandex.Cloud services](https://cloud.yandex.com/). The provider needs
to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the Yandex.Cloud provider
provider "yandex" {
  token                    = "auth_token_here"
  service_account_key_file = "path_to_service_account_key_file"
  cloud_id                 = "cloud_id_here"
  folder_id                = "folder_id_here"
  zone                     = "ru-central1-a"
}

// Create a new instance
resource "yandex_compute_instance" "default" {
  ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `token` - (Optional) Security token used for authentication in Yandex.Cloud.

  This can also be specified using environment variable `YC_TOKEN`.

* `service_account_key_file` - (Optional) Path to file that contains service account key data.

  This can also be specified using environment variable `YC_SERVICE_ACCOUNT_KEY_FILE`.
  You can read how to create service account key file [here][yandex-service-account-key].

~> **NOTE:** Only one of `token` or `service_account_key_file` can be specified.

* `cloud_id` - (Required) The ID of the [cloud][yandex-cloud] to apply any resources to.

  This can also be specified using environment variable `YC_CLOUD_ID`.

* `folder_id` - (Required) The ID of the [folder][yandex-folder] to operate under, if not specified by a given resource.

  This can also be specified using environment variable `YC_FOLDER_ID`.

* `zone` - (Optional) The default [availability zone][yandex-zone] to operate under, if not specified by a given resource.

  This can also be specified using environment variable `YC_ZONE`.
  
* `max_retries` - (Optional) This is the maximum number of times an API call is retried, in the case where requests 
  
  are being throttled or experiencing transient failures. The delay between the subsequent API calls increases 
  exponentially.


[yandex-cloud]: https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy#cloud
[yandex-folder]: https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy#folder
[yandex-zone]: https://cloud.yandex.com/docs/overview/concepts/geo-scope
[yandex-service-account-key]: https://cloud.yandex.com/docs/iam/operations/iam-token/create-for-sa#keys-create
