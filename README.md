# Introduction
`gdrivesql` is a module running on top of golang language that enable developers/devOps to automate daily process to backup both filesystems and databases(.sql) into `Google Drive` without having to do it manually. All these are done by executing golang binary file(executable file).

Since this package using `Google Drive` as a place to store backups files, its required few steps for authorization for the first time. For the subsequent execution, authorization no longer needed.

Until now, only `MYSQL` and `postgreSQL` are supported. For the backup, it doesn't has to backup both databases and filesystems at the same time. You may choose either one or both.

## Demo
View demo below to get the whole picture and overview how things works.

### Get Client Credentials
Getting client credentials and placed directly(credentials.json) into root folder of `gdrivesql`. Head over to [Google Drive Api Quickstart](https://developers.google.com/drive/api/v3/quickstart/go) to get client credential file.

![1](https://i.imgur.com/pT3SOzY.gif)

### Get Consent Link
Get consent link by using `gdriveauth` executable to paste into browser URL

![2](https://i.imgur.com/iRqjcsq.gif)

### Get Authorization Key
Get authorization key

![3](https://i.imgur.com/fcmbAkb.gif)

### Paste Authorization Key
Paste authorization key back to terminal. This time, file `token.json` automatically created for you and placed under root directory

![4](https://i.imgur.com/BuKDvAb.gif)

### Get Google Drive ID
This module require Google Drive ID to upload backup files.

![5](https://i.imgur.com/gWL6fgB.gif)

### Backup and Upload Files
Use `gdrivesql` executable file to create a backup file and upload into google drive's folder

![6](https://i.imgur.com/9LNTh3V.gif)

### Done
Check google drive's folder for uploaded files

![7](https://i.imgur.com/hL6Eetr.gif)


## Configurations
This module provide 3 config files:

- **Database config** : List of database's name to export into .sql
- **Filesystem config** : List of absolute filesystem's(project) path that need to archive to backup into google drive
- **Google drive config** : Defined which filesystem and databases files(from above config) that able to upload into google drive

### Database Config
```yaml
---
# List of connection's name
connections:
  - name: "connection_a"
    driver: "mysql"
    host: "127.0.0.1"
    port: "3306"
    user: "root"
    password:

  - name: "connection_b"
    driver: "mysql"
    host: "external_ip_address"
    port: "3306"
    user: "root"
    password: "root@1234"

  - name: "connection_c"
    driver: "postgres"
    host: "127.0.0.1"
    port: "5432"
    user: "postgres"
    password:
       

# List of available database's name that need to export.
databases:
  - connection: connection_a # Will use `connection_a`
    list:
      - DatabaseA
      - DatabaseB

  - connection: connection_c # will use `connection_c`
    list:
      - DatabaseC
```

From the above config, this module will export `DatabaseA.sql`, `DatabaseB.sql` and `DatabaseC.sql` into `temp` folder. 

** Notes for **postgreSQL**, you may need to create `.pgpass` file under home directory. On windows the file is named `%APPDATA%\postgresql\pgpass.conf` while on linux/unix should be `~/.pgpass`. 

This file should contain lines of the following format:

```
hostname:port:database:username:password
```

On Unix systems, the permissions on a password file must disallow any access to world or group; achieve this by a command such as `chmod 0600 ~/.pgpass`. If the permissions are less strict than this, the file will be ignored. On Microsoft Windows, it is assumed that the file is stored in a directory that is secure, so no special permissions check is made.


### Filesystem Config
```yaml
---
# List of filesystem path's name.
path: # Absolute path to filesystem that need to be archive
  - "/Users/metallurgical/projects/automate" # "/path/to/the/folder/to/archive"
```

From the above config, this module will compress and create `automate.tar.gz` into `temp` folder.

### Google Drive Config
```yaml
---
# List of available database's name.
config:
  - folder: "automate" # Folder name to create inside temp folder to store backup files
    filesystem: true # Set to true to backup and upload along with
    driveid: "Google Drive ID" # Backup archive will be stored under this google drive's folder
    files: # Archived and database name that will stored under "automate" folder. E.g: automate.tar.gz and dbname.sql
      - "automate" # Filesystem: this name must be matched with folder defined inside filesystem.yaml(if exist)
      - "DatabaseA" # Database name defined inside database.yaml
      - "DatabaseB" # Database name defined inside database.yaml
```

From the above config, this module will move file `automate.tar.gz`, `DatabaseA.sql`, `DatabaseB.sql` into `automate`(depend on `folder` option) folder and finally compress those folder into `automate.tar.gz`. 

`gdrivesql` module will upload `automate.tar.gz` into google drive's folder(depend on `driveid` option) with the name `backup.tar.gz`. 

## Installation & Usage
Make sure you install golang in your server. Clone this repository somewhere. Head over into cloned repository and run below command to build:

```
$ cd gdrivesql
```

Build `gdriveauth` to get `token.json`. **This command must be execute first before able to upload into google drive(required to run one time only)**

```
$ go build cmd/gdriveauth/gdriveauth.go
$ ./gdriveauth
```

Build `gdrivesql` to do start backup filesystem and database

```
$ go build cmd/gdrivesql/gdrivesql.go
$ ./gdrivesql
```

Build `cleanup` to remove all leftovers files inside `temp` folder

```
$ go build cmd/cleanup/cleanup.go
$ ./cleanup
```