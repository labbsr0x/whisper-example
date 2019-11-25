# Web Example

This is a small example on how to use whisper-client on a web application.

## Usage

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
