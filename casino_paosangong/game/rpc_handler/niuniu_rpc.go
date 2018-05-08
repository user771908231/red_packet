package rpc_handler

import (
	"google.golang.org/grpc"
	"net"
	"casino_common/proto/ddproto"
	"golang.org/x/net/context"
	"casino_paosangong/service/paosangongService"
)

type NiuRpcSrv struct{}

//监听rpc服务
func LisenAndServeNiuniuRpc(addr string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := grpc.NewServer()
	ddproto.RegisterNiuniuRpcServer(s, ddproto.NiuniuRpcServer(NiuRpcSrv{}))
	go s.Serve(lis)
	return s, nil
}

//创建房间
func (r NiuRpcSrv) CreateRoom(ctx context.Context,req *ddproto.NiuCreateDeskReq) (*ddproto.NiuEnterDeskAck, error) {
	return paosangongService.CreateDeskHandler(req, nil), nil
}
