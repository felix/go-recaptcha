// Package recaptcha interacts with the reCaptcha verification API
package recaptcha

import (
    "errors"
    "net/http"
    "net/url"
    "encoding/json"
)

const (
    verifyURL = "https://www.google.com/recaptcha/api/siteverify"
)

var (
    ErrSecretMissing   = errors.New("The secret parameter is missing.")
    ErrSecretInvalid   = errors.New("The secret parameter is invalid or malformed.")
    ErrResponseMissing = errors.New("The response parameter is missing.")
    ErrResponseInvalid = errors.New("The response parameter is invalid or malformed.")
)

type Response struct {
    Success     bool     `json:"success"`
    // Errors from reCaptcha service
    ErrorCodes  []string `json:"error-codes"`
    // Our errors to caller
    Errors      []error  `json:"-"`
}

// Translate error code into error
func (r *Response) getErrorFromCode(code string) error {
    switch code {
    case "missing-input-secret":
        return ErrSecretMissing
    case "invalid-input-secret":
        return ErrSecretInvalid
    case "missing-input-response":
        return ErrResponseMissing
    case "invalid-input-response":
        return ErrResponseInvalid
    default:
        return nil
    }
}

// Translate reCaptcha error codes into errors
func (r *Response) populateErrors() {
    for _, code := range r.ErrorCodes {
        r.Errors = append(r.Errors, r.getErrorFromCode(code))
    }
}

// Public verfication function
// It requires the reCaptcha secret, the response from the widget.
// It returns the reCaptcha response as a struct.
func Verify(secret string, response string, remoteip string) (*Response, error) {

    values := make(url.Values)
    values.Set("secret", secret)
    values.Set("response", response)
    values.Set("remoteip", remoteip)

    resp, err := http.PostForm(verifyURL, values)
    if err != nil {
        return nil, err

    }

    defer resp.Body.Close()

    // Collect the response into struct
    recaptchaResponse := new(Response)
    err = json.NewDecoder(resp.Body).Decode(recaptchaResponse); if err != nil {
        return nil, err
    }

    recaptchaResponse.populateErrors()

    return recaptchaResponse, nil
}
