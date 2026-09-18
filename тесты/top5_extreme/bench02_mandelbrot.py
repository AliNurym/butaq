import time

def mandelbrot(w, h, max_iter):
    total = 0
    y = 0
    while y < h:
        cy = -1.5 + (3.0 * y) / h
        x = 0
        while x < w:
            cx = -2.0 + (3.0 * x) / w
            zx = 0.0
            zy = 0.0
            i = 0
            while i < max_iter:
                zx2 = zx * zx
                zy2 = zy * zy
                if zx2 + zy2 > 4.0:
                    break
                zy = 2.0 * zx * zy + cy
                zx = zx2 - zy2 + cx
                i += 1
            total += i
            x += 1
        y += 1
    return total

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = mandelbrot(500, 500, 500)
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{ans}")
    print(f"TIME_MS:{ms:.2f}")
