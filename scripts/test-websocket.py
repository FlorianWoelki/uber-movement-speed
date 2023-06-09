import json
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
                        "data": {
                            "year": 2020,
                            "month": 1,
                            "day": 1,
                            "hour": 1,
                            "utc_timestamp": "2020-01-01T09:00:00.000Z",
                            "start_junction_id": "8e555723c3dff79036c7a8c0cef6b32a80763c9f",
                            "end_junction_id": "2278ad9374ec96c35a0d769bc8a275f6355b55da",
                            "osm_way_id": 40722998,
                            "osm_start_node_id": 62385707,
                            "osm_end_node_id": 4927951349,
                            "speed_mph_mean": 26.636,
                            "speed_mph_stddev": 4.483,
                        },
                    }
                )
            )
            result = await websocket.recv()
            print(f"Received message from websocket: {result}")

    print("Connecting to websocket URL")
    asyncio.get_event_loop().run_until_complete(start_client(url))


if __name__ == "__main__":
    main()
