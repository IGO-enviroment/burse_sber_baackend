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
	AuthorizationScopes = "authorization.Scopes"
	OriginScopes        = "origin.Scopes"
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

	ctx = context.WithValue(ctx, AuthorizationScopes, []string{"Bearer"})

	ctx = context.WithValue(ctx, OriginScopes, []string{})

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

	ctx = context.WithValue(ctx, AuthorizationScopes, []string{"Bearer"})

	ctx = context.WithValue(ctx, OriginScopes, []string{})

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

	"H4sIAAAAAAAC/6xWzW7jNhB+FWHaUyBbThMUgW5224PRNi2S5mT4QItjm6lEsiSVxDUEtOf2DfoSRYEe",
	"dhe7zyC/0YKUrB9LdgzsXgxJM/zmm5mPM95CJBIpOHKjIdyCJIokaFC5N0btL0UdKSYNExxCmH4LPjD7",
	"JIlZgw+cJAih9fVB4W8pU0ghNCpFH3S0xoRYkKVQCTHWjxvwwWwkFi+4QgVZlu2dXWBC6b1JKXLjSCkh",
	"URmGzqYLQ0HQYOIevlS4hBC+COp0ghIvKA9AVoUlSpGNfU85e0KlmdlMXaodWs2UZm13v2Yyr5DF4hEj",
	"FypSSGzJ7lCncU8aTH9jPbAZdyFEjITb4yX2dwlhccNDG8X4qkOs5e03sPuYxWJl+3dICPeh2v1+0Kg8",
	"LIEPWPggidbPQtEj5yqz/0oC+wDVgT7i+pgk8EiVfFgypc2tU2iPNSYnjAmjNMajZrkWHG/TZIHq9f7s",
	"06vpNIK3QrWB+6pgxK/I71AKrrFbChJFqPUv1qeXNr5IplBPeY/cfXgZrMSg/vr1dSeTJn4TrcvUNgyj",
	"VDGzubc3seSXmrVQ7HdSyGRbjJI1EuoIlMNk3PKq761k3+OmmBaML0VXdT+msWGDJYvRWwgWo5IxMegt",
	"hfJ+ksjHP0+9e4kRW7LIYQ8tODOxRZ/UB8AHd9Ed5uVwNBzZ2gmJnEgGIVwNR8MrJ1ezdmkFNq+gvllC",
	"my65H5zZ4SgX3Q6d6qutMmozEXRjj0aCm1LsF8GFG3zVMD017QoOWVY0TjuZOIpfjUafgtuSnUPvJLdC",
	"ast0XQRqmyeEemWGLWlAOOuIYgYTJMrKP/O3IBRzVZ3NMysystLWwyoE5hYpqKdyQCgdNPdDfx/GlHoc",
	"n73K87Al42r96JONIVLGpZCCR13o+bxqNhbcWa06HemsNXiwkDrbsNvSfQ08Qum5rbU+l12fW2G8B77v",
	"swX7DBp4qBpvleAA1ZP76zI7jJ//k3/I3+b/5u92f+fvd3/lb7zdn/n/uz/y/+wv+JCq2A4iY2QYBLGI",
	"SLwW2oQ3o5tR8HQJ2Tz7GAAA///OPbHFLAkAAA==",
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
