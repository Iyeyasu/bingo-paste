# Bingo

Bingo is a light-weight pastebin written with Go. The goal of the project is to keep both the server-side and client-side as small and simple as possible in response to the large number of bloated pastebins out there.

TODO Screenshots

## Features

* Pasting
* Binning
* Light/dark theme
* User authenticaton (optional)
* Visibility (optional)
* Expiring pastes (optional)
* Syntax highlighting (optional)
* Lightweight and fast
* Single executable
* Javascript free

## Getting started

### Configuration

An example configuration can be found in [config/bingo.example.yml](./config/bingo.example.yml). An example docker compose file can be found in `docker-compose.example.yml`. Make sure you mount your custom configuration file in your `docker-compose.yml` file.

### Building and running

#### Docker (preferred)

To build the server in Docker:

```bash
make docker
```

After building, you should have a Docker image called `bingo`. Run your `docker-compse.yml` using:

```bash
docker-compose up -d
```

The service should now be up and running in the port you mapped in your compose file (in the example file it's `http://localhost:8080`).

#### Plain Go

You must have PostgreSQL and optionally Redis setup on your host machine and change the related setting in the `bingo.example.yml` file. 

To build the server on your host:

```bash
make release
```

The binary can be found in `build/server`. Make sure to pass the configuration file to the program:

```bash
build/server <path-to-config-file>
```

The service should now be up and running in the port defined in your configuration file.


### Using

The default username and password for an admin user are `admin` and `admin`. Note that you can disable user management and make the service anonymous from the configuration file with the setting `auth: enabled: false`.

## TODO

- LDAP
- Minimize templates

LOW PRIO:
- SMTP
    - Send register mail/password recovery?
- Backend
    - MySQL/SQLite options?
