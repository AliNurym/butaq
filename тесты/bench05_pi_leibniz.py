import time

def compute_pi(n):
    total = 0.0
    sign = 1.0
    k = 0
    while k < n:
        total += sign / (2.0 * k + 1.0)
        sign = -sign
        k += 1
    return total * 4.0

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = compute_pi(5000000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans:.6f}")
    print(f"TIME_MS:{ms:.2f}")
