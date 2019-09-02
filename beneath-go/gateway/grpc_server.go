package gateway

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/beneath-core/beneath-go/core/timeutil"

	"github.com/beneath-core/beneath-go/control/auth"
	"github.com/beneath-core/beneath-go/control/model"
	"github.com/beneath-core/beneath-go/core/jsonutil"
	"github.com/beneath-core/beneath-go/core/queryparse"
	"github.com/beneath-core/beneath-go/db"
	pb "github.com/beneath-core/beneath-go/proto"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientVersionSpec struct {
	DeprecatedVersion  string
	WarningVersion     string
	RecommendedVersion string
}

func (s clientVersionSpec) IsZero() bool {
	return s.DeprecatedVersion == "" || s.WarningVersion == "" || s.RecommendedVersion == ""
}

var clientSpecs = map[string]clientVersionSpec{
	"beneath-python": clientVersionSpec{
		DeprecatedVersion:  "0.0.0",
		WarningVersion:     "0.0.0",
		RecommendedVersion: "0.0.1",
	},
}

// ListenAndServeGRPC serves a gRPC API
func ListenAndServeGRPC(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(auth.GRPCInterceptor),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_auth.StreamServerInterceptor(auth.GRPCInterceptor),
			grpc_recovery.StreamServerInterceptor(),
		),
	)
	pb.RegisterGatewayServer(server, &gRPCServer{})

	log.Printf("gRPC server running on port %d\n", port)
	return server.Serve(lis)
}

// gRPCServer implements pb.GatewayServer
type gRPCServer struct{}

func (s *gRPCServer) GetStreamDetails(ctx context.Context, req *pb.StreamDetailsRequest) (*pb.StreamDetailsResponse, error) {
	// get auth
	secret := auth.GetSecret(ctx)

	// get instance ID
	instanceID := model.FindInstanceIDByNameAndProject(ctx, req.StreamName, req.ProjectName)
	if instanceID == uuid.Nil {
		return nil, status.Error(codes.NotFound, "stream not found")
	}

	// get stream details
	stream := model.FindCachedStreamByCurrentInstanceID(ctx, instanceID)
	if stream == nil {
		return nil, status.Error(codes.NotFound, "stream not found")
	}

	// check permissions
	if !secret.ReadsProject(stream.ProjectID) {
		return nil, grpc.Errorf(codes.PermissionDenied, "token doesn't grant right to read this stream")
	}

	// return
	return &pb.StreamDetailsResponse{
		CurrentInstanceId: instanceID.Bytes(),
		ProjectId:         stream.ProjectID.Bytes(),
		ProjectName:       stream.ProjectName,
		StreamName:        stream.StreamName,
		KeyFields:         stream.Codec.GetKeyFields(),
		AvroSchema:        stream.Codec.GetAvroSchemaString(),
		Public:            stream.Public,
		External:          stream.External,
		Batch:             stream.Batch,
		Manual:            stream.Manual,
	}, nil
}

func (s *gRPCServer) ReadRecords(ctx context.Context, req *pb.ReadRecordsRequest) (*pb.ReadRecordsResponse, error) {
	// get auth
	secret := auth.GetSecret(ctx)

	// read instanceID
	instanceID := uuid.FromBytesOrNil(req.InstanceId)
	if instanceID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "instance_id not valid UUID")
	}

	// get cached stream
	stream := model.FindCachedStreamByCurrentInstanceID(ctx, instanceID)
	if stream == nil {
		return nil, status.Error(codes.NotFound, "stream not found")
	}

	// check permissions
	if !secret.ReadsProject(stream.ProjectID) {
		return nil, grpc.Errorf(codes.PermissionDenied, "token doesn't grant right to read from this stream")
	}

	// check limit is valid
	if req.Limit == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "limit cannot be 0")
	} else if req.Limit > maxRecordsLimit {
		return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprintf("limit exceeds maximum of %d", maxRecordsLimit))
	}

	// get key range where clause (if no where, it will be nil)
	var whereQuery queryparse.Query
	if req.Where != "" {
		// read where
		var where map[string]interface{}
		err := jsonutil.UnmarshalBytes([]byte(req.Where), &where)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "couldn't parse 'where' -- is it valid JSON?")
		}

		// make query
		whereQuery, err = queryparse.JSONToQuery(where)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprintf("couldn't parse 'where' query: %s", err.Error()))
		}
	}

	// adapt key range based on after (for pagination), if present
	var afterQuery queryparse.Query
	if req.After != "" {
		// read after
		var after map[string]interface{}
		err := jsonutil.UnmarshalBytes([]byte(req.After), &after)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "couldn't parse 'after' -- is it valid JSON?")
		}

		// make query
		afterQuery, err = queryparse.JSONToQuery(after)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprintf("couldn't parse 'after' query: %s", err.Error()))
		}
	}

	// set key range
	keyRange, err := stream.Codec.MakeKeyRange(whereQuery, afterQuery)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	// read rows from engine
	response := &pb.ReadRecordsResponse{}
	err = db.Engine.Tables.ReadRecordRange(ctx, instanceID, keyRange, int(req.Limit), func(avroData []byte, timestamp time.Time) error {
		response.Records = append(response.Records, &pb.Record{
			AvroData:  avroData,
			Timestamp: timeutil.UnixMilli(timestamp),
		})
		return nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	// done
	return response, nil
}

func (s *gRPCServer) ReadLatestRecords(ctx context.Context, req *pb.ReadLatestRecordsRequest) (*pb.ReadRecordsResponse, error) {
	// get auth
	secret := auth.GetSecret(ctx)

	// read instanceID
	instanceID := uuid.FromBytesOrNil(req.InstanceId)
	if instanceID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "instance_id not valid UUID")
	}

	// get cached stream
	stream := model.FindCachedStreamByCurrentInstanceID(ctx, instanceID)
	if stream == nil {
		return nil, status.Error(codes.NotFound, "stream not found")
	}

	// check permissions
	if !secret.ReadsProject(stream.ProjectID) {
		return nil, grpc.Errorf(codes.PermissionDenied, "token doesn't grant right to read from this stream")
	}

	// check isn't batch
	if stream.Batch {
		return nil, grpc.Errorf(codes.InvalidArgument, "cannot get latest records for batch streams")
	}

	// check limit is valid
	if req.Limit == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "limit cannot be 0")
	} else if req.Limit > maxRecordsLimit {
		return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprintf("limit exceeds maximum of %d", maxRecordsLimit))
	}

	// get before as time
	var before time.Time
	if req.Before != 0 {
		before = timeutil.FromUnixMilli(req.Before)
	} else {
		before = time.Time{}
	}

	// read rows from engine
	response := &pb.ReadRecordsResponse{}
	err := db.Engine.Tables.ReadLatestRecords(ctx, instanceID, int(req.Limit), before, func(avroData []byte, timestamp time.Time) error {
		response.Records = append(response.Records, &pb.Record{
			AvroData:  avroData,
			Timestamp: timeutil.UnixMilli(timestamp),
		})
		return nil
	})
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	// done
	return response, nil
}

func (s *gRPCServer) WriteRecords(ctx context.Context, req *pb.WriteRecordsRequest) (*pb.WriteRecordsResponse, error) {
	// get auth
	secret := auth.GetSecret(ctx)

	// read instanceID
	instanceID := uuid.FromBytesOrNil(req.InstanceId)
	if instanceID == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "instance_id not valid UUID")
	}

	// get stream info
	stream := model.FindCachedStreamByCurrentInstanceID(ctx, instanceID)
	if stream == nil {
		return nil, status.Error(codes.NotFound, "stream not found")
	}

	// check permissions
	if !secret.WritesStream(stream) {
		return nil, grpc.Errorf(codes.PermissionDenied, "token doesn't grant right to write to this stream")
	}

	// check records supplied
	if len(req.Records) == 0 {
		return nil, status.Error(codes.InvalidArgument, "records cannot be empty")
	}

	// check each record is valid
	for idx, record := range req.Records {
		// set timestamp to current timestamp if it's 0
		if record.Timestamp == 0 {
			record.Timestamp = timeutil.UnixMilli(time.Now())
		}

		// check it decodes
		decodedData, err := stream.Codec.UnmarshalAvro(record.AvroData)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("record at index %d doesn't decode: %v", idx, err.Error()))
		}

		// compute key (only used to check size)
		keyData, err := stream.Codec.MarshalKey(decodedData)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error encoding record at index %d: %v", idx, err.Error()))
		}

		// check size
		err = db.Engine.CheckSize(len(keyData), len(record.AvroData))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("error encoding record at index %d: %v", idx, err.Error()))
		}
	}

	// write request to engine
	err := db.Engine.Streams.QueueWriteRequest(ctx, req)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
	}

	return &pb.WriteRecordsResponse{}, nil
}

func (s *gRPCServer) SendClientPing(ctx context.Context, req *pb.ClientPing) (*pb.ClientPong, error) {
	spec := clientSpecs[req.ClientId]
	if spec.IsZero() {
		return nil, grpc.Errorf(codes.InvalidArgument, "unrecognized client ID")
	}

	status := "deprecated"
	if req.ClientVersion == spec.RecommendedVersion {
		status = "stable"
	}

	secret := auth.GetSecret(ctx)
	return &pb.ClientPong{
		Authenticated:      !secret.IsAnonymous(),
		Status:             status,
		RecommendedVersion: spec.RecommendedVersion,
	}, nil
}
