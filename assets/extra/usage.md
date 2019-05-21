# How to Use GoSock binaries

This is a simple guide to use the distributed binaries e.g (.exe for windows)
The distributed binaries can be found at the [releases](https://github.com/deven96/gosock/releases)

## Usage

- Click the binary file to run it on a default "localhost:8008" address
- To change the address and port, open cmd-prompt/terminal
- Navigate to the location of the file

``` bash
    # change port to 9000 for windows
    gosock_windows_amd64.exe -addr=":9000"

    # linux
    ./gosock_linux_amd64 -addr=":9000"

    # mac
    ./gosock_darwin -addr=":9000"
```

To host the chat application while accessing it across the network :

- Connect the hosting device (windows/mac/linux) to a wifi/router network
- Check the device IP

```bash
    # linux/mac
    ifconfig
    Look for Network adapter
    Under device, look for "inet" section and next to it will be your computer's IP address

    # windows
    ipconfig
    Look for "Default Gateway" under your network adapter for your router's IP address. 
    Look for "IPv4 Address" under the same adapter section to find your computer's IP address.
```

- Start GoSock server on "IP:desiredport"

```bash
    # using IP 192.168.10.23 and port 9000
    gosock_windows_amd64.exe -addr="192.168.10.23:9000"

    # NOTE: for *nix users, allow traffic on that port
    sudo ufw allow 9000
```

- Connect multiple other devices to the wifi/router network
- Access the chatroom by navigating to the "IP:desiredport" on browser of your choice
