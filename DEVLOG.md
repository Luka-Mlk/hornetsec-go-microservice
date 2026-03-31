# Devlog

This devlog's intent is to track my ideas, thoughts, vizualizations and pivots I make as I solve this home assignment

## Requirement

The requirements mention a RESTful API in go, refering to the payload as 'documents' containing an ID, name and description.

Initial thoughts to this project look like an S3 clone due to refering to the payload as a document. 

That would mean I need to handle byte streams or similar ways of sending large files.

But second read the service feels like its a middleman service to fetch document metadata. 
Similar to how S3 stores metadata in a relational db and serves it from the most optimal source

## Implementation brainstorming

Since this is a microservice we can have an RPC interface and HTTP interface in case it's a user facing or service facing microservice. 
It isn't specified in the requirements

I would start with an easy to test entrypoint so I can test and log as I go.

### Standard library or framework for handling network?

For http endpoints I will use the default http/net standard library

I'll also do a RPC implementation using gRPC with google's standard library 

Using interfaces for a request is a must to avoid writing double the endpoints

### Mini error handling framework?

To keep error handling and responses clean a small error handling framework can be built

Or at least a light mapping of responses and errors

### Keeping files in memory 

#### What if service is shut down? 

To avoid data loss because the data is kept in memory i want to save the information stored.

This requires some kind of persistence.

When though?

Periodically checking if there is changes to the data and if there is saving to disc?

However requirements state only in memory. 

Will leave interface open to improvement but closed to change

#### Avoiding race conditions?

Simple mutex on storage struct?

Possible queue implementation, however concept is too simple for now

#### Gracefull shutdown?

Make sure all requests finish before closing service/look into?

### Scaling?

Initially considered redis for in memory store.

But in order not to add any external dependencies this step is skipped for an in memory key value map as a store

### Call tracking/Call ID?

In order to keep track of calls a simple call id can be generated upon entry to the service, maybe UUID?

So i dont end up with tons of different id's ill have the service check if there is already an existing id and if not create a new one.

Utilize context
