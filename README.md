# [Radio-T](https://radio-t.com/) Dumper

A.K.A. Http copier, because you can use it with different
sites, and it will successfully detect file type and download
content.

In an endless loop with a break on delay
(to avoid creating a heavy load on the target server), 
the application will wait for a `OK 200` response, after which it will copy the response
body into the generated file.

## Features

* File type detection (without losing bytes)
* Waiting for broadcast availability 
* Delay for avoiding creating a heavy load on the target server
* Automatic generation of human-readable file names
* Unlimited number of files (Limited by disk space)
* Unit tests

## Configuration

You can use `dumper.yml` file or **environment variables** for 
configuration.

| yml              | env              | default    | description                                                  |
|------------------|------------------|------------|--------------------------------------------------------------|
| source_url       | SOURCE_URL       |            | Address from where to download information                   |
| file_prefix      | FILE_PREFIX      |            | Prefix for files                                             |
| file_date_format | FILE_DATE_FORMAT | 02_01_2006 | Date format for files                                        |
| output_directory | OUTPUT_DIRECTORY | .          | Directory where to save files                                |
| timeout          | TIMEOUT          | 10         | Delay between sending requests if response status is not 200 |
| log_level        | LOG_LEVEL        | info       | Log level                                                    |

## Run with [docker compose](https://docs.docker.com/compose/)

Before starting the application, you can change output
directory volume in docker-compose.yml file from `./records`
to your desired directory.
In container all files will be
in the `/tmp/output` directory.

To start the application, run:

```shell
docker-compose up -d
```
