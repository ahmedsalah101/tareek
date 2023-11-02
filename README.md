# tareek

A Reverse Proxy and a load balancer that utilizes go reverse proxies to distribute the load over predefined servers. The used algorithm is the static balancing algorithm (round Robin).

Tareek also supports multi host reverse proxy with a config file to map routes to reverse proxies and Load balancers. Tareek is written in Go.

## Build and Run

```
make run
```

## Example Config

redirect `127.0.0.1/api/` to `http://www.bing.com`

```json
{
    "reverse-proxies": [
        {
            "host": "127.0.0.1",
            "endpoint": "/api",
            "target": "http://www.google.com"
        },
}
```

balance the load over `127.0.0.1/rev1` using the predefined targets:

```json
{
    ...
    "load-balancers": [
        {
            "host": "127.0.0.1",
            "endpoint": "/rev1",
            "targets": [
                "http://www.github.com",
                "http://www.google.com",
                "http://www.bing.com"
            ]
        },
        ...
    ]
}
```

# Acknowledge

the first steps in this project is based on Akhil Sharma's [video](https://www.youtube.com/watch?v=ZSDYx9eOiqo)
