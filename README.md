
# scooter-spotter
[![semver](https://img.shields.io/badge/semver-1.0.0-blue.svg?cacheSeconds=2592000)](https://github.com/kakkoyun/scooter-spotter/releases) [![Maintenance](https://img.shields.io/maintenance/yes/2019.svg)](https://github.com/kakkoyun/scooter-spotter/commits/master) [![Drone](https://cloud.drone.io/api/badges/kakkoyun/scooter-spotter/status.svg)](https://cloud.drone.io/kakkoyun/scooter-spotter) [![Go Doc](https://godoc.org/github.com/kakkoyun/scooter-spotter?status.svg)](http://godoc.org/github.com/kakkoyun/scooter-spotter) [![Go Report Card](https://goreportcard.com/badge/github.com/kakkoyun/scooter-spotter)](https://goreportcard.com/report/github.com/kakkoyun/scooter-spotter) [![](https://images.microbadger.com/badges/image/kakkoyun/scooter-spotter.svg)](https://microbadger.com/images/kakkoyun/scooter-spotter) [![](https://images.microbadger.com/badges/version/kakkoyun/scooter-spotter.svg)](https://microbadger.com/images/kakkoyun/scooter-spotter)

An API service to search for available Scooters.

## How does it work

`scooter-spotter` is a small binary program that searches and finds available scooters in given time interval (default: 1 second).

It is a simple map-reduce task system which schedules `max * 2` and collects results from backing API.
If the given time-out interval is exceeded or if necessary maximum limit is reached, it returns collected results.

## Examples

```console
$ curl 'localhost:4000/scooters?max=10'
[{id: 1, battery_level: 90, available_for_rent: true}]
```

## Usage

### Using executable (with CLI args)

```console
$ ./scooter-spotter -h

Usage of scooter-spotter:
  -grace-period int
    	grace period to wait for connections to drain, in secods (default 5)
  -port int
    	host port to listen (default 4000)

```

### Using Docker (with Environment variables)

```bash
$ docker run --rm \
      -e SCOOTER_SEARCH_API_URL=mybeautifulseachapi.com \
      -e SCOOTER_SEARCH_API_TIMEOUT=10 \
      kakkoyun/scooter-spotter
```

## Variables

- `SCOOTER_SEARCH_API_URL`: Backend Seach API URL
- `SCOOTER_SEARCH_API_TIMEOUT`: Timeout value in seconds to execute search queries

## Development

### Build Binary

Build the binary with the following commands:

```console
$ go build .
```

### Build Docker

Build the docker image with the following commands:

```console
$ make docker-build
```

### Release Docker

```console
$ make docker-push
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/kakkoyun/scooter-spotter/tags).

## License and Copyright

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details
