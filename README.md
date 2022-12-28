# slack-channel-to-spotify-playlist
A collection of code to automate adding tracks from a Slack channel to a Spotify playlist

## Notes

### Docker Images for Lambdas

To containerize one of the lambdas, build the Docker image by pointing to the `Dockerfile` in the corresponding directory, e.g.

```
docker build --platform linux/amd64 -t sp-auth:0.1 -f ./lambda/spotify/authorize/Dockerfile .  
```

To run locally:

```
docker run --platform linux/amd64 --rm -p 9000:8080 sp-auth:0.1  
```
