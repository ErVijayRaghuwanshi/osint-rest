import requests
import time
import statistics
from concurrent.futures import ThreadPoolExecutor, as_completed

URL = "http://osint.local/api/x/userinfo"
TIMEOUT = 5
CONCURRENCY = 6

usernames = [
    "UrbanPulse",
    "DailyOrbit",
    "EchoVerse",
    "MindSpark",
    "NextWaveHQ",
    "TrueVibesOnly",
    "PixelNomad",
    "BoldNarrative",
    "QuietlyLoud",
    "SwiftThoughts",
    "NovaSignals",
    "TheIdeaLab",
    "FreshAngle",
    "CosmicNotes",
    "ThinkInPublic",
    "ModernMotive",
    "VibeArchitect",
    "SignalAndNoise",
    "BrightContext",
    "TheDailyFrame"
]


def test_username(username):
    params = {"screen_name": username}
    start = time.perf_counter()
    try:
        response = requests.get(URL, params=params, timeout=TIMEOUT)
        elapsed = time.perf_counter() - start
        return {
            "username": username,
            "status": response.status_code,
            "time": elapsed,
            "success": response.status_code == 200
        }
    except Exception as e:
        elapsed = time.perf_counter() - start
        return {
            "username": username,
            "status": None,
            "time": elapsed,
            "success": False,
            "error": str(e)
        }


results = []
start_test = time.perf_counter()

with ThreadPoolExecutor(max_workers=CONCURRENCY) as executor:
    futures = [executor.submit(test_username, u) for u in usernames]
    for future in as_completed(futures):
        results.append(future.result())

total_time = time.perf_counter() - start_test

# Metrics
times = [r["time"] for r in results]
successes = sum(1 for r in results if r["success"])

print("\n=== Endpoint Performance Report ===")
print(f"Endpoint: {URL}")
print(f"Usernames tested: {len(usernames)}")
print(f"Concurrency: {CONCURRENCY}")
print(f"Total test time: {total_time:.2f}s")
print(f"Successful responses: {successes}")
print(f"Failed responses: {len(usernames) - successes}")
print(f"Average latency: {statistics.mean(times):.4f}s")
print(f"Min latency: {min(times):.4f}s")
print(f"Max latency: {max(times):.4f}s")

print("\nPer-username results:")
for r in sorted(results, key=lambda x: x["time"], reverse=True):
    print(f"{r['username']:15} | {r['status']} | {r['time']:.4f}s")
