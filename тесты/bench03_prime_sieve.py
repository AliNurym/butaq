import time

def count_primes(limit):
    is_prime = [True] * (limit + 1)
    is_prime[0] = is_prime[1] = False
    p = 2
    while p * p <= limit:
        if is_prime[p]:
            i = p * p
            while i <= limit:
                is_prime[i] = False
                i += p
        p += 1
    count = 0
    p = 2
    while p <= limit:
        if is_prime[p]:
            count += 1
        p += 1
    return count

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = count_primes(200000)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
