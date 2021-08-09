#include <iostream>
#include <memory>
#include <string>
using namespace std;

#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>

#include "br.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

using pb::LotteryQuery;
using pb::LotteryRequest;
using pb::LotteryResponse;


class LotteryQueryService final : public LotteryQuery::Service 
{
	Status FetchLotteryInfo(ServerContext* context, const LotteryRequest* req, LotteryResponse* res) override 
	{
		string prefix("Hello ");
		//res->set_message(prefix + req->name());

		return Status::OK;
	}
};


void RunServer() 
{
	string server_address("0.0.0.0:8090");
	LotteryQueryService service;

	grpc::EnableDefaultHealthCheckService(true);
	grpc::reflection::InitProtoReflectionServerBuilderPlugin();

	ServerBuilder builder;
	builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
	builder.RegisterService(&service);
	
	unique_ptr<Server> server(builder.BuildAndStart());
	cout << "Server listening on " << server_address << std::endl;

	server->Wait();
}

int main(int argc, char** argv) 
{
	RunServer();

	return 0;
}

