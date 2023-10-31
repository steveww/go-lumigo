# go-lumigo
OTel Lumigo Go example

See how easy it is to add OTel instrumentation and send it to [Lumigo](https://lumigo.io) for visualisation.

## Prereqisites
To run this example you will need [Docker](https://docs.docker.com/engine/install/) with `docker-compose` installed on your computer.

## Fire It Up
Before you can run the example it needs to be built.

```bash
$ docker-compose build
```

This will take a little while depending on your internet connection speed and the speed of your computer. Once it has been built you can run the example. Get your unique Lumigo token from the dashboard under `Settings -> Tracing`, it is shown under the **Manual Tracing** section.

```bash
$ export LUMIGO_TOKEN="your-token-goes-here"
$ docker-compose up
```

To have traces sent, you will need some load. Run the `test.sh` script to generate a light load, use the `-l` flag to have it loop continuously.

```bash
$ ./test.sh -l
```

Head over to the Lumigo dashboard to see the results.

## Tidy Up
Once you've seen how Lumigo provides truely great insight, you can tidy eveything up.

```bash
$ docker-compose down
$ docker image prune -a
```
