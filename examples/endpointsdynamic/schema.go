package endpointsdynamic

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

var nodeDefinitions *relay.NodeDefinitions
var siteType *graphql.Object
var endpointType *graphql.Object
var queryType *graphql.Object

var resolverMap = map[string]func(string) interface{}{
	"Site":     GetSite,
	"Endpoint": GetEndpoint,
}

func siteIDFetcher(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
	if s, ok := obj.(*Site); ok {
		return s.SiteID, nil
	}
	return "", errors.New("Not valid type")
}

func endpointIDFetcher(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
	if s, ok := obj.(Endpoint); ok {
		return s.EndpointID, nil
	}
	if s, ok := obj.(*Endpoint); ok {
		return s.EndpointID, nil
	}
	return "", errors.New("Not valid type")
}

// Schema is exported, defined in init()
var Schema graphql.Schema

func init() {

	siteType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Site",
		Description: "site",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(*Site)
			return ok
		},
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Site", siteIDFetcher),
			"clientId": &graphql.Field{
				Type: graphql.String,
			},
			"siteId": &graphql.Field{
				Type: graphql.String,
			},
			"siteCode": &graphql.Field{
				Type: graphql.String,
			},
			"siteName": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"sites": &graphql.Field{
				Type: siteType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetSites(), nil
				},
			},
		},
	})

	/**
	 * Finally, we construct our schema (whose starting query type is the query
	 * type we defined above) and export it.
	 */
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
		//Mutation: mutationType,
	})
	if err != nil {
		// panic if there is an error in schema
		panic(err)
	}
}

// UpdateSchema with dynamic Endpoint type. Site fetch need to paginate endpoints
func UpdateSchema() {

	/**
	 * We get the node interface and field from the relay library.
	 *
	 * The first method is the way we resolve an ID to its object. The second is the
	 * way we resolve an object that implements node to its type.
	 */
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			// resolve id from global id
			resolvedID := relay.FromGlobalID(id)
			if f, ok := resolverMap[resolvedID.Type]; ok {
				return f(resolvedID.ID), nil
			}
			return nil, errors.New("Unknown node type")
		},
	})

	endpointType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Endpoint",
		Description: "endpoint",
		IsTypeOf: func(p graphql.IsTypeOfParams) bool {
			_, ok := p.Value.(Endpoint)
			_, ok2 := p.Value.(*Endpoint)
			return ok || ok2
		},
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Endpoint", endpointIDFetcher),
			"clientId": &graphql.Field{
				Type: graphql.String,
			},
			"siteId": &graphql.Field{
				Type: graphql.String,
			},
			"endpointId": &graphql.Field{
				Type: graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	endpointConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "Endpoint",
		NodeType: endpointType,
	})

	siteType.AddFieldConfig("endpoints", &graphql.Field{
		Type: endpointConnectionDefinition.ConnectionType,
		Args: relay.ConnectionArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// convert args map[string]interface into ConnectionArguments
			args := relay.NewConnectionArguments(p.Args)

			endpoints := []interface{}{}
			if site, ok := p.Source.(*Site); ok {
				for _, e := range site.Endpoints {
					endpoints = append(endpoints, e)
				}
			}
			return relay.ConnectionFromArray(endpoints, args), nil
		},
	})
	queryType.AddFieldConfig("endpoints", &graphql.Field{
		Type: endpointType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return GetEndpoints(), nil
		},
	})
	queryType.AddFieldConfig("node", nodeDefinitions.NodeField)

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		// panic if there is an error in schema
		panic(err)
	}
}
