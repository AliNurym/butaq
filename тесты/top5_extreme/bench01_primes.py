import time

def count_primes(n):
    count = 0
    i = 2
    while i <= n:
        is_prime = 1
        d = 2
        while d * d <= i:
            q = i // d
            if i - q * d == 0:
                is_prime = 0
                break
            d += 1
        if is_prime == 1:
            count += 1
        i += 1
    return count

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = count_primes(1000000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
