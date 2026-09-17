import time

def mul(A, B):
    return [
        A[0] * B[0] + A[1] * B[2],
        A[0] * B[1] + A[1] * B[3],
        A[2] * B[0] + A[3] * B[2],
        A[2] * B[1] + A[3] * B[3]
    ]

def power(M, p):
    res = [1.0, 0.0, 0.0, 1.0]
    base = [M[0], M[1], M[2], M[3]]
    while p > 0:
        if p % 2 == 1:
            res = mul(res, base)
        base = mul(base, base)
        p //= 2
    return res

if __name__ == "__main__":
    t0 = time.perf_counter()
    ans = 0
    T = [1.0, 1.0, 1.0, 0.0]
    k = 0
    while k < 500000:
        res = power(T, 30)
        ans = res[0]
        k += 1
    t1 = time.perf_counter()
    ms = (t1 - t0) * 1000.0
    print(f"RESULT:{int(ans)}")
    print(f"TIME_MS:{ms:.2f}")
