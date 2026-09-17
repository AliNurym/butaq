import time

def run_reduction():
    N = 100000
    arr = [0.0] * N
    i = 0
    while i < N:
        arr[i] = float((i * 3 + 7) % 1000)
        i += 1

    total = 0.0
    p = 0
    while p < 50:
        i = 0
        while i < N:
            total += arr[i]
            i += 1
        p += 1
    return total

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = run_reduction()
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{int(ans)}")
    print(f"TIME_MS:{ms:.2f}")
