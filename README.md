# BMC Library

This library provides functionality to communicate with the [Dell CS24-SC](https://aramblinggeek.com/on-the-dell-cs24-sc-server/) BMC "Embedded server manager." (I did not write that article). [Ebay link](https://www.ebay.com/p/141699065)

The BMC allows you to perform hardware functions on the server remotely such as turning the server on and off.  The BMC uses the outdated cipher suite `TLS_RSA_WITH_RC4_128_SHA`, making it difficult to connect to in modern web browsers.  I made this library to make it easy to automate some BMC actions.

## CLI Usage

Each flag can also be set as an environment variable (all caps)

```txt
Usage of ./CS24-SC-BMC:
  -Action string
        Action to perform on server. Options are: info, start, stop, reset, monitor (default "info")
  -IP string
        IP of server to connect to
  -Password string
        Password for BMC
  -Port uint
        Port of server to connect to (default 443)
  -Username string
        Username for BMC
```

Docker image available at `vertoforce/bmc-cs24-sc`

```sh
docker run vertoforce/bmc-cs24-sc -IP=10.0.0.10 -Username=root -Password=password -Action=info
```

## Library Usage

Example of connecting to and starting a server

```go
c, _ := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
c.Start(context.Background())
```

## Server monitoring

This tool can also be used to monitor a dell cs24-sc server.  If you set the `action` to `monitor`, every 30 seconds it will log out information about the server.  You can capture these logs and send them to an ELK server using the provided docker-compose.yml.

To bring up the docker stack, make sure you set all the env vars defined in the docker-compose.yml file then run the following command.

```sh
docker stack deploy -c docker-compose.yml server-monitor
```

## BMC API Endpoints

### Login

Logs in to the BMC.  The session cookie is `PHPSESSID`

- URL: `/cgi_bin/login.cgi`
- Method: `POST`
- Body Parameters (`application/x-www-form-urlencoded`)
  - `quser` - string
  - `qpass` - string

### Set Power Control

Sets the power state of the server (start, stop, reset).

- URL: `/cgi_bin/ipmi_set_powercontrol.cgi`
- Method: `POST`
- Body Parameters (`application/x-www-form-urlencoded`)
  - `power_option` - One of `poweron`, `poweroff`, `powerreboot`

### Get temperatures

- URL: `/cgi_bin/ipmi_get_info.cgi?operation=temperature`
- Method: `GET`

Then parse out the temperatures from the table.  See [bmc/temp.go](bmc/temp.go) for implementation.

## Get certificate of your server

If you need to get the certificate of your server, use this command replacing `IP`.

```sh
openssl s_client -cipher "RC4" -connect IP:443 -showcerts
```
