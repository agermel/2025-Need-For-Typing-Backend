from decouple import config
import grpc
from concurrent import futures
import sys
import os

sys.path.append(
    os.path.join(os.path.dirname(__file__), "..", "grpc", "python")
)

import llm_pb2
import llm_pb2_grpc
from langchain_community.chat_models import ChatOpenAI
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler

class GenerateService(llm_pb2_grpc.GenerateServicer):
    def GenerateStream(self,request,context):
        prompt = f"please write an essay about {request.topic} for 200 words in one paragraph. The generate response can only include one paragraph."
        llm = ChatOpenAI(
            model="gpt-4o-mini",
            openai_api_key=config("OPENAI_API_KEY"),
            base_url="https://api.chatanywhere.tech/v1",
            streaming=True,
            callbacks=[StreamingStdOutCallbackHandler()]
        )

        full_response = ""
        for chunk in llm.stream([{"role": "user", "content": prompt}]):
            full_response += chunk.content
            yield llm_pb2.GenerateResponse(response=full_response)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    llm_pb2_grpc.add_GenerateServicer_to_server(GenerateService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("GenerateService 已启动，监听端口 50051")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
        