import time

def gcd(a, b):
    while b != 0:
        a, b = b, a % b
    return a

if __name__ == "__main__":
    t0 = time.perf_counter()
    total = 0
    i = 1
    while i <= 1000000:
        total += gcd(i, i + 17)
        i += 1
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{total}")
    print(f"TIME_MS:{ms:.2f}")
