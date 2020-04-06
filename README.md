# BMC Library

This library provides functionality to communicate with the Dell CS24-SC BMC "Embedded server manager."

The BMC uses the outdated cipher suite `TLS_RSA_WITH_RC4_128_SHA`, making it difficult to connect to in modern web browsers.  I made this library to make it easy to automate some BMC actions.

## CLI Usage

```txt
Usage of ./bmco:
  -Action string
        Action to perform on server. Options are: info, start, stop, reset (default "info")
  -IP string
        IP of server to connect to
  -Password string
        Password for BMC
  -Port uint
        Port of server to connect to (default 443)
  -Username string
        Username for BMC
```

## Library Usage

Example of connecting to and starting a server

```go
c, _ := New(context.Background(), os.Getenv("IP"), 443, os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
c.Start(context.Background())
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

## Get certificate of your server

If you need to get the certificate of your server, use this command replacing `IP`.

```sh
openssl s_client -cipher "RC4" -connect IP:443 -showcerts
```
