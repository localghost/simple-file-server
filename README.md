# simple-file-server

HTTP server serving files.

## Endpoints

* `/files/[PATH]`: for accessing served files
* `/filelist`: returns flat list of served files (each file in new line)
* `/shutdown`: instructs the server to shutdown itself
* `/health`: for checking if server is up
