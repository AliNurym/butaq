import time

def compute_sum_squares(n):
    total = 0
    mod = 1000000007
    i = 1
    while i <= n:
        total = (total + i * i) % mod
        i += 1
    return total

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = compute_sum_squares(10000000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
