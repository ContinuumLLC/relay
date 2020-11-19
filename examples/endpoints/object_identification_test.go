package endpoints_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/graphql-go/relay/examples/endpoints"
)

func TestObjectIdentification_TestFetching_CorrectlyFetchesTheIDAndTheNameOfThesites(t *testing.T) {
	query := `
        query sitesQuery {
          sites {
            id
            siteId
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"id":     "U2l0ZTo2NTAwMDIz",
				"siteId": "6500023",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpoints.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestObjectIdentification_TestFetching_CorrectlyRefetchesThesites(t *testing.T) {
	query := `
	query sitesRefetchQuery {
		node(id: "U2l0ZTo2NTAwMDIz") {
		  id
		  ... on Site {
			siteId
		  }
		}
	  }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":     "U2l0ZTo2NTAwMDIz",
				"siteId": "6500023",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpoints.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestObjectIdentification_TestFetching_CorrectlyRefetchesTheEndpoint(t *testing.T) {
	query := `
	query {
		node(id: "RW5kcG9pbnQ6Mzc4NjY1MjQtY2M5MS00ZDY0LWI1ZGItYjkxMmVhZjQzMzll") {
		  id
		  ... on Endpoint {
			endpointId
		  }
		}
	  }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":         "RW5kcG9pbnQ6Mzc4NjY1MjQtY2M5MS00ZDY0LWI1ZGItYjkxMmVhZjQzMzll",
				"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339e",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpoints.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
