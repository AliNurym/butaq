import time

def bitwise_benchmark(n):
    state = 123456789
    acc = 0
    for i in range(n):
        # 32-bit/64-bit PRNG mixing
        state = (state * 1664525 + 1013904223) & 0xFFFFFFFF
        acc = (acc + (state ^ (state >> 11))) & 0xFFFFFFFF
    return acc

start = time.perf_counter()
res = bitwise_benchmark(20_000_000)
elapsed = (time.perf_counter() - start) * 1000.0

print(f"RESULT: {res}")
print(f"PYTHON_TIME: {elapsed:.2f} ms")
