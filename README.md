# echo-service

## Usage

From the repository root, run:

```sh
docker build -t echo-service:0.1.0 .
docker run --rm -p 8080:8080 echo-service:0.1.0
```

Then in another terminal instance, run:

```sh
curl -s http://localhost:8080/hello | jq .
```

and you should get

```json
{
  "method": "GET",
  "path": "/hello",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "User-Agent": [
      "curl/8.12.1"
    ]
  },
  "host": "localhost:8080",
  "remote_addr": "[::1]:65325"
}
```
