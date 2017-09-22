# simple-file-server

Simple HTTP server for serving files.

## Endpoints

Currently, following endpoints are available:

* `/files/[PATH]`: for accessing served files
* `/filelist`: returns flat list of served files (each file in new line)
* `/shutdown`: instructs the server to shutdown itself
* `/health`: for checking if server is up