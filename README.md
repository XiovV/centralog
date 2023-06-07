# Overview
`centralog` is a CLI tool used for reading logs from docker containers. It can read logs in real time from multiple containers, 
as well as filter them by certain criteria.

It's ideal for use cases where you want to read logs from multiple containers at once. Centralog will aggregate the logs 
and display them, making it seem like you are reading logs from only one container.

# Setup
Centralog consist of two parts: the `agent`, which runs on your server, and the CLI tool, which runs on your personal machine.

## Agent setup
Running the `agent` is as simple as running a single docker command:
```shell
sudo docker run -v ~/centralog/data:/data -p 8080:8080 \
-v /var/run/docker.sock:/var/run/docker.sock \ 
-e CONTAINERS=logserver1,logserver2 \
--name centralog xiovv/centralog:latest
```

The `CONTAINERS` environment variable tells the `agent` which containers it should monitor. If you want it to monitor all of your containers, then set the `MONITOR_ALL_CONTAINERS` to true.

Output:
```
Your new API key is: PD7Xk8WzmSpWl6MgeWl0lkMI6YBBqY8KWUN2457kr
2023-06-07T14:26:16.762+0200    INFO    agent/main.go:52        storing containers      {"containers": ["logserver1", "logserver2"]}
2023-06-07T14:26:16.764+0200    INFO    agent/main.go:72        initialising log listener...
2023-06-07T14:26:16.764+0200    INFO    agent/server.go:72      getting container       {"containerName": "logserver1"}
2023-06-07T14:26:16.768+0200    INFO    agent/server.go:72      getting container       {"containerName": "logserver2"}
2023-06-07T14:26:16.770+0200    INFO    agent/main.go:79        server is listening for requests...     {"port": "8080"}
```

Take note of that API key written at the very top of the output, we will need it for setting up the CLI. And that's it, the agent is ready.

## CLI Setup
### Build and install from source
Clone the repository:
```shell
git clone https://github.com/XiovV/centralog.git
```
Build the CLI binary:
```shell
cd centralog/cmd/cli
go build
```
Now you should see a new binary called `cli`.

### Adding a new node:
To add a new node, run the `./cli add node` command, and an input prompt will show up:
```
? Enter your node's URL: localhost:8080
? Enter your node's API key: PD7Xk8WzmSpWl6MgeWl0lkMI6YBBqY8KWUN2457kr
? Enter your node's custom name: myNewNode
Node myNewNode added successfully
```

If any of the inputs is invalid (such as if the provided URL is unreachable, or if the API key is incorrect) you will receive an error message, so you cannot add a node with incorrect data by accident.

### Monitoring logs
In order to monitor logs in real time, run this command:
```shell
$ ./cli logs myNewNode
```
```
container: logserver1 | timestamp: 07/06/2023, 14:42:59 | message: 2023-06-07T12:42:59+0000 DEBUG This is a debug log that shows a log that can be ignored.
container: logserver2 | timestamp: 07/06/2023, 14:42:59 | message: 2023-06-07T12:42:59+0000 WARN A warning that should be ignored is usually at this level and should be actionable.
container: logserver2 | timestamp: 07/06/2023, 14:43:02 | message: 2023-06-07T12:43:02+0000 INFO This is less important than debug log and is often used to provide context in the current task.
container: logserver1 | timestamp: 07/06/2023, 14:43:03 | message: 2023-06-07T12:43:03+0000 ERROR An error is usually an exception that has been caught and not handled.
container: logserver2 | timestamp: 07/06/2023, 14:43:03 | message: 2023-06-07T12:43:03+0000 WARN A warning that should be ignored is usually at this level and should be actionable.
container: logserver1 | timestamp: 07/06/2023, 14:43:04 | message: 2023-06-07T12:43:04+0000 DEBUG This is a debug log that shows a log that can be ignored.
container: logserver1 | timestamp: 07/06/2023, 14:43:05 | message: 2023-06-07T12:43:05+0000 INFO This is less important than debug log and is often used to provide context in the current task.
````
You should see the logs being displayed as they are coming in real time. As you can see, both `logserver1` and `logserver2`'s logs are being displayed.

### Listing nodes
If you'd like to list your nodes and check their statuses, run this command:
```shell
$ ./cli list nodes
```
```
NAME          CONTAINERS     STATUS     
myNewNode     2/2            UP
```

### Listing containers
If you'd like to list the containers of a specific node, run this command:
```shell
$ ./cli list containers myNewNode
```
```
NAME           STATUS
logserver1     running
logserver2     running
```

