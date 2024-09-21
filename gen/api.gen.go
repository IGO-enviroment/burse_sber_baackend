// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package gen

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

const (
	OriginScopes = "origin.Scopes"
)

// AddStudent defines model for addStudent.
type AddStudent struct {
	Students     []Student `json:"students"`
	UniversityId int       `json:"universityId"`
}

// CreationResult defines model for creationResult.
type CreationResult struct {
	IsCreated    bool   `json:"isCreated"`
	StudentEmail string `json:"studentEmail"`
}

// Login defines model for login.
type Login struct {
	// Email User email
	Email string `json:"email"`

	// Password User password
	Password string `json:"password"`
}

// Student defines model for student.
type Student struct {
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	MiddleName  string `json:"middleName"`
	PhoneNumber string `json:"phoneNumber"`
}

// TokenReponse defines model for tokenReponse.
type TokenReponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

// Id defines model for id.
type Id = int

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = Login

// AddStudentsJSONRequestBody defines body for AddStudents for application/json ContentType.
type AddStudentsJSONRequestBody = AddStudent

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /auth/login)
	Login(w http.ResponseWriter, r *http.Request)

	// (POST /university/add-students)
	AddStudents(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, OriginScopes, []string{"*"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// AddStudents operation middleware
func (siw *ServerInterfaceWrapper) AddStudents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.AddStudents(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = mux.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc(options.BaseURL+"/auth/login", wrapper.Login).Methods("POST")

	r.HandleFunc(options.BaseURL+"/university/add-students", wrapper.AddStudents).Methods("POST")

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RWzW7bRhB+FWLaU0GJcm0UBm9y24PQ1i7s+iTosOKOpHXJ3c3u0rYiEEjOyRvkJYIA",
	"OSRB8gzUGwW7pPgjUrYPzkXQav6+me/bWW0gEokUHLnREG5AEkUSNKjciVH7SVFHiknDBIcQJn+AD8x+",
	"k8SswAdOEoTQ+vqg8EXKFFIIjUrRBx2tMCE2yUKohBjrxw34YNYSiwMuUUGWZTtnV5hQemVSitw4UEpI",
	"VIahs+nCUAA0mLgvPytcQAg/BXU7QZkvKAMgq8oSpcjanlPOblFpZtYT12oHVrOladvdr5HMqsxifoOR",
	"KxUpJHZkl6jTuKcNpn+3HtisOxciRsJteJn7z4SwuOGhjWJ82QHW8vYbufuQxWJp+dsHhLtSbb6vNSoP",
	"y8R7KHyQROs7oeiBuMrsP9LArkAV0AdcH5IEHpiSDwumtDl3Cu2xxuQBY8IojfGgWa4Ex/M0maN6nJ9d",
	"ezWcRvFWqXbivikY8T/yS5SCa+yOgkQRav2f9emFjfeSKdQT3iN3H+4HSzGof/3tpNNJM38zWxepJQyj",
	"VDGzvrI3scSXmpVQ7CUpZLIpVskKCXUAymUybnnV91ayv9BdXKFYKeLe+IvC3Am0oBhfiK5c/0ljwwYL",
	"FqM3FyxGJWNi0FsI5V1I5ON/J96VxIgtWORADW1yZmKb/awOAB/chnA5j4aj4cihlciJZBDC8XA0PHY6",
	"Nys3j8AOJKivpNCmC+5vUXRjiXbV7baqfrX0oDZngq5taCS4KW8JkTIu8QY3uph3vZIf2pkFoCwr6NdO",
	"bA7vr6PRsxVpKdmV6rS9RGoHeFJUbZvPCPXK3ltqg3C6qQQyhV9glll1kqW2RystmFn/oF7nAaF00HxY",
	"+nkYU+pxvPMqz31KxtW7pX8QMY2X8RnYedL7ufeSdZ7RLnG7GXiE0qcSaH2Ouj7nwnjXfLc1bLKs5vK6",
	"ItAy6iSgbt1/l+l+nvxd/i3/nL/Pv2zf5l+3b/JP3vZ1/nH7Kv9gP8GHVMV2kxgjwyCIRUTildAmPB2d",
	"joLbI8hm2fcAAAD//zNfcoUtCQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
