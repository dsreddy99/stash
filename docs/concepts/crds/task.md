---
title: Task Overview
menu:
  product_stash_0.8.3:
    identifier: task-overview
    name: Task
    parent: crds
    weight: 35
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: concepts
---

> New to Stash? Please start [here](/docs/concepts/README.md).

# Task

## What is Task

Individual [Function](/docs/concepts/crds/function.md) perform a step of a backup or restore process. An entire backup or restore process needs ordered execution of one or more functions. A `Task` is a Kubernetes `CustomResourceDefinition (CRD)` which specifies such order of the functions along with their inputs for a backup or restore process in Kubernetes native way.

When you install Stash, some default `Tasks` will be automatically created for supported targets. However, you can create your own `Task` to customize or extend backup/restore process.

You can also add one or more steps in the `Task` to execute them as pre-backup or post-backup hook. Stash will execute these hooks in the order you have specified.

## Task CRD Specification

Like other official Kubernetes resources, `Task` has `TypeMeta`, `ObjectMeta` and `Spec` sections. However, unlike other Kubernetes resources, it does not have a `Status` section.

A sample `Task` object to backup a PostgreSQL database is shown below,

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: Task
metadata:
  name: pg-backup
spec:
  steps:
  - name: pg-backup
    params:
    - name: outputDir # specifies where to write output file
      value: /tmp/output
    - name: secretVolume # specifies where backend secret has been mounted
      value: secret-volume
  - name: update-status
    params:
    - name: outputDir # specifies where previous step wrote output file. it will read that file and update status of respective resources accordingly.
      value: /tmp/output
  volumes:
  - name: secret-volume
    secret:
      secretName: ${REPOSITORY_SECRET_NAME}
```

This `Task` uses two functions to backup a PostgreSQL database. The first step indicates `pg-backup` function that dump PostgreSQL database and upload the dumped file. The second step indicates `update-status` function which update status of `BackupSession` and `Repository` crd for respective backup.

Here, we are going to describe some important sections of a `Task` crd.

### Task `Spec` Section

Task object holds following fields in `.spec` section.

#### spec.steps

`spec.steps` section specifies list of functions and their parameters in the order they should be executed. You can also template this section using the [variables](/docs/concepts/crds/functions.md#stash-provided-variables) that Stash can resolve itself. Stash will resolve all the variables and create a pod definition with the container specification specified in the respective `Functions` mentioned in `steps` section.

Each `step` consist of following fields:

- **name :** `name` specifies the name of the `Function` that will executed in this step.
- **params :** `params` specifies an optional list of variables names and their value that Stash should use to resolve respective `Function`. If you use a variable in `Function` specification whose value Stash can not provide, you can pass the value of that variable using this `params` section. You have to specify following fields for a variable:
  - **name :** `name` of the variable.
  - **value :** value of the variable.

In the sample `Task` task specification that has been shown above, we have used `outputDir` variable in `pg-backup` function but Stash can not provide it's value. So, we have passed the value using `params` section of the `Task` object.

>Stash executes the `Functions` in the order they appear in `spec.steps` section. All the functions excepts the last one will be used to create `init-container` specification and the last function will be used create `container` specification for respective backup job. This guarantee an ordered execution of the functions.

#### spec.volumes

`spec.volumes` specifies a list of volumes that should be mounted in the respective job created for this `Task`. In the sample we have shown above, we need to mount storage secret to the backup job. So, we have added the secret volume in `spec.volumes` section. Note that, we have used `REPOSITORY_SECRET_NAME` variable as secret name. This variable will be resolved by Stash from `Repository` specification.

## Next Steps

- Learn how Stash backup databases using `Function-Task` model from [here](/docs/guides/databases/backup.md).
- Learn how Stash backup stand alone PVC using `Function-Task` model from [here](/docs/guides/volumes/backup.md).
