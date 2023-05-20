import sys
from datetime import datetime, timedelta
import random
import time

# Sample data for street segment speeds.
segment_speeds = {
    "8f4827ebed3c2e66f50daef967d5e91daadd8d98": {
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
    "df089962c85c67603f2e90c969931d72ab02c1ee": {
        "year": 2020,
        "month": 1,
        "day": 1,
        "hour": 1,
        "utc_timestamp": "2020-01-01T09:00:00.000Z",
        "start_junction_id": "5bea0e8381e051830525c0aba0141bde108dc02d",
        "end_junction_id": "2278ad9374ec96c35a0d769bc8a275f6355b55da",
        "osm_way_id": 40722998,
        "osm_start_node_id": 5780849015,
        "osm_end_node_id": 4927951349,
        "speed_mph_mean": 25.459,
        "speed_mph_stddev": 3.585,
    },
    "f32dbf217023581f429d56330be2a16410bc2809": {
        "year": 2020,
        "month": 1,
        "day": 30,
        "hour": 8,
        "utc_timestamp": "2020-01-30T16:00:00.000Z",
        "start_junction_id": "8aaf6ad421333ad741cb1d5de1e3fa83e7d0e908",
        "end_junction_id": "8bbd97259361a23d374e60e6582019550fc58e0f",
        "osm_way_id": 417094233,
        "osm_start_node_id": 4714793573,
        "osm_end_node_id": 1014244233,
        "speed_mph_mean": 27.761,
        "speed_mph_stddev": 3.679,
    },
}


def get_speed(segment: str, time_interval: float) -> int:
    """
    Simulates the speed and timestamp for a given street segment.

    Args:
        segment_id (str): The ID of the street segment.
        time_interval (float): The time interval between speed updates in seconds.

    Returns:
        dict or None: A dictionary containing speed information for the segment,
                      or None if speed information is not available for the segment.
    """
    if segment in segment_speeds:
        segment_speed = segment_speeds[segment]
        # Simulate some random variation in speed.
        base_speed = segment_speed["speed_mph_mean"]
        speed_variation = random.uniform(-1, 1) * segment_speed["speed_mph_stddev"]
        speed = base_speed + speed_variation

        # Simulate some random variation in speed standard deviation.
        segment_speed["speed_mph_stddev"] = segment_speed[
            "speed_mph_stddev"
        ] * random.uniform(0.9, 1.1)

        # Update the timestamp.
        timestamp = datetime.strptime(
            segment_speed["utc_timestamp"], "%Y-%m-%dT%H:%M:%S.%fZ"
        )
        updated_timestamp = timestamp + timedelta(seconds=time_interval)
        segment_speed["utc_timestamp"] = updated_timestamp.strftime(
            "%Y-%m-%dT%H:%M:%S.%fZ"
        )

        # Update the year, month, day, and hour based on the updated timestamp.
        segment_speed["year"] = updated_timestamp.year
        segment_speed["month"] = updated_timestamp.month
        segment_speed["day"] = updated_timestamp.day
        segment_speed["hour"] = updated_timestamp.hour

        speed_update = {
            "year": segment_speed["year"],
            "month": segment_speed["month"],
            "day": segment_speed["day"],
            "hour": segment_speed["hour"],
            "utc_timestamp": segment_speed["utc_timestamp"],
            "start_junction_id": segment_speed["start_junction_id"],
            "end_junction_id": segment_speed["end_junction_id"],
            "osm_way_id": segment_speed["osm_way_id"],
            "osm_start_node_id": segment_speed["osm_start_node_id"],
            "osm_end_node_id": segment_speed["osm_end_node_id"],
            "speed_mph_mean": speed,
            "speed_mph_stddev": segment_speed["speed_mph_stddev"],
        }

        return speed_update
    return None


def send_speed_updates(time_interval: float) -> None:
    """
    Simulates sending real-time speed updates for street segments.
    Continuously generates speed updates and prints the current speed for each segment.

    Args:
        time_interval (float): The time interval between speed updates in seconds.
    """
    while True:
        for segment, _ in segment_speeds.items():
            speed_update = get_speed(segment, time_interval)
            if speed_update is not None:
                print(speed_update)
            else:
                print("Speed information not available for the given segment.")
        time.sleep(time_interval)


if __name__ == "__main__":
    time_interval = 5

    try:
        # Check and parse if a time interval was provided as a command line argument.
        if len(sys.argv) > 1:
            time_interval = float(sys.argv[1])
            if time_interval <= 0:
                raise ValueError
    except ValueError:
        print("Invalid time interval. Please enter a positive number.")
        sys.exit(1)

    send_speed_updates(time_interval)
