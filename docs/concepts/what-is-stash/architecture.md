---
title: Stash Architecture
description: Stash Architecture
menu:
  product_stash_0.8.3:
    identifier: architecture-concepts
    name: Architecture
    parent: what-is-stash
    weight: 20
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: concepts
---

# Stash Architecture

Stash is a Kubernetes operator for [restic](https://restic.net/). In the heart of Stash, it has Kubernetes [controller](https://book.kubebuilder.io/basics/what_is_a_controller.html). It uses [Custom Resources Definition(CRD)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)  to specify targets and behaviours of backup and restore process in Kubernetes native way. A simplified architecture of Stash is shown below.

<figure align="center">
  <img alt="Stash Architecture" src="/docs/images/concepts/stash_architecture.svg">
  <figcaption align="center">Fig: Stash Architecture</figcaption>
</figure>

## Components

Stash consists of various components that implements backup and restore logic. This section will give a brief overview of such components.

### Stash Operator

When an user install Stash, it creates a Kubernetes [Deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) typically named `stash-operator`. This deployment controls entire backup and restore process. `stash-operator` deployment run two containers. One is `operator` which run controllers and other necessary stuff and another is `pushgateway` which is a Prometheus [pushgateway](https://github.com/prometheus/pushgateway).

#### Operator

`operator` container runs all the controllers as well as an [Aggregated API Server](https://kubernetes.io/docs/tasks/access-kubernetes-api/setup-extension-api-server/).

##### Controllers

Controllers watches various Kubernetes resources as well as the custom resources introduced by Stash. It applies the backup or restore logics when respective backup or restore is configured for a target resource.

##### Aggregated API Server

Aggregated API Server self-host validating and mutating [webhooks](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/) and run an Extension API Server for Snapshot.

**Mutating Webhook:** Mutating Webhook is used to inject backup `sidecar` or restore `init-container` into a workload if it is configured for backup or restore. It is also used for defaulting custom resoruces.

**Validating Webhook:** Validating Webhook is used to validate custom resources has been defined properly.

#### Pushgateway

`pushgateway` container runs Prometheus [pushgateway](https://github.com/prometheus/pushgateway). All the backup sidecars/jobs and restore init-containers/jobs send Prometheus metrics to this pushgateway after completing backup or restore process. Prometheus server can scrap those metrics from this pushgateway.

### Backend

Backend is the storage where Stash stores backup files. It can be a cloud storage like GCS bucket, AWS S3, Azure Blob Storage etc. or a Kubernetes persistent volume like [HostPath](https://kubernetes.io/docs/concepts/storage/volumes/#hostpath), [PersistentVolumeClaim](https://kubernetes.io/docs/concepts/storage/volumes/#persistentvolumeclaim), [NFS](https://kubernetes.io/docs/concepts/storage/volumes/#nfs) etc. To know more about backend, please visit [here](/docs/guides/backends/overview.md).

### CronJob

When an user creates a [BackupConfiguration](#backupconfiguration) object, Stash creates a CronJob with the schedule specified there. At each scheduled time, this CronJob trigger a backup for the targeted workload.

### Backup Sidecar / Backup Job

When an user creates a [BackupConfiguration](#backupconfiguration) object, Stash inject an `sidecar` to the target if it is a workload (i.e. `Deployment`, `DaemonSet`, `StatefulSet` etc.). This `sidecar` takes backup when the respective CronJob triggers a backup. If the target is a database or stand alone volume, Stash creates a job to take backup at each trigger.

### Restore Init-Container / Restore Job

When an user creates a [RestoreSession](#restoresession) object, Stash inject an `init-container` to the target if it is a workload (i.e. `Deployment`, `DaemonSet`, `StatefulSet` etc.). This `init-container` perform restore process when the respective workload pods restart. If the target is a database or stand alone volume, Stash creates a job to restore.

### Custom Resources

Stash uses [Custom Resources Definition(CRD)](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)  to specify targets and behaviours of backup and restore process in Kubernetes native way. This section will give a brief overview of the custom resources used by Stash.

#### Repository

`Repository` specifies the backend information where backed up data will be stored. User have to create `Repository` object for each backup target. Only one target can be backed up into one `Repository`. For details about `Repository`, please visit [here](/docs/concepts/crds/repository.md).

#### BackupConfiguration

`BackupConfiguration` specifies the backup target, behaviours (schedule, retention policy etc.), `Repository` object that holds backend information etc. User have to create one `BackupConfiguration` object for each backup target. When a `BackupConfiguration` is created, Stash creates a CronJob for it and inject backup sidecar to the target if the target is a workload (i.e. Deployment, DaemonSet, StatefulSet etc.). For more details about `BackupConfiguration`, please visit [here](/docs/concepts/crds/backupconfiguration.md).

#### BackupSession

`BackupSession` object is created by respective CronJob at each backup schedule. It points to respective `BackupConfiguration`. Controller that run inside backup sidecar (in case of backup through job it is stash operator itself) will watch this `BackupSession` object and start taking backup instantly. User also can create a `BackupSession` object manually to trigger backup instantly. For more details about `BackupSession`, please visit [here](/docs/concepts/crds/backupsession.md).

#### RestoreSession

`RestoreSession` specifies where to restore and the the `Repository` that stores backed up data. User have to create a `RestoreSession` object when she want to restore. When a `RestoreSession` is created, Stash inject an `init-container` into the target workload (lunch a job if the target is not an workload) to restore. For more details about `RestoreSession`, please visit [here](/docs/concepts/crds/restoresession.md).

#### Function

A `Function` is a template for a container that performs only a specific action. For example, `pg-backup` `Function` take backup of a PostgreSQL database where `pg-restore` `Function` only restore PostgreSQL database from backed up data. `Function` and `Task` enables user to extend or customize backup/restore process. For more details about `Function`, please visit [here](/docs/concepts/crds/function.md).

#### Task

A complete backup or restore process may consist of several steps. For example, in order to backup a PostgreSQL database we first need to dump the database then upload the dumped file to backend then we need to update `Repository` and `BackupSession` status and send Prometheus metrics. A `Task` is an ordered collection of multiple `Function` that performs such individual step. For more details about `Task`, please visit [here](/docs/concepts/crds/task.md).

#### BackupConfigurationTemplate

`BackupConfigurationTemplate` enables users to provide a template for `Repository` and `BackupConfiguration` object. Then, she just need to add some annotations to the workload she want to backup. Stash will automatically create respective `Repository` and `BackupConfiguration` according to the template. In this way, users can create a single template for all similar kind of workloads and backup them just adding some annotations. In Stash parlance, we call it **default backup**. For more details about `BackupConfigurationTemplate`, please visit [here](/docs/concepts/crds/backupconfiguration_template.md).

#### AppBinding

`AppBinding` holds necessary information to connect with a database. For more details about `AppBinding`, please visit [here](/docs/concepts/crds/appbinding.md).
