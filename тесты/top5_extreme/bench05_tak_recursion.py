import time

def tak(x, y, z):
    if y < x:
        return tak(tak(x - 1, y, z), tak(y - 1, z, x), tak(z - 1, x, y))
    return z

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = tak(33, 16, 8)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
