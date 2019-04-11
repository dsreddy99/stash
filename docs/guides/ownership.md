---
title: File Ownership | Stash
description: Handling Restored File Ownership in Stash
menu:
  product_stash_0.8.3:
    identifier: file-ownership-stash
    name: File Ownership
    parent: guides
    weight: 50
product_name: stash
menu_name: product_stash_0.8.3
section_menu_id: guides
---

# Handling Restored File Ownership in Stash

Stash preserve permission bits of the restored files. However, it may change ownership (owner `uid` and `gid`) of restored files in some cases. This tutorial will explain when and how ownership of restored files can be changed. Then, we will explain how we can avoid or resolve this problem.

## Understanding Backup and Restore Behaviour

At first, let's understand how backup and restore behave in different scenario. A table with some possible backup and restore scenario is given below. We have run different container as different user in different scenario using [SecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/).

| Case  | Original File Owner | Backup `sidecar` User | Backup Succeed? | Restore `init-container` User | Restore Succeed? | Restored File Owner | Restored File Editable to Original Owner? |
| :---: | :-----------------: | :-------------------: | :-------------: | :---------------------------: | :--------------: | :-----------------: | :---------------------------------------: |
|   1   |        root         |     stash (1005)      |    &#10003;     |          stash(1005)          |     &#10003;     |        1005         |                 &#10003;                  |
|   2   |        2000         |      stash(1005)      |    &#10003;     |          stash(1005)          |     &#10003;     |        1005         |                 &#10007;                  |
|   3   |        2000         |         root          |    &#10003;     |             root              |     &#10003;     |        2000         |                 &#10003;                  |
|   4   |        2000         |         root          |    &#10003;     |          stash(1005)          |     &#10007;     |          -          |                     -                     |
|   5   |        2000         |         3000          |    &#10003;     |          stash(1005)          |     &#10007;     |          -          |                     -                     |
|   6   |        2000         |         3000          |    &#10003;     |             root              |     &#10003;     |        2000         |                 &#10003;                  |
|   7   |        2000         |         3000          |    &#10003;     |             3000              |     &#10003;     |        3000         |                 &#10007;                  |

If we look at the table carefully, we will notice following behaviours:

1. User of backup `sidecar` does not have any effect on backup. It just need read permission.
2. User of restore container must have permission to read backup repository. It has to run either as same user as backup sidecar or `root` user. Otherwise, restore will fail.
3. If restore container run as `root` user then original ownership of the restored files are preserved.
4. If restore container run as `non-root` user then restored files ownership is changed to restore container's user and restored files become read only to original user unless it was `root` user.

So, we can see when we run restore container as `non-root` user, it rises some serious concerns as restored files become read only to original user. Next section will discuss how we can avoid or fix this problem.

## Avoid or Fix Ownership Issue

As we have seen when file ownership get changed, the restored files can be unusable to a user. We need to avoid or fix this issue.

At first, let's check who could be possible user's of restored files. There could be two scenario for restored files user.

1. Restored file user is same as original user.
2. Restored file user is different than original user.

### Restored files user is same as original user

This is likely to be the most common scenario. Generally, the same application will use the restored files whose data was backed up. In this case, if your cluster supports running container as `root` user, then it is fairly easy to avoid this issue. We just need to run restore container as `root` user. However, things get little more complicated when your cluster does not support running container as `root` user. In that case, we can do followings:

- Run backup container and restore container as same user as the target container.
- Change ownership of restored files using `chown` from restore container after restore is completed.

For first method, we can achieve this configuring SecurityContext under `RuntimeSetting` of `BackupConfiguration` and `RestoreSession` object. A sample `BackupConfiguration` and `RestoreSession` objects configured SecurityContext to run as same user as original user (let original user is 2000) is shown bellow.

**BackupConfiguration:**

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: BackupConfiguration
metadata:
  name: deployment-backup
  namespace: demo
spec:
  repository:
    name: local-repo
  schedule: "* * * * *"
  target:
    ref:
      apiVersion: apps/v1
      kind: Deployment
      name: stash-demo
    volumeMounts:
    - name: source-data
      mountPath: /source/data
    directories:
    - /source/data
  runtimeSettings:
    container:
      securityContext:
        runAsUser: 2000
        runAsGroup: 2000
  retentionPolicy:
    name: 'keep-last-5'
    keepLast: 5
    prune: true
```

**RestoreSession:**

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: RestoreSession
metadata:
  name: deployment-restore
  namespace: demo
spec:
  repository:
    name: local-repo
  rules:
  - paths:
    - /source/data
  target:
    ref:
      apiVersion: apps/v1
      kind: Deployment
      name: stash-demo
    volumeMounts:
    - name:  source-data
      mountPath:  /source/data
  runtimeSettings:
    container:
      securityContext:
        runAsUser: 2000
        runAsGroup: 2000
```

Second method is necessary when backup container was not run as same user as the target container. This is similar to the process where restored file user is different than original user. In this case, we have to change the ownership of restored files using `chown` after restore. This method allows us to change ownership to not only to original user but also to any user. Next section will explain this process.

### Restored file user is different than original user

This is advance use case. If you want to use restored files as different user than original one, then you have to configure stash to change ownership to target user after restore.

You can provide `UID` and `GID` of expected owner in `spec.rules` section of `RestoreSession` object using `owner.uid` and `owner.gid` field of each rules. If you run restore container as root user (which is default behaviour) then you don't need any additional configurations. Otherwise, you have to give `CAP_CHOWN` capability to to restore container using `container.securityContext.capabilities.add` field of `RuntimeSetting` of `RestoreSession` object.

A sample `RestoreSession` object is shown below which run restore container as `non-root` and configured to change restored file ownership after restore.

```yaml
apiVersion: stash.appscode.com/v1beta1
kind: RestoreSession
metadata:
  name: deployment-restore
  namespace: demo
spec:
  repository:
    name: local-repo
  rules:
  - paths:
    - /source/data
    owner: # specify uid and gid of expected owner
      uid: 5000
      gid: 5000
  target:
    ref:
      apiVersion: apps/v1
      kind: Deployment
      name: stash-demo
    volumeMounts:
    - name:  source-data
      mountPath:  /source/data
  runtimeSettings:
    container:
      securityContext: # don't need this part if your cluster allow to run container as root user
        runAsUser: 3000
        runAsGroup: 3000 
        capabilities:
          add: ["CHOWN"]
```
