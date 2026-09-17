import time
import sys

def fib(n):
    if n <= 1:
        return n
    return fib(n - 1) + fib(n - 2)

if __name__ == "__main__":
    n = 25
    if len(sys.argv) > 1:
        n = int(sys.argv[1])
    
    t0 = time.perf_counter()
    ans = fib(n)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
