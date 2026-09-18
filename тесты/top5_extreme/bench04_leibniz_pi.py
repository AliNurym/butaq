import time

def calc_pi(n):
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
    ans = calc_pi(50000000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans:.5f}")
    print(f"TIME_MS:{ms:.2f}")
