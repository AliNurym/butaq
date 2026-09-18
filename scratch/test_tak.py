import time

def tak(x, y, z):
    if y < x:
        return tak(tak(x - 1, y, z), tak(y - 1, z, x), tak(z - 1, x, y))
    return z

start = time.perf_counter()
res = tak(28, 16, 8)
elapsed = (time.perf_counter() - start) * 1000.0

print(f"RESULT: {res}")
print(f"PYTHON_TIME: {elapsed:.2f} ms")
