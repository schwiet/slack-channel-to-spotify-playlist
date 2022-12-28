# slack-channel-to-spotify-playlist
A collection of code to automate adding tracks from a Slack channel to a Spotify playlist

## Notes

### Docker Images for Lambdas

#### Simple

Use docker compose to build the images and run the containers:

```
docker compose up
```

#### Manual

To containerize one of the lambdas, build the Docker image by pointing to the `Dockerfile` in the corresponding directory, e.g.

```
docker build --platform linux/amd64 -t sp-auth:0.1 -f ./lambda/spotify/authorize/Dockerfile .  
```

To run locally:

```
docker run --platform linux/amd64 --rm -p 9000:8080 sp-auth:0.1  
```

> NOTE: each of these lambdas expect environment variable(s) to access secrets. When running locally, you can inject these by adding them to the environment before running `docker run` or `docker compose up`