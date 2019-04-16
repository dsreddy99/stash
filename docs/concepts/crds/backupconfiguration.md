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

A `BackupConfiguration` is a Kubernetes `CustomResourceDefinition (CRD)` which specifies the backup target, behaviours (schedule, retention policy etc.) and name of the `Repository` object that holds backend information in Kubernetes native way.

Users have to create a `BackupConfiguration` object for each backup target. `BackupConfiguration` object has 1-1 mapping with target. Thus, only one target can be backed up using one `BackupConfiguration`.

## BackupConfiguration CRD Specification

Like other official Kubernetes resources, `BackupConfiguration` object has `TypeMeta`, `ObjectMeta` and `Spec` sections. However, unlike other Kubernetes resources, it does not have a `Status` section.

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
  # task: workload-backup # task field is not required for workload data backup but it is necessary for database backup.
  schedule: "*1 * * * *" # backup at every minutes
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

Here, we are going to describe some important sections of `BackupConfiguration` CRD.

### BackupConfiguration `Spec` Section

BackupConfiguration object holds following fields in `.spec` section.

#### spec.repository

`spec.repository.name` indicates the `Repository` crd name that hold necessary backend information where backed up data will be stored.

#### spec.schedule

`spec.schedule` is a [cron expression](https://en.wikipedia.org/wiki/Cron) that specifies the schedule of backup. Stash creates a Kubernetes [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) with this schedule.

#### spec.task

`spec.task` specifies the name and parameter of `Task` template to use to backup target.

- **spec.task.name:** `spec.task.name` indicates the name of the `Task` template to use for this backup process.
- **spec.task.params:** `spec.task.params` is an array of custom parameters to use to configure the task.

> `spec.task` section is not necessary for backing up workload data (i.e. Deployment, DaemonSet, StatefulSet etc.). However, it is necessary to backup database and stand alone PVC.

#### spec.paused

`spec.paused` can be used as `enable/disable` switch for backup. If it is set `true`, Stash will not take any backup of the target specified by this BackupConfiguration.

#### spec.target

`spec.target` field indicates the target of backup. This section consist of following fields.

- **spec.target.ref**
`spec.target.ref` refers to target of backup. Users have to specify `apiVersion`, `kind` and `name` of the target. Stash will use this information to inject sidecar or create backup job for respective target.

- **spec.target.directories**
`spec.target.directories` specifies list of target directories of backup.

- **spec.target.volumeMounts**
`spec.target.volumeMounts` list of volumes that contains these directories. Stash will mount these directories inside sidecar container or backup job.

#### spec.runtimeSettings

`spec.runtimeSettings` allows to configure runtime environment for backup sidecar or job. You can specify runtime settings in both pod level and container level.

- **spec.runtimeSettings.container** 
  `spec.runtimeSettings.container` is used to configure backup sidecar/job in container level. You can configure following container level parameters,
  
  |       Field       |                                                                                                           Usage                                                                                                            |
  | :---------------: | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
  |    `resources`    | Compute resources required by sidecar container or backup job. To know how to manage resources for containers, please visit [here](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/). |
  |  `livenessProbe`  | Periodic probe of backup sidecar/job container's liveness. Container will be restarted if the probe fails.                                                                                                                 |
  | `readinessProbe`  | Periodic probe of backup sidecar/job container's readiness. Container will be removed from service endpoints if the probe fails.                                                                                           |
  |    `lifecycle`    | Actions that the management system should take in response to container lifecycle events.                                                                                                                                  |
  | `securityContext` | Security options for container it should run with. For more details, please visit [here](https://kubernetes.io/docs/concepts/policy/security-context/).                                                                    |
  |      `nice`       | Set CPU scheduling priority for backup process. For more details about `nice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#nice).                                                         |
  |     `ionice`      | Set I/O scheduling class and priority for backup process. For more details about `ionice`, please visit [here](https://www.askapache.com/optimize/optimize-nice-ionice/#ionice).                                           |

- **spec.runtimeSettings.pod**

  `spec.runtimeSettings.pod` is used to configure backup sidecar/job in pod level. You can configure following pod level parameters,
  | Field  |  Usage |
  |---|---|
  |   |   |


#### spec.tempDir

#### spec.retentionPolicy

`spec.retentionPolicies` defines an array of retention policies for old snapshots. Retention policy options are below.

| Policy        | Value   | restic forget flag | Description                                                                                        |
|---------------|---------|--------------------|----------------------------------------------------------------------------------------------------|
| `name`        | string  |                    | Name of retention policy provided by users. This is used in file groups to refer to a policy.       |
| `keepLast`    | integer | --keep-last n      | Never delete the n last (most recent) snapshots                                                    |
| `keepHourly`  | integer | --keep-hourly n    | For the last n hours in which a snapshot was made, keep only the last snapshot for each hour.      |
| `keepDaily`   | integer | --keep-daily n     | For the last n days which have one or more snapshots, only keep the last one for that day.         |
| `keepWeekly`  | integer | --keep-weekly n    | For the last n weeks which have one or more snapshots, only keep the last one for that week.       |
| `keepMonthly` | integer | --keep-monthly n   | For the last n months which have one or more snapshots, only keep the last one for that month.     |
| `keepYearly`  | integer | --keep-yearly n    | For the last n years which have one or more snapshots, only keep the last one for that year.       |
| `keepTags`    | array   | --keep-tag <tag>   | Keep all snapshots which have all tags specified by this option (can be specified multiple times). [`--tag foo,tag bar`](https://github.com/restic/restic/blob/master/doc/060_forget.rst) style tagging is not supported. |
| `prune`       | bool    | --prune            | If set, actually removes the data that was referenced by the snapshot from the repository.         |
| `dryRun`      | bool    | --dry-run          | Instructs `restic` to not remove anything but print which snapshots would be removed.              |

You can set one or more of these retention policy options together. To learn more, read [here](
https://restic.readthedocs.io/en/latest/manual.html#removing-snapshots-according-to-a-policy).

## Next Steps

- Learn how to create `BackupConfiguration` crd for different backends from [here](/docs/guides/backends/overview.md).
- Learn how Stash backup workloads data from [here](/docs/guides/workloads/backup.md).
- Learn how Stash backup databases from [here](/docs/guides/databases/backup.md).
