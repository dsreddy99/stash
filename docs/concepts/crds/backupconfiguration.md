---
title: BackupConfiguration Overview
menu:
  product_stash_0.8.3:
    identifier: backupconfiguration-overview
    name: BackupConfiguration
    parent: crds
    weight: 15
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: concepts
---

> New to Stash? Please start [here](/docs/concepts/README.md).

# BackupConfiguration

## What is BackupConfiguration

A `BackupConfiguration` is a Kubernetes `CustomResourceDefinition (CRD)` which specifies the backup target, behaviors (schedule, retention policy etc.) and `Repository` object that holds backend information in Kubernetes native way.

You have to create a `BackupConfiguration` object for each backup target. `BackupConfiguration` has 1-1 mapping with the target. Thus, only one target can be backed up using one `BackupConfiguration`.

## BackupConfiguration CRD Specification

Like other official Kubernetes resources, `BackupConfiguration` has `TypeMeta`, `ObjectMeta` and `Spec` sections. However, unlike other Kubernetes resources, it does not have a `Status` section.

A sample `BackupConfiguration` object to backup a Deployment's data is shown below,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: BackupConfiguration
metadata:
  name: demo-backup
  namespace: demo
spec:
  repository:
    name: local-repo
  # task:
  #   name: workload-backup # task field is not required for workload data backup but it is necessary for database backup.
  schedule: "* * * * *" # backup at every minutes
  paused: false
  target:
    ref:
      apiVersion: apps/v1
      kind: Deployment
      name: stash-demo
    directories:
    - /source/data
    volumeMounts:
    - name: source-data
      mountPath: /source/data
  runtimeSettings:
    container:
      resources:
        requests:
          memory: 256M
        limits:
          memory: 256M
      securityContext:
        runAsUser: 2000
        runAsGroup: 2000
      nice:
        adjustment: 5
      ionice:
        class: 2
        classData: 4
    pod:
      imagePullSecrets:
      - name:  my-private-registry-secret
      serviceAccountName: my-backup-svc
  tempDir:
    medium: Memory
    sizeLimit:  2Gi
    disableCaching: false
  retentionPolicy:
    name: 'keep-last-5'
    keepLast: 5
    prune: true
```

Here, we are going to describe some important sections of `BackupConfiguration` crd.

### BackupConfiguration `Spec` Section

`BackupConfiguration` object holds following fields in `.spec` section.

#### spec.repository

`spec.repository.name` indicates the `Repository` crd name that holds necessary backend information where the backed up data will be stored.

#### spec.schedule

`spec.schedule` is a [cron expression](https://en.wikipedia.org/wiki/Cron) that specifies the schedule of backup. Stash creates a Kubernetes [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) with this schedule.

#### spec.task

`spec.task` specifies the name and parameters of [Task](/docs/concepts/crds/task.md) template to use to backup the target.

- **spec.task.name:** `spec.task.name` indicates the name of the `Task` template to use for this backup process.
- **spec.task.params:** `spec.task.params` is an array of custom parameters to use to configure the task.

> `spec.task` section is not necessary for backing up workload data (i.e. Deployment, DaemonSet, StatefulSet etc.). However, it is necessary for backing up databases and stand-alone PVC.

#### spec.paused

`spec.paused` can be used as `enable/disable` switch for backup. If it is set `true`, Stash will not take any backup of the target specified by this BackupConfiguration.

#### spec.target

`spec.target` field indicates the target of backup. This section consist of the following fields:

- **spec.target.ref :** `spec.target.ref` refers to the target of backup. You have to specify `apiVersion`, `kind` and `name` of the target. Stash will use this information to inject a sidecar to the target or to create a backup job for it.

- **spec.target.directories :** `spec.target.directories` specifies list of directories to backup.

- **spec.target.volumeMounts :** `spec.target.volumeMounts` list of volumes and their `mountPath` that contains the target directories. Stash will mount these volumes inside sidecar container or backup job.

#### spec.runtimeSettings

`spec.runtimeSettings` allows to configure runtime environment for backup sidecar or job. You can specify runtime settings in both pod level and container level.

- **spec.runtimeSettings.container**
  
  `spec.runtimeSettings.container` is used to configure backup sidecar/job in container level. You can configure the following container level parameters:

|       Field       |                                                                                                           Usage                                                                                                            |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `resources`       | Compute resources required by sidecar container or backup job. To know how to manage resources for containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/). |
| `livenessProbe`   | Periodic probe of backup sidecar/job container's liveness. Container will be restarted if the probe fails.                                                                                                                 |
| `readinessProbe`  | Periodic probe of backup sidecar/job container's readiness. Container will be removed from service endpoints if the probe fails.                                                                                           |
| `lifecycle`       | Actions that the management system should take in response to container lifecycle events.                                                                                                                                  |
| `securityContext` | Security options that backup sidecar/job's container should run with. For more details, please visit [here](https://kubernetes.io/docs/concepts/policy/security-context/).                                                 |
| `nice`            | Set CPU scheduling priority for backup process. For more details about `nice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#nice).                                                         |
| `ionice`          | Set I/O scheduling class and priority for backup process. For more details about `ionice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#ionice).                                           |

- **spec.runtimeSettings.pod**

  `spec.runtimeSettings.pod` is used to configure backup job in pod level. You can configure the following pod level parameters,

|             Field              |                                                                                                                  Usage                                                                                                                   |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `serviceAccountName`           | Name of the `ServiceAccount` to use for backup job. Stash sidecar will use the same `ServiceAccount` as the target.                                                                                                                      |
| `nodeSelector`                 | Selector which must be true for backup job pod to fit on a node.                                                                                                                                                                         |
| `automountServiceAccountToken` | Indicates whether a service account token should be automatically mounted into the backup pod.                                                                                                                                           |
| `nodeName`                     | NodeName is used to request to schedule backup job's pod onto a specific node.                                                                                                                                                           |
| `securityContext`              | Security options that backup job's pod should run with. For more details, please visit [here](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).                                                               |
| `imagePullSecrets`             | A list of secret names in the same namespace that will be used to pull image from private Docker registry. For more details, please visit [here](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/). |
| `affinity`                     | Affinity and anti-affinity to schedule backup job's pod in the desired node. For more details, please visit [here](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity).                       |
| `schedulerName`                | Name of the scheduler that should dispatch the backup job.                                                                                                                                                                               |
| `tolerations`                  | Taints and Tolerations to ensure that backup job's pod is not scheduled in inappropriate nodes. For more details about `toleration`, please visit [here](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/).       |
| `priorityClassName`            | Indicates the backup job pod's priority class. For more details, please visit [here](https://kubernetes.io/docs/concepts/configuration/pod-priority-preemption/).                                                                        |
| `priority`                     | Indicates the backup job pod's priority value.                                                                                                                                                                                           |
| `readinessGates`               | Specifies additional conditions to be evaluated for Pod readiness. For more details, please visit [here](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-readiness-gate).                                          |
| `runtimeClassName`             | RuntimeClass is used for selecting the container runtime configuration. For more details, please visit [here](https://kubernetes.io/docs/concepts/containers/runtime-class/)                                                             |
| `enableServiceLinks`           | EnableServiceLinks indicates whether information about services should be injected into pod's environment variables.                                                                                                                     |

#### spec.tempDir

Stash mounts an `emtpyDir` for holding temporary files. It is also used for `caching` for faster backup performance. You can configure the `emptyDir` using `spec.tempDir` section. You can also disable `caching` using this field. Following fields are configurable in `spec.tempDir` section:

- **spec.tempDir.medium :** Specifies the type of storage medium should back this directory.
- **spec.tempDir.sizeLimit :** Maximum limit of storage for this volume.
- **spec.tempDir.disableCaching :** Disable caching while backup. This may negatively impact backup performance.

#### spec.retentionPolicy

`spec.retentionPolicy` specifies the policy to follow for cleaning old snapshots. Following options are available to configure retention policy:

|    Policy     |  Value  | `restic` forget command flag |                                            Description                                             |
| ------------- | ------- | ---------------------------- | -------------------------------------------------------------------------------------------------- |
| `name`        | string  |                              | Name of retention policy. You can provide any name.                                                |
| `keepLast`    | integer | --keep-last n                | Never delete the **n** last (most recent) snapshots.                                               |
| `keepHourly`  | integer | --keep-hourly n              | For the last **n** hours in which a snapshot was made, keep only the last snapshot for each hour.  |
| `keepDaily`   | integer | --keep-daily n               | For the last **n** days which have one or more snapshots, only keep the last one for that day.     |
| `keepWeekly`  | integer | --keep-weekly n              | For the last **n** weeks which have one or more snapshots, only keep the last one for that week.   |
| `keepMonthly` | integer | --keep-monthly n             | For the last **n** months which have one or more snapshots, only keep the last one for that month. |
| `keepYearly`  | integer | --keep-yearly n              | For the last **n** years which have one or more snapshots, only keep the last one for that year.   |
| `keepTags`    | array   | --keep-tag <tag>             | Keep all snapshots which have all tags specified by this option (can be specified multiple times). |
| `prune`       | bool    | --prune                      | If set `true`, Stash will cleanup unreferenced data from the backend.                              |
| `dryRun`      | bool    | --dry-run                    | Stash will not remove anything but print which snapshots would be removed.                         |

## Next Steps

- Learn how to configure `BackupConfiguration` to backup workloads data from [here](/docs/guides/workloads/backup.md).
- Learn how to configure `BackupConfiguration` to backup databases from [here](/docs/guides/databases/backup.md).
- Learn how to configure `BackupConfiguration` to backup stand-alone PVC from [here](/docs/guides/volumes/backup.md).
