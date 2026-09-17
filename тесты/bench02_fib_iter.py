import time

def fib_iter(n):
    if n <= 1:
        return n
    a = 0
    b = 1
    i = 2
    while i <= n:
        c = a + b
        a = b
        b = c
        i += 1
    return b

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = 0
    k = 0
    while k < 500000:
        ans = fib_iter(80)
        k += 1
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
