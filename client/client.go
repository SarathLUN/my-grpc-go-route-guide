package main

import (
	"context"
	"flag"
	pb "github.com/SarathLUN/my-grpc-go-route-guide/route-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"log"
	"time"
)

var (
	tls = flag.Bool("tls", false, "connection use TLS if true, else plain TCP")
	caFile = flag.String("ca_file","","the file containing the CA root cert file")
	serverAddr = flag.String("server_addr","localhost:10000", "the server in the format of host:port")
	serverHhostOverride = flag.String("server_host_override","x.test.example.com","the server use to verify host returned by the TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile=data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHhostOverride)
		if err != nil {
			log.Fatalf("failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v",err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)
	// looking for a valid feature
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	// feature missing
	printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})
}

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d %d)",point.Latitude,point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v",client,err)
	}
	log.Println(feature)
}

