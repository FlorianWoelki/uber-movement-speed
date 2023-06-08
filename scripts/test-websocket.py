import json
import boto3
import asyncio
import websockets


url = "ws://localhost:4510"


def main():
    async def start_client(uri: str):
        async with websockets.connect(uri) as websocket:
            print("Sending message to websocket")
            await websocket.send(
                json.dumps(
                    {
                        "action": "kinesis-data-forwarder",
                        "data": {"id": "ws-id", "title": "WS book title"},
                    }
                )
            )
            result = await websocket.recv()
            print(f"Received message from websocket: {result}")

    print("Connecting to websocket URL")
    asyncio.get_event_loop().run_until_complete(start_client(url))


if __name__ == "__main__":
    main()
