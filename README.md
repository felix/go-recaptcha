go-recaptcha
============

About
-----

This package handles reCaptcha (http://www.google.com/recaptcha) form
submissions in Go (http://golang.org/).

Usage
-----

Install the package in your environment using `go get`:

```
go get github.com/felix/go-recaptcha
```

or as part of your source if using [gb](http://getgb.io/).

To use it within your own code, import "github.com/felix/go-recaptcha" and call:

```
recaptcha.Verify(secret, response, remoteip)
```

for each reCaptcha form input you need to check, using the values obtained by
reading the form's POST parameters.

The recaptcha.Verify() function returns a Response struct defined as:

```
type Response struct {
    Success     bool     `json:"success"`
    // Errors from reCaptcha service
    ErrorCodes  []string `json:"error-codes"`
    // Our errors to caller
    Errors      []error  `json:"-"`
}
```
