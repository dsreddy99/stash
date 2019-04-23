---
title: BackupConfigurationTemplate Overview
menu:
  product_stash_0.8.3:
    identifier: backuptemplate-overview
    name: BackupConfigurationTemplate
    parent: crds
    weight: 40
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: concepts
---

> New to Stash? Please start [here](/docs/concepts/README.md).

# BackupConfigurationTemplate

## What is BackupConfigurationTemplate

Stash uses 1:1:1 mapping between `Repository`, `BackupConfiguration` and the target. So, whenever you want to backup a target, you have to create a `Repository` and `BackupConfiguration` object. This could be tiresome job when you are trying to backup similar types of target and `Repository` and `BackupConfiguration` has only slight difference. To mitigate this problem, Stash provides a way to specify a template for those two objects through `BackupConfigurationTemplate` crd.

A `BackupConfigurationTemplate` is a Kubernetes `CustomResourceDefinition (CRD)` which specifies a template for `Repository` and `BackupConfiguration` for similar types of backup target in Kubernetes native way.

A user have to create only one  `BackupConfigurationTemplate` for all similar types of workload (i.e. Deployment, DaemonSet, StatefulSet etc.). Then, she just needs to add some annotations in the target workload. Stash will automatically create respective `Repository`, `BackupConfiguration` object using the template. In Stash parlance, we call this process as **default backup**.

## BackupConfigurationTemplate CRD Specification

Like other official Kubernetes resources, `BackupConfigurationTemplate` has `TypeMeta`, `ObjectMeta` and `Spec` sections. However, unlike other Kubernetes resources, it does not have a `Status` section.

A sample `BackupConfigurationTemplate` object to backup a Deployment's data through default backup is shown below,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: BackupConfigurationTemplate
metadata:
  name: workload-backup-template
spec:
  # ============== Template for Repository ==========================
  backend:
    gcs:
      bucket: stash-backup
      prefix: stash/${TARGET_NAMESPACE}/${TARGET_KIND}/${TARGET_NAME}
    storageSecretName: gcs-secret
  wipeOut: false
  # ============== Template for BackupConfiguration =================
  schedule: "* * * * *"
  # task: # no task section is required for workload data backup
  #   name: workload-backup
  runtimeSettings:
    container:
      securityContext:
        runAsUser: 2000
        runAsGroup: 2000
  tempDir:
    medium: "Memory"
    size:  "1Gi"
    disableCache: false
  retentionPolicy:
    name: 'keep-last-5'
    keepLast: 5
    prune: true
```

The sample `BackupConfigurationTemplate` that has been shown above can be used to backup Deployments, DaemonSets, StatefulSets, ReplicaSets and ReplicationControllers. You only have to add some annotations to these workloads. For more details about what annotations you have to add to the target, please visit [here](/docs/guides/default-backup/overview.md).

Here, we are going to describe some important sections of `BackupConfigurationTemplate` crd.

### BackupConfigurationTemplate `Spec` Section

We can divide BackupConfigurationTemplate's `.spec` section into two parts. One part specifies template for `Repository` object and other specifies template for `BackupConfiguration` object.

#### Repository Template

You can configure `Repository` template using `spec.backend` field and `spec.wipeOut` field.

- **spec.backend :** `spec.backend` section is backend specification similar to [spec.backend](/docs/concepts/crds/repository.md#specbackend) section of `Repository` crd. There is only one difference. You can now template `prefix` section (`subPath` for local volume) of the backend to store backed up data of different workloads at different directory of the backend. You can use following variables to template `spec.backend` section:
    |       Variable       |            Usage            |
    | -------------------- | --------------------------- |
    | `TARGET_API_VERSION` | API version of the target   |
    | `TARGET_KIND`        | Resource kind of the target |
    | `TARGET_NAMESPACE`   | Namespace of the target     |
    | `TARGET_NAME`        | Name of the target          |

    If we use the sample `BackupConfigurationTemplate` that has been shown above to backup a Deployment named `my-deploy` of `test` namespace, the backed up file will be stored in `stash/test/deployment/my-deploy` directory of the `stash-backup` bucket. If we want to backup a ReplicaSet with name `my-rs` of same namespace, the backed up data will be stored in `/stash/test/replicaset/my-rs` directory of the backend.

- **spec.wipeOut :** `spec.wipeOut` indicates whether Stash should delete backed up data from backend if user delete respective `Repository` created by Stash for a target. For more details, please visit [here](/docs/concepts/crds/repository.md#specwipeout).

#### BackupConfiguration Template

You can set a template for the `BackupConfiguration` object that will be created for respective target using following fields:

- **spec.schedule :** `spec.schedule` is the schedule that will be used to create `BackupConfiguration` for respective target. For more details, please visit [here](/docs/concepts/crds/backupconfiguration.md#specschedule).

- **spec.task :** `spec.task` specifies the name and parameters of [Task](/docs/concepts/crds/task.md) template to use to backup the target.

- **spec.runtimeSettings :** `spec.runtimeSettings` allows to configure runtime environment for backup sidecar or job. For more details, please visit [here](/docs/concepts/crds/backupconfiguration.md#specruntimesettings).

- **spec.tempDir :** `spec.tempDir` specifies the temporary volume setting that will be used to create respective `BackupConfiguration` object. For more details, please visit [here](/docs/concepts/crds/backupconfiguration.md#spectempdir).

- **spec.retentionPolicy :** `spec.retentionPolicy` specifies the retention policies that will be used to create respective `BackupConfiguration` object. For more details, please visit [here](/docs/concepts/crds/backupconfiguration.md#specretentionpolicy).

## Next Steps

- Learn how to use `BackupConfigurationTemplate` for default backup of workloads data from [here](/docs/guides/default-backup/workload.md).
- Learn how to use `BackupConfigurationTemplate` for default-backup of database from [here](/docs/guides/default-backup/database.md).
- Learn how to use `BackupConfigurationTemplate` for default-backup of stand alone PVC from [here](/docs/guides/default-backup/volume.md).