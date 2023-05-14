import random
import time

# Sample data for street segment speeds.
segment_speeds = {
    "Segment1": 25,
    "Segment2": 20,
    "Segment3": 30,
    "Segment4": 15,
}


def get_speed(segment: str) -> int:
    """
    Simulates the speed for a given street segment.

    Args:
        segment (str): The name of the street segment.

    Returns:
        int or None: The simulated speed in mph, or None if speed information is not available for the segment.
    """
    if segment in segment_speeds:
        base_speed = segment_speeds[segment]
        # Simulate some random variation in speed.
        speed_variation = random.randint(-5, 5)
        speed = base_speed + speed_variation
        return speed
    else:
        return None


def send_speed_updates() -> None:
    """
    Simulates sending real-time speed updates for street segments.
    Continuously generates speed updates and prints the current speed for each segment.
    """
    while True:
        for segment in segment_speeds:
            speed = get_speed(segment)
            if speed is not None:
                print(f"The current speed on {segment} is {speed} mph.")
            else:
                print("Speed information not available for the given segment.")
        time.sleep(5)  # Update speed every 5 seconds.


send_speed_updates()
