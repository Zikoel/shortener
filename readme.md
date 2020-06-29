# URL Shortener

A very simple url shortener redis based

## Quick start

1. run `make docker`
2. run `docker-compose up`

## How to compile

1. run `make build`

## How to test

1. run `make test`

## API documentation

A very synthetic APIs usage can be found at the address `/api/usge` provided directly by the app

## Application parameters

**SERVER_PORT:** Where the server will listen for incoming requests (default: 5000)  
**REDIS_HOST:** The host were the redis server is listing (default: localhost)  
**REDIS_PORT:** The port were the redis server is listing (default: 6379)  
**REDIS_PASSWORD:** The password for the redis connection (default: "")  
**IN_MEMORY_PERSISTENCE:** The flag to switch from redis memory to in-memory persistence (actually non used) (default: false)  

## Algorithm consideration

We have few cases to care about:
* The user suggests a key
    * No collision: The key doesn't exist yet, so we can simply add it
    * There is collision with another key and, the URL content is the same: Nothing to do here, we can safely return the suggested key
    * There is collision with another key and, the URL content is not the same: Here we need to reject the newer suggested key because is already taken
* The user doesn't suggest a key, so we calculate an hash code from the URL text:
    * No collision: The key doesn't exist yet, so we can simply add it
    * There is collision with another key, the URL content is the same: Nothing to do here, we can safely return the calculated key
    * There is collision with another key and, the URL content is not the same: Here we need to generate a non-existing yet random key and persist the URL with that

## Tests consideration

The tests cover only the part where there is more logic, the other part should be tested but out of scope of the experiment. 

## References

[Project structure](https://github.com/golang-standards/project-layout)  
[Makefile inspiration](https://sohlich.github.io/post/go_makefile)
