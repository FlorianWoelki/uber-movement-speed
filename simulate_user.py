import sys
from datetime import datetime, timedelta
import random
import time


segment_speeds = {
    "segment1": {
        "start_junction_id": "junction1",
        "end_junction_id": "junction2",
        "osm_way_id": "way1",
        "osm_start_node_id": "node1",
        "osm_end_node_id": "node2",
        "driver_id": "driver1",
    },
    "segment2": {
        "start_junction_id": "junction2",
        "end_junction_id": "junction3",
        "osm_way_id": "way2",
        "osm_start_node_id": "node2",
        "osm_end_node_id": "node3",
        "driver_id": "driver2",
    },
    "segment3": {
        "start_junction_id": "junction3",
        "end_junction_id": "junction4",
        "osm_way_id": "way3",
        "osm_start_node_id": "node3",
        "osm_end_node_id": "node4",
        "driver_id": "driver3",
    },
}


def get_speed(segment: str, driver_speeds: dict) -> int:
    """
    Simulates the speed and timestamp for a given street segment and driver.

    Args:
        segment_id (str): The ID of the street segment.
        driver_speeds (dict): Dictionary containing driver-specific speeds.

    Returns:
        dict or None: A dictionary containing speed information for the segment and driver,
                      or None if speed information is not available for the segment or driver.
    """
    if segment in segment_speeds:
        segment_speed = segment_speeds[segment]
        if segment_speed["driver_id"] in driver_speeds:
            driver_speed = driver_speeds[segment_speed["driver_id"]]

            # Simulate some random variation in speed.
            base_speed = driver_speed["speed_mph"]
            speed_variation = random.uniform(-1, 1) * driver_speed["speed_variation"]
            speed = base_speed + speed_variation

            # Simulate some random variation in speed standard deviation.
            driver_speed["speed_variation"] = driver_speed[
                "speed_variation"
            ] * random.uniform(0.9, 1.1)

            speed_update = {
                "start_junction_id": segment_speed["start_junction_id"],
                "end_junction_id": segment_speed["end_junction_id"],
                "osm_way_id": segment_speed["osm_way_id"],
                "osm_start_node_id": segment_speed["osm_start_node_id"],
                "osm_end_node_id": segment_speed["osm_end_node_id"],
                "driver_id": segment_speed["driver_id"],
                "speed_mph": speed,
                "speed_variation": driver_speed["speed_variation"],
            }

            return speed_update
    return None


def simulate_driving(time_interval: float, driver_speeds: dict) -> None:
    """
    Simulates a user driving on a street.
    Generates speed updates and prints the current speed for each segment and driver.

    Args:
        time_interval (float): The time interval between speed updates in seconds.
        driver_speeds (dict): Dictionary containing driver-specific speeds.
    """
    current_time = datetime.now()

    while True:
        print(f"Current time: {current_time}")
        for segment, _ in segment_speeds.items():
            speed_update = get_speed(segment, driver_speeds)
            if speed_update is not None:
                print(
                    f"Segment: {segment}, Driver: {speed_update['driver_id']}, Speed: {speed_update['speed_mph']} mph"
                )
            else:
                print(f"Speed information not available for segment: {segment}")

        print("---------------------------")
        time.sleep(time_interval)
        current_time += timedelta(seconds=time_interval)


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

    driver_speeds = {
        "driver1": {"speed_mph": 45, "speed_variation": 5},
        "driver2": {"speed_mph": 50, "speed_variation": 3},
        "driver3": {"speed_mph": 55, "speed_variation": 4},
    }

    simulate_driving(time_interval, driver_speeds)
