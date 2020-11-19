package endpoints

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
)

var nodeDefinitions *relay.NodeDefinitions
var siteType *graphql.Object
var endpointType *graphql.Object

// exported schema, defined in init()
var Schema graphql.Schema

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

func init() {

	/**s
	 * We get the node interface and field from the relay library.
	 *
	 * The first method is the way we resolve an ID to its object. The second is the
	 * way we resolve an object that implements node to its type.
	 */
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			// resolve id from global id
			resolvedID := relay.FromGlobalID(id)

			// based on id and its type, return the object
			switch resolvedID.Type {
			case "Site":
				return GetSite(resolvedID.ID), nil
			case "Endpoint":
				return GetEndpoint(resolvedID.ID), nil
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			// based on the type of the value, return GraphQLObjectType
			switch p.Value.(type) {
			case *Site:
				return siteType
			default:
				return endpointType
			}
		},
	})

	endpointType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Endpoint",
		Description: "endpoint",
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

	siteType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Site",
		Description: "site",
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
			"endpoints": &graphql.Field{
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
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"sites": &graphql.Field{
				Type: siteType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetSites(), nil
				},
			},
			"node": nodeDefinitions.NodeField,
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
