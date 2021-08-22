package main

import (
	"context"
	"flag"
	pb "github.com/SarathLUN/my-grpc-go-route-guide/route-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"io"
	"log"
	"math/rand"
	"time"
)

var (
	tls                 = flag.Bool("tls", false, "connection use TLS if true, else plain TCP")
	caFile              = flag.String("ca_file", "", "the file containing the CA root cert file")
	serverAddr          = flag.String("server_addr", "localhost:10000", "the server in the format of host:port")
	serverHhostOverride = flag.String("server_host_override", "x.test.example.com", "the server use to verify host returned by the TLS handshake")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHhostOverride)
		if err != nil {
			log.Fatalf("failed to create TLS credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRouteGuideClient(conn)
	// looking for a valid feature
	//printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	// feature missing
	//printFeature(client, &pb.Point{Latitude: 0, Longitude: 0})

	// Looking for features between 40, -75 and 42, -73.
	//printFeatures(client, &pb.Rectangle{
	//	Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
	//	Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	//})

	// RecordRoute
	runRecordRoute(client)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("looking for features with %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeature(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeature(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeature(_) = _, %v", client, err)
		}
		log.Printf("Feature: name=%q, point=(%v, %v)", feature.GetName(), feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}
}

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v", client, err)
	}
	log.Println(feature)
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client pb.RouteGuideClient) {
	// create a random number of random point
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("traversing %d points", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRev() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route Summary: %v", reply)
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}
