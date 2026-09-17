import time

def golden_ratio_reciprocal(n):
    x = 1.0
    i = 0
    while i < n:
        x = 1.0 + 1.0 / x
        i += 1
    return x

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = golden_ratio_reciprocal(10000000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans:.8f}")
    print(f"TIME_MS:{ms:.2f}")
