REST API Microservice Demo
==========================

Let's learn some Golang with a simple REST API microservice.

Implemented using chi for routing, bbolt for data persistance, swagger for sdk
and ui generation, testify for testing.


Architecture
============

Since this is a very simple microservice, there wasn't much architecture needed.
There is just one simple resource managed. Below is a simple high level diagram
describing the architecture.

```base
                                        ________      _________      __________
  ClientApplication [SDK] ------------ |        |    |         |    |          |
                           < HTTP >    |  HTTP  |    | Message |    | Database |
  ClientApplication [SDK] ------------ | Server | -- | Service | -- |          |
                                       |________|    |_________|    |__________|
                   "/messages"             |
                   "/messages/{messageId}" |
                                           |
                                           |
  Web Browser ------------------------------
                   "/swagger"
```

Documentation
=============

The RESTful API is described in an OpenAPI 3.0 document found here: ([openapi.json](https://github.com/brandonto/rest-api-microservice-demo/docs/openapi.json))


Running
=======

```bash
usage: ./rest-api-microservice-demo db_path [port] [db_bucket_name]
```

Testing
=======

```bash
go test ./...
```
