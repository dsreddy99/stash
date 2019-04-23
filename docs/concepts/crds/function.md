---
title: Function Overview
menu:
  product_stash_0.8.3:
    identifier: function-overview
    name: Function
    parent: crds
    weight: 30
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: concepts
---

> New to Stash? Please start [here](/docs/concepts/README.md).

# Function

## What is Function

A complete backup or restore process may consist of several steps. For example, in order to backup a PostgreSQL database we first need to dump the database then upload the dumped file to the backend then we need to update Repository and BackupSession status and send Prometheus metrics. In Stash, we call such individual step as a `Function`.

A `Function` is a Kubernetes `CustomResourceDefinition (CRD)` which basically specifies a template for a container that perform only a specific action. For example, `pg-backup` function only dump and upload dumped file into the backend where `update-status` function update status of respective `BackupSession` and `Repository` and send Prometheus metrics to pushgateway based on the output of `pg-backup` function.

When you install Stash, some default `Functions` will be automatically created for supported targets. However, you can create your own function to customize or extend backup/restore process.

## Function CRD Specification

Like other official Kubernetes resources, `Function` has `TypeMeta`, `ObjectMeta` and `Spec` sections. However, unlike other Kubernetes resources, it does not have a `Status` section.

A sample `Function` object to backup a PostgreSQL is shown below,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: Function
metadata:
  name: pg-backup
spec:
  image: appscode/stash:pg
  args:
  - backup-pg
  - --provider=${REPOSITORY_PROVIDER:=}
  - --bucket=${REPOSITORY_BUCKET:=}
  - --endpoint=${REPOSITORY_ENDPOINT:=}
  - --path=${REPOSITORY_PREFIX:=}
  - --secret-dir=/etc/repository/secret
  - --scratch-dir=/tmp
  - --hostname=${HOSTNAME:=host-0}
  - --pg-args=${pgArgs:=}
  - --namespace=${NAMESPACE:=default}
  - --app-binding=${TARGET_NAME:=}
  - --retention-keep-last=${RETENTION_KEEP_LAST:=0}
  - --retention-prune=${RETENTION_PRUNE:=false}
  - --output-dir=${outputDir:=}
  - --enable-cache=${ENABLE_CACHE:=true}
  - --max-connections=${MAX_CONNECTIONS:=0}
  volumeMounts:
  - name: ${secretVolume}
    mountPath: /etc/repository/secret
  runtimeSettings:
    container:
      resources:
        requests:
          memory: 256M
        limits:
          memory: 256M
      securityContext:
        runAsUser: 5000
        runAsGroup: 5000
```

A sample Function that update `BackupSession` and `Repository`  status and send metrics to Prometheus is shown below,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: Function
metadata:
  name: update-status
spec:
  image: appscode/stash:pg
  args:
  - update-status
  - --namespace=${NAMESPACE:=default}
  - --repository=${REPOSITORY_NAME:=}
  - --backup-session=${BACKUP_SESSION:=}
  - --restore-session=${RESTORE_SESSION:=}
  - --output-dir=${outputDir:=}
  - --enable-status-subresource=${ENABLE_STATUS_SUBRESOURCE:=false}
```

Here, we are going to describe some important sections of a `Function` crd.

### Function `Spec` Section

Function object holds following fields in `.spec` section.

#### spec.image

`spec.image` specifies docker image to use to create container using template specified in this Function.

#### spec.command

`spec.command` specifies the commands to be executed by the container. Docker image's `ENTRYPOINT` will be executed if no commands are specified.

#### spec.args

`spec.args` specifies a list of arguments that will be passed to the entrypoint. You can template this section using variable. Stash will resolve all variables before creating respective container. A variable should following patterns:

- ${VARIABLE_NAME:=default-value}
- ${VARIABLE_NAME:=}

In the first case, if Stash can't resolve the variable, default value will be used in place of this variable. In the second case, if Stash can't resolve the variable, empty string will be used to replace the variable.

##### Stash Provided Variables

Stash can provide values of following variables utilizing `BackupConfiguration`, `BackupSession`, `RestoreSession`, `Repository`, `Task`, `Function`, `BackupConfigurationTemplate` etc.

|    Environment Variable     |                                                                                      Usage                                                                                       |
| :-------------------------: | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|         `NAMESPACE`         | Namespace of backup or restore job                                                                                                                                               |
|      `BACKUP_SESSION`       | Name of the respective BackupSession object                                                                                                                                      |
|      `RESTORE_SESSION`      | Name of the respective RestoreSession object                                                                                                                                     |
|      `REPOSITORY_NAME`      | Name of the Repository object that holds respective backend information                                                                                                          |
|    `REPOSITORY_PROVIDER`    | Type of storage provider. i.e. gcs, s3, aws, local etc.                                                                                                                          |
|  `REPOSITORY_SECRET_NAME`   | Name of the secret that holds the credentials to access the backend                                                                                                              |
|     `REPOSITORY_BUCKET`     | Name of the bucket where backed up data will be stored                                                                                                                           |
|     `REPOSITORY_PREFIX`     | A prefix of directory inside bucket where backed up data will be stored                                                                                                          |
|    `REPOSITORY_ENDPOINT`    | URL of S3 compatible Minio/Rook server                                                                                                                                           |
|      `REPOSITORY_URL`       | URL of the REST server for REST backend                                                                                                                                          |
|         `HOSTNAME`          | An identifier for the backed up data. If multiple pods backup in same Repository (i.e. StatefulSet or DaemonSet) this host name is to used identify data of the individual host. |
|        `TARGET_NAME`        | Name of the target of backup or restore                                                                                                                                          |
|    `TARGET_API_VERSION`     | API version of the target of backup or restore                                                                                                                                   |
|        `TARGET_KIND`        | Kind of the target of backup or restore                                                                                                                                          |
|     `TARGET_NAMESPACE`      | Namespace of target object for backup or restore                                                                                                                                 |
|     `TARGET_MOUNT_PATH`     | Directory where target PVC will be mounted in stand alone PVC backup or restore                                                                                                  |
|    `TARGET_DIRECTORIES`     | Array of directories that are subject to backup                                                                                                                                  |
|    `RESTORE_DIRECTORIES`    | Array of directories that are subject to restore                                                                                                                                 |
|     `RESTORE_SNAPSHOTS`     | Name of the snapshot that will be restored                                                                                                                                       |
|    `RETENTION_KEEP_LAST`    | Number of latest snapshot to keep                                                                                                                                                |
|   `RETENTION_KEEP_HOURLY`   | Number of hourly snapshot to keep                                                                                                                                                |
|   `RETENTION_KEEP_DAILY`    | Number of daily snapshot to keep                                                                                                                                                 |
|   `RETENTION_KEEP_WEEKLY`   | Number of weekly snapshot to keep                                                                                                                                                |
|  `RETENTION_KEEP_MONTHLY`   | Number of monthly snapshot to keep                                                                                                                                               |
|   `RETENTION_KEEP_YEARLY`   | Number of yearly snapshot to keep                                                                                                                                                |
|    `RETENTION_KEEP_TAGS`    | Keep only those snapshot that has these tags                                                                                                                                     |
|      `RETENTION_PRUNE`      | Specify weather to remove data of old snapshot completely from the backend                                                                                                       |
|     `RETENTION_DRY_RUN`     | Specify weather to run cleanup in test mode                                                                                                                                      |
|       `ENABLE_CACHE`        | Specify weather to use cache while backup or restore                                                                                                                             |
|      `MAX_CONNECTIONS`      | Specifies number of parallel connections to upload/download data to/from backend                                                                                                 |
|      `NICE_ADJUSTMENT`      | Adjustment value to configure `nice` to throttle the load on cpu.                                                                                                                |
|       `IONICE_CLASS`        | Name of the `ionice` class                                                                                                                                                       |
|     `IONICE_CLASS_DATA`     | Value of the `ionice` class data                                                                                                                                                 |
| `ENABLE_STATUS_SUBRESOURCE` | Specifies weather crd has subresource enabled                                                                                                                                    |

If you want to use a variable that is not present this table, you have to provide it's value in `spec.task.params` section of `BackupConfiguration` crd.

#### spec.workDir

`spec.workDir` specifies the container's working directory. If this field is not specified, container's runtime default will be used might be configured in respective docker file.

#### spec.ports

`spec.ports` specifies list of ports to expose from the respective container created for this function.

#### spec.env

`spec.env` section specifies list of environment variable to set in the container created for this function.

#### spec.envFrom

`spec.envFrom` allows to set environment variables to the container created for this function from a Secret or ConfigMap.

#### spec.volumeMounts

`spec.volumeMounts` specifies a list of volume names and their `mountPath` that will be mounted into the container created for this function.

#### spec.volumeDevices

`spec.volumeDevices` is the list of block devices to be used by the container created for this function.

#### spec.runtimeSettings

`spec.runtimeSettings.container` allows to configure runtime environment of backup job in container level. You can configure following container level parameters,

  |       Field       |                                                                                                           Usage                                                                                                            |
  | :---------------: | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
  |    `resources`    | Compute resources required by sidecar container or backup job. To know how to manage resources for containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/). |
  |  `livenessProbe`  | Periodic probe of backup sidecar/job container's liveness. Container will be restarted if the probe fails.                                                                                                                 |
  | `readinessProbe`  | Periodic probe of backup sidecar/job container's readiness. Container will be removed from service endpoints if the probe fails.                                                                                           |
  |    `lifecycle`    | Actions that the management system should take in response to container lifecycle events.                                                                                                                                  |
  | `securityContext` | Security options that backup sidecar/job's container should run with. For more details, please visit [here](https://kubernetes.io/docs/concepts/policy/security-context/).                                                 |
  |      `nice`       | Set CPU scheduling priority for backup process. For more details about `nice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#nice).                                                         |
  |     `ionice`      | Set I/O scheduling class and priority for backup process. For more details about `ionice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#ionice).                                           |

#### spec.podSecurityPolicyName

If you are using a [PSP enabled cluster](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) and the function need any specific permission then you can specify the PSP name using `spec.podSecurityPolicyName` field. Stash will add this PSP in the respective RBAC roles created for this function.

>Note that, Stash operator can't give permission to use a PSP to a backup job if operator itself does not have permission to use it. So, if you want to specify PSP name in this section, make sure to add this that to `stash-operator` ClusterRole too. For more details about using PSP in Stash, please visit [here](/docs/setup/psp.md).

## Next Steps

- Learn how to use `Function` to create a `Task` from [here](/docs/concepts/crds/task.md).
