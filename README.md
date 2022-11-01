# ctf-browser-visitor

Bot visitor for XSS challenges in CTF

The example ctf-browser-visitor binary will allow you to queue up visits for sites with the following request:

```
http://localhost:8080/?job={"url":"url_to_visit","cookies":{"key":"value"}}
```


Alternatively, this project includes a wrapper around the selenium API with a worker pool. See [main.go](https://github.com/tnek/ctf-browser-visitor/blob/master/main.go) for an example of how to use this.

