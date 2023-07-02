# transac_api

Transac_api is an api interface to the jump_database users, invoice and transaction management.

**Version used:** go 1.20

The project exposes four routes on the port 8000:
- GET `/health`: Writes 200 response - used to check if server is alive
- GET `/users`: Get a list of all users.
- POST `/invoice`: Add invoice for user.
- POST `/transaction`: Handle paid transaction for an invoice.

You can read the [openapi file](./openapi.yaml) for more information.

## Start project
You can source `build.sh` file to use building/running helpers.
```shell
source build.sh
```

### Building
To build project in docker, use `build` or
```shell
docker build -t "transac_api":latest -f ./Dockerfile .
```

### Running
To run project in a docker container, using docker-compose, use `run` or
```shell
docker-compose up 
```
The project runs on PORT 8000.

### Testing?
As this backend test was all database logic and API web server, 
I did not implement unit tests, or an integration test suite. 

## CRA

| Durée     | Note                                        |
|-----------|---------------------------------------------|
| 40min     | Architecture du projet, server de base fini |
| 1h10min   | Users + Invoice finis                       |
| 1h5min    | Toutes routes fonctionnelles                |
| 1h        | Dockerfile + build.sh + polishing           |
| TOTAL ~4h |                                             |

