package endpointsdynamic_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/graphql-go/relay/examples/endpointsdynamic"
)

func init() {
	endpointsdynamic.UpdateSchema()
}

func TestConnection_TestFetching_CorrectlyFetchesTheFirstShipOfThesites(t *testing.T) {
	query := `
        query {
          sites {
            siteId,
            endpoints(first: 1) {
              edges {
                node {
					endpointId
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"siteId": "6500023",
				"endpoints": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339e",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpointsdynamic.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesTheFirstTwoendpointsOfThesitesWithACursor(t *testing.T) {
	query := `
        query {
          sites {
            siteId,
            endpoints(first: 2) {
              edges {
                cursor,
                node {
					endpointId
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"siteId": "6500023",
				"endpoints": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjA=",
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339e",
							},
						},
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjE=",
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339u",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpointsdynamic.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesTheNextTwoEndpointsOfThesitesWithACursor(t *testing.T) {
	query := `
	query  {
		sites {
		  siteId,
		  endpoints(first: 2, after: "YXJyYXljb25uZWN0aW9uOjA=") {
			edges {
			  cursor,
			  node {
				endpointId
			  }
			}
		  }
		}
	  }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"siteId": "6500023",
				"endpoints": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjE=",
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339u",
							},
						},
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjI=",
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4338u",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpointsdynamic.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesNoendpointsOfThesitesAtTheEndOfTheConnection(t *testing.T) {
	query := `
	query  {
		sites {
		  siteId,
		  endpoints(first: 2, after: "YXJyYXljb25uZWN0aW9uOjQ=") {
			edges {
			  cursor,
			  node {
				endpointId
			  }
			}
		  }
		}
	  }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"siteId": "6500023",
				"endpoints": map[string]interface{}{
					"edges": []interface{}{},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpointsdynamic.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyIdentifiesTheEndOfTheList(t *testing.T) {
	query := `
	query  {
		sites {
		  siteId,
		  originalendpoints: endpoints(first: 2) {
			edges {
			  node {
				endpointId
			  }
			}
			pageInfo {
			  hasNextPage
			}
		  }
		  moreendpoints: endpoints(first: 1 after: "YXJyYXljb25uZWN0aW9uOjE=") {
			edges {
			  node {
				endpointId
			  }
			}
			pageInfo {
			  hasNextPage
			}
		  }
		}
	  }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"sites": map[string]interface{}{
				"siteId": "6500023",
				"originalendpoints": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339e",
							},
						},
						map[string]interface{}{
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4339u",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": true,
					},
				},
				"moreendpoints": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"endpointId": "37866524-cc91-4d64-b5db-b912eaf4338u",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        endpointsdynamic.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
