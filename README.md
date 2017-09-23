# simple-file-server

HTTP server serving files.

## Endpoints

* `/files/[PATH]`: for accessing served files
* `/filelist`: returns flat list of served files (each file in new line), supported query parameters:
  * `file`
    * `any`: return files and directories
    * `file`: return only files
    * `dir`: return only directories
  * `recursive`
    * `yes`: get list for the entire file tree
    * `no`: list specific directory
  * `startswith`: base directory to list (if omitted root directory is used)
* `/shutdown`: instructs the server to shutdown itself
* `/health`: for checking if server is up

### `/filelist` examples

#### Return all served files

```
curl "localhost:3000/filelist"
```

#### Return files and directories only from the root dir

```
curl "localhost:3000/filelist?recursive=no&type=any"
```

#### Return files only from specific sub-directory

```
curl "localhost:3000/filelist?startswith=/subdir&recursive=no"
```
