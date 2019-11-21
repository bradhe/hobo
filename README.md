:information_desk_person: **Welcome to Hobo!** :wave:

Hobo is a service for autocompleting cities. Most of the heavy lifting is done
in ElasticSearch, but Hobo provides a standard data model and interface for
city data. It also is responsible for providing a common ID scheme for cities.

## Should I use this?

Probably not. It's meant to support some other services I'm working on. The
data set to actually search is not provided, for instance.

## Usage

There's a docker-compose file provided. Just run `docker-compose up` to get the
latest version of our Docker Hub.

### Commands

There are two commands that are supported.

#### `hobo serve` 

The `serve` command starts up the HTTP server for hobo.

#### `hobo import`

The `import` command imports a snapshot of the data set directly in to
ElasticSearch. Think of it like a migration.
