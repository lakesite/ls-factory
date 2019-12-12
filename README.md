# ls-factory #

🏭 A simple task queue. 🏭

This package is not for production use.

## usage ##

The factory package provides you with a dispatcher, which can queue
jobs for processing by workers.  It uses [logrus](https://github.com/sirupsen/logrus)
for logging.  Usage is straight forward:

```
	// Create a dispatcher:
	dispatcher := factory.CreateNewDispatcher(log)

	// Enqueue n*2 jobs:
	n := 10
	for i:= 0; i<n; i++ {
		dispatcher.QueueJob(process1)
		dispatcher.QueueJob(process2)
	}

	// Tell the dispatcher to start processing jobs with n workers:
	if err := dispatcher.Start(n); err != nil {
		log.Fatal("Dispatch startup error: %s", err.Error())
	}

	for dispatcher.Running() {
		if dispatcher.Finished() {
			log.Info("All jobs finished.")
			break
		}
	}
```

A full example is provided in [examples](examples/main.go)

## testing ##

    $ go test

## license ##

MIT