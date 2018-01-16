# drone-metronome

Drone plugin for deploying jobs to [Metronome](https://dcos.github.io/metronome/).

## Docker

Build the Docker image with the following commands:

```
docker build --rm=true -t quintoandar/drone-metronome .
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e PLUGIN_SERVER=http://master.mesos:9000 \
  -e PLUGIN_METRONOMEFILE=metronome.yaml \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  quintoandar/drone-metronome
```

