# Introduction
`gdrivesql` is a module running on top of golang language that enable developers/devOps to automate daily process to backup both filesystems and databases(.sql) into `Google Drive` without having to do it manually. All these are done by executing golang binary file(executable file).

Since this package using `Google Drive` as a place to store backups files, its required few steps for authorization for the first time. For the subsequents execution, authorization no longer needed.

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
Paste authorization key back to terminal. This time, file `token.json` automcatically created for you and placed under root directory

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
# List of available database's name.
name: # Database's name that will be exported into .sql format
  - DatabaseA
  - DatabaseB
```

From the above config, this module will export `DatabaseA.sql` and `DatabaseB.sql` into `temp` folder.

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