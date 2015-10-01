package gqlrelay

import (
	"encoding/json"
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/gorilla/schema"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeGraphQL        = "application/graphql"
	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

var decoder = schema.NewDecoder()

type Handler struct {
	Schema *types.GraphQLSchema
	render *render.Render
}
type requestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getRequestOptions(r *http.Request) *requestOptions {

	query := r.URL.Query().Get("query")
	if query != "" {

		// get variables map
		var variables map[string]interface{}
		variablesStr := r.URL.Query().Get("variables")
		json.Unmarshal([]byte(variablesStr), variables)

		return &requestOptions{
			Query:         query,
			Variables:     variables,
			OperationName: r.URL.Query().Get("operationName"),
		}
	}
	if r.Method != "POST" {
		return &requestOptions{}
	}
	if r.Body == nil {
		return &requestOptions{}
	}

	// TODO: improve Content-Type handling
	contentTypeStr := r.Header.Get("Content-Type")
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]

	switch contentType {
	case ContentTypeGraphQL:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &requestOptions{}
		}
		return &requestOptions{
			Query: string(body),
		}
	case ContentTypeFormUrlEncoded:
		var opts requestOptions
		err := r.ParseForm()
		if err != nil {
			return &requestOptions{}
		}
		err = decoder.Decode(&opts, r.PostForm)
		if err != nil {
			return &requestOptions{}
		}
		return &opts
	case ContentTypeJSON:
		fallthrough
	default:
		var opts requestOptions
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &opts
		}
		err = json.Unmarshal(body, &opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			json.Unmarshal(body, &optsCompatible)
			json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return &opts
	}
}

// ServeHTTP provides an entry point into executing graphQL queries
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get query
	opts := getRequestOptions(r)

	// execute graphql query
	resultChannel := make(chan *types.GraphQLResult)
	params := gql.GraphqlParams{
		Schema:         *h.Schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
	}
	go gql.Graphql(params, resultChannel)
	result := <-resultChannel

	// render result
	h.render.JSON(w, http.StatusOK, result)
}

type HandlerConfig struct {
	Schema *types.GraphQLSchema
	Pretty bool
}

func NewHandlerConfig() *HandlerConfig {
	return &HandlerConfig{
		Schema: nil,
		Pretty: true,
	}
}

func NewHandler(p *HandlerConfig) *Handler {
	if p == nil {
		p = NewHandlerConfig()
	}
	if p.Schema == nil {
		panic("undefined graphQL schema")
	}
	r := render.New(render.Options{
		IndentJSON: p.Pretty,
	})
	return &Handler{
		Schema: p.Schema,
		render: r,
	}
}
