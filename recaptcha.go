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
    ErrSecretMissing   = errors.New("reCaptcha secret is missing.")
    ErrSecretInvalid   = errors.New("reCaptcha secret is invalid.")
    ErrResponseMissing = errors.New("reCaptcha response is missing.")
    ErrResponseInvalid = errors.New("reCaptcha response is invalid.")
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

func Verify(secret string, response string, remoteip string) (*Response, error) {

    values := make(url.Values)
    values.Set("secret", secret)
    values.Set("response", response)

    resp, err := http.Get(verifyURL + "?" + values.Encode())
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
