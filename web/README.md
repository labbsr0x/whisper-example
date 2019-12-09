# Web Example

This is a small example of how to use whisper-client on a web application. It is as simple as a home page and a dashboard page.

## Usage (Locally)

Run the docker stack.
    
```bash
docker-compose up -d
```

### Requirements

* Docker 19.03.4 or later
* Docker Compose 1.24.1 or later

## Development

1. Add Hydra to Hosts

    ```bash
    sudo echo "127.0.0.1 hydra" >> /etc/hosts
    ```

2. Run the necessary applications;

    ```bash
    docker-compose up -d local
    ```

3. Run the commands below;

    ```bash
    go build && ./whisper-example
    ```
   
### Requirements
   
* Docker 19.03.4 or later
* Docker Compose 1.24.1 or later
* Golang version X (if not using docker to build and run)

## Observations

+ Whisper Helper only shows a recommended way to use whisper client APIs;
+ Whisper Helper provides a recommended routes for post-login and post-logout and the respective haldlers;
+ It is necessary to have a post logout URL to remove the cookies. Whisper Logout Link only erase the session cookies on the server and the cookies on the user browser still valid;
