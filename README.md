# Why these kind of structure

I am experimenting with go-project structures for a while and think that for the small projects a provided project-layout 
is one of the best. 
All the folders have understandable naming that helps to find the reason of folders existing. I thought a bit 
about using cleaning architecture, but decided not to used it, because it's massive and sometimes overhelmed.

# Observability

`Jaeger` is using for tracing. `Zap` is using for logging. `Prometheus` is using for metrics. I decided not to include scripts
for running `Jaeger/Prometheus` collectors instances, because It's not a hard task and I had limited amount of time.
I added metrics for incoming/successful/failed request and tend to think that this is a minimum basic of metrics. 

# Unit-testing

I had limited amount of time, so decided to write only tests for `controllers.Email` just to show my approach. 
For getting 100% encapsulation logger can be passed as an interface 
(it'd be needed to create interface for logger, but I'm not sure, that it should be done). For `malgun`-implementation
tests can be written without any troubles because of using interfaces too. For others implementations some interfaces 
should be created first and then there won't be any problems with writing tests.

I prefer to use code-generated mocks, because there are no `{}interface` as an argument for them and we have only one
instance of the mock and it's being used in all project. That helps us to maintain code simplicity

For running tests call `make test`

# Focus-area
> Scale. Implement a dummy console email provider that takes at least 100 ms to respond
> and scale your service to be able to send 100 emails per second or more.

I created one more provider for the sender and just put time.Wait with random interval inside of it.
Because service works well with even a single instance, I decided to scale it horizontally. I made it in a very simple way.
I just created docker-compose file with `nginx` like a `load-balancer` and 3 instances of the service. If I misunderstood the task, I can do something different.


# Running service

#### Running 3 nodes + load-balancer:

`make deploy` or
`cd deploy && docker-compose up`

`cd deploy` - Lazy way of running docker-compose :)


#### Running 1 node (Docker):
`make docker` or
```
docker build -t email_sender .
docker run -p 5678:5678 email_sender
```

#### Running with out containers:

`go run ./cmd/ -config "./configuration"`

or

`make go-run`