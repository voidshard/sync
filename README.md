### Sync

Simple WaitGroup backed by a Pub/Sub


#### Why

Wanted to implement something like task groups for [asynq](https://github.com/hibiken/asynq).


#### Usage

Using the WaitGroup to fan-in
```golang
group := "test"

ps, _ := sync.NewPubSub(sync.NewOptions())
sink := sync.New(ps)

wg, _ := sink.WaitGroup(context.Background(), group, 512)
wg.Add(2)

wg.Wait(time.Second * 60)
```

Any number of children in the fan-out
```golang
group := "test"

ps, _ := sync.NewPubSub(sync.NewOptions())
sink := sync.New(ps)
	
sink.Done(context.TODO(), group, "one")

sink.Done(context.TODO(), group, "two")
```


#### Implementations

Currently implemented is Redis, but the interface is simple enough to support any pub/sub provider. The New() function accepts the aforementioned interface, so you can pass in your own.

