#include <iostream>
#include <chrono>

struct Mat2 {
    double a, b, c, d;
};

Mat2 mul(const Mat2& X, const Mat2& Y) {
    return {
        X.a * Y.a + X.b * Y.c,
        X.a * Y.b + X.b * Y.d,
        X.c * Y.a + X.d * Y.c,
        X.c * Y.b + X.d * Y.d
    };
}

Mat2 power(Mat2 M, int p) {
    Mat2 res = {1.0, 0.0, 0.0, 1.0};
    Mat2 base = M;
    while (p > 0) {
        if (p % 2 == 1) res = mul(res, base);
        base = mul(base, base);
        p /= 2;
    }
    return res;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    double ans = 0;
    Mat2 T = {1.0, 1.0, 1.0, 0.0};
    for (int k = 0; k < 500000; k++) {
        Mat2 res = power(T, 30);
        ans = res.a;
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << (long long)ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
