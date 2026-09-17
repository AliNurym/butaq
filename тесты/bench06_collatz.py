import time

def collatz_len(n):
    steps = 0
    while n > 1:
        if n % 2 == 0:
            n //= 2
        else:
            n = 3 * n + 1
        steps += 1
    return steps

if __name__ == "__main__":
    t0 = time.perf_counter()
    max_steps = 0
    i = 1
    while i <= 100000:
        s = collatz_len(i)
        if s > max_steps:
            max_steps = s
        i += 1
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{max_steps}")
    print(f"TIME_MS:{ms:.2f}")
