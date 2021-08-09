#include <iostream>
#include <memory>
#include <string>
using namespace std;

#include <grpcpp/grpcpp.h>

#include "br.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

using pb::LotteryQuery;
using pb::LotteryRequest;
using pb::LotteryResponse;


class LotteryQueryClient
{
public:
	LotteryQueryClient(std::shared_ptr<Channel> channel)
		: stub_(LotteryQuery::NewStub(channel)) 
	{
	}

	bool FetchLotteryInfo(LotteryResponse& res, int tp)
	{
		LotteryRequest req;
		req.set_type(tp);

		ClientContext context;

		Status status = stub_->FetchLotteryInfo(&context, req, &res);

		if (status.ok()) 
		{
			return true;
		}
		else 
		{
			std::cout << status.error_code() << ": " << status.error_message() << std::endl;

			return false;
		}
	}

private:
	std::unique_ptr<LotteryQuery::Stub>		stub_;
};

int main(int argc, char** argv) 
{
	string target_str;
	string arg_str("--target");

	if (argc > 1) 
	{
		string arg_val = argv[1];
		size_t start_pos = arg_val.find(arg_str);

		if (start_pos != std::string::npos) 
		{
			start_pos += arg_str.size();

			if (arg_val[start_pos] == '=') 
			{
				target_str = arg_val.substr(start_pos + 1);
			}
			else 
			{
				std::cout << "The only correct argument syntax is --target=" << std::endl;
				return 0;
			}
		}
		else 
		{
			std::cout << "The only acceptable argument is --target=" << std::endl;
			return 0;
		}
	}
	else 
	{
		target_str = "localhost:8090";
	}

	LotteryQueryClient client(grpc::CreateChannel(target_str, grpc::InsecureChannelCredentials()));

	LotteryResponse res;
	bool b = client.FetchLotteryInfo(res, 0);

	//cout << "Greeter received: " << res << std::endl;

	getchar();

	return 0;
}
