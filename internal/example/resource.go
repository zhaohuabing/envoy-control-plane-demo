// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package example

import (
	"time"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
)

const (
	AllClusterName     = "echo-service"
	Subset1ClusterName = "echo-service-subset1"
	Subset2ClusterName = "echo-service-subset2"
	RouteName          = "echo-service"
	ListenerName       = "test-listener"
	ListenerPort       = 8080
	UpstreamHost1      = "instance1"
	UpstreamHost2      = "instance2"
	UpstreamHost3      = "instance3"
	UpstreamHost4      = "instance4"
	UpstreamPort       = 5678
)

func makeCluster() []types.Resource {
	return []types.Resource{
		&cluster.Cluster{
			Name:                 AllClusterName,
			ConnectTimeout:       durationpb.New(5 * time.Second),
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STRICT_DNS},
			LbPolicy:             cluster.Cluster_ROUND_ROBIN,
			LoadAssignment: makeEndpoint(AllClusterName, []string{UpstreamHost1,
				UpstreamHost2, UpstreamHost3, UpstreamHost4}),
			DnsLookupFamily: cluster.Cluster_V4_ONLY,
		},
		&cluster.Cluster{
			Name:                 Subset1ClusterName,
			ConnectTimeout:       durationpb.New(5 * time.Second),
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STRICT_DNS},
			LbPolicy:             cluster.Cluster_ROUND_ROBIN,
			LoadAssignment: makeEndpoint(Subset1ClusterName, []string{UpstreamHost1,
				UpstreamHost2}),
			DnsLookupFamily: cluster.Cluster_V4_ONLY,
		},
		&cluster.Cluster{
			Name:                 Subset2ClusterName,
			ConnectTimeout:       durationpb.New(5 * time.Second),
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STRICT_DNS},
			LbPolicy:             cluster.Cluster_ROUND_ROBIN,
			LoadAssignment: makeEndpoint(Subset2ClusterName, []string{UpstreamHost3,
				UpstreamHost4}),
			DnsLookupFamily: cluster.Cluster_V4_ONLY,
		},
	}
}

func makeEndpoint(clusterName string, upstreamHosts []string) *endpoint.ClusterLoadAssignment {
	endpoints := &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{},
		}},
	}
	for i := 0; i < len(upstreamHosts); i++ {
		endpoints.Endpoints[0].LbEndpoints = append(endpoints.Endpoints[0].LbEndpoints, &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{
				Endpoint: &endpoint.Endpoint{
					Address: &core.Address{
						Address: &core.Address_SocketAddress{
							SocketAddress: &core.SocketAddress{
								Address: upstreamHosts[i],
								PortSpecifier: &core.SocketAddress_PortValue{
									PortValue: UpstreamPort,
								},
							},
						},
					},
				},
			},
		})
	}
	return endpoints
}

func makeRoute(routeName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name: routeName,
		VirtualHosts: []*route.VirtualHost{{
			Name:    "echo-service",
			Domains: []string{"*"},
			Routes: []*route.Route{
				{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/subset1",
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: Subset1ClusterName,
							},
						},
					},
				},
				{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/subset2",
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: Subset2ClusterName,
							},
						},
					},
				},
				{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/",
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: "echo-service",
							},
						},
					},
				}},
		}},
	}
}

func makeHTTPListener(listenerName string, route string) *listener.Listener {
	routerConfig, _ := anypb.New(&router.Router{})
	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: route,
			},
		},
		HttpFilters: []*hcm.HttpFilter{{
			Name:       wellknown.Router,
			ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: routerConfig},
		}},
	}
	pbst, err := anypb.New(manager)
	if err != nil {
		panic(err)
	}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{

					Address: "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
	}
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds-server"},
				},
			}},
		},
	}
	return source
}

func GenerateSnapshot() *cache.Snapshot {
	snap, _ := cache.NewSnapshot("1",
		map[resource.Type][]types.Resource{
			resource.ClusterType:  makeCluster(),
			resource.RouteType:    {makeRoute(RouteName)},
			resource.ListenerType: {makeHTTPListener(ListenerName, RouteName)},
		},
	)
	return snap
}
