# SFTPGO SAL

## Initialization

- to use sftpgo sal you need to have a jwt, this jwt can be generated by visiting this page [jwt-token](http://localhost:8080/api/v2/token)
- create a new sal as follows

```
args := sftpgo.SFTPGOClientArgs{
	url: "http://<HOSTNAME>/api/v2",
	jwt: "<JWT>"
}
mut cl := sftpgo.new(args)
```

## User management

### Create a new user

```
mut user := sftpgo.User {
    username: "test_user"
    email: "test_email@test.com"
    password: "test_password"
    public_keys: ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDTwULSsUubOq3VPWL6cdrDvexDmjfznGydFPyaNcn7gAL9lRxwFbCDPMj7MbhNSpxxHV2+/iJPQOTVJu4oc1N7bPP3gBCnF51rPrhTpGCt5pBbTzeyNweanhedkKDsCO2mIEh/92Od5Hg512dX4j7Zw6ipRWYSaepapfyoRnNSriW/s3DH/uewezVtL5EuypMdfNngV/u2KZYWoeiwhrY/yEUykQVUwDysW/xUJNP5o+KSTAvNSJatr3FbuCFuCjBSvageOLHePTeUwu6qjqe+Xs4piF1ByO/6cOJ8bt5Vcx0bAtI8/MPApplUU/JWevsPNApvnA/ntffI+u8DCwgP"]
    permissions: {"/": ["*"]}
    status: 1
}
// add user
created_user := cl.add_user(user)  or {
    logger.error("Failed to add user: $err")
    exit(1)
}
```

### Get user

```
returned_user := cl.get_user("test_user") or {
    logger.error("failed to get user: $err")
    exit(1)
}
logger.info("got user: $returned_user")
```

### Update user

```
// update user
user.email = "test_email@modified.com"
cl.update_user(user)  or {
    logger.error("failed to update user: $err")
    exit(1)
}
```

### Delete user

```
cl.delete_user("test_user") or {
    logger.error("failed to update user: $err")
    exit(1)
}
```

## Folder management

### Create folder

```
// create folder struct
mut folder := sftpgo.Folder{
    name: "folder2"
    mapped_path: "/folder2"
    description: "folder 2 description"
}
```

### Get folder

```
returned_folder = cl.get_folder(folder.name) or {
    logger.error("failed to get folder: $err")
    exit(1)
}
logger.info("folder: $returned_folder")
```

### Update folder

```
folder.description = "folder2 description modified"
cl.update_folder(folder)  or {
    logger.error("failed to update folder: $err")
    exit(1)
}
```

### List folders

```
//list all folders
folders := cl.list_folders() or {
    logger.error("failed to list folder: $err")
    exit(1)
}
logger.info("folders: $folders")
```

### Delete folder

```
cl.delete_folder(folder.name) or {
    logger.error("failed to update user: $err")
    exit(1)
}
```

## Roles management

### Add role

```
// create role struct
mut role := sftpgo.Role{
    name: "role1"
    description: "role 1 description"
    users: []
    admins: []
}

//add Role
created_role := cl.add_role(role)  or {
    logger.error("failed to add role: $err")
    exit(1)
}
logger.info("role created: $created_role")

```

### Get role

```
//get role
returned_role := cl.get_role(role.name) or {
    logger.error("failed to get folder: $err")
    exit(1)
}
logger.info("role: $returned_role")
```

### Update role

```
//update role
role.description = "role1 description modified"
cl.update_role(role)  or {
    logger.error("failed to update role: $err")
    exit(1)
}
```

### Delete role

```
//delete role
deleted := cl.delete_role(role.name) or {
    logger.error("failed to update role: $err")
    exit(1)
}
logger.info("role deleted: $deleted")
```

### List roles

```
// list existing roles
roles := cl.list_roles() or {
    logger.error("failed to list roles: $err")
    exit(1)
}
logger.info("roles: $roles")
```

## Events

### Get fs events

```
fs_events := cl.get_fs_events(0, 0, 100, "DESC") or {
    logger.error("failed to list fs events: $err")
    exit(1)
}

logger.info("fs_events: $fs_events")
```

### Ge provider events

```
provider_events := cl.get_provider_events(0, 0, 100, "DESC") or {
    logger.error("failed to list provider events: $err")
    exit(1)
}
logger.info("provider_events: $provider_events")
```