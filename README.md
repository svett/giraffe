# Giraffe

This package is an extension of standard `net/http` package in
[Go](http://www.golang.org). It contains utilities that facilitate handling of
common HTTP operations.

![alt text](https://i.imgsafe.org/2fc6802cf9.png "Giraffe Logo")

# Installation

```sh
$ go get github.com/svett/giraffe
```

# Usage

You can encode an object as JSON in simplified way:

```Go
encoder := giraffe.NewHTTPEncoder(responseWriter)
encoder.EncodeJSON(map[string]string{"username": "root", "password": "swordfish"})
```

It can be encoded with padding as well:

```Go
encoder := giraffe.NewHTTPEncoder(responseWriter)
encoder.EncodeJSONP("login", map[string]string{"username": "root", "password": "swordfish"})
```

A similar operation can be performed for a byte array:

```Go
encoder := giraffe.NewHTTPEncoder(responseWriter)
encoder.EncodeData([]byte("gopher))
```

A plain text can be encoded as well:

```Go
encoder := giraffe.NewHTTPEncoder(responseWriter)
encoder.EncodeText("Hello World")
```

*MIT License*
