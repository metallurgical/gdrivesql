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
