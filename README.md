# simple-file-server

HTTP server serving files.

## Parameters

* `-h`: print help message
* `-served-dir`: directory to serve files from (default: current working directory)
* `-listen-addr`: address to listen on (default: 0.0.0.0)
* `-listen-port`: port to listen on (default: 3000)

## Endpoints

* `/files/[PATH]`: access served files
* `/filelist`: return flat list of served files (each file in new line), supported query parameters:
  * `type`
    * `any`: return files and directories
    * `file`: return only files (default)
    * `dir`: return only directories
  * `recursive`
    * `yes`: list entire file tree (default)
    * `no`: list specific directory only
  * `startswith`: base directory to list (if omitted root directory is used)
* `/shutdown`: instruct the server to shutdown itself
* `/health`: check if server is up

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
