#include <iostream>
#include <chrono>
#include <cstdlib>

long long fib_iter(long long n) {
    if (n <= 1) return n;
    long long a = 0;
    long long b = 1;
    long long i = 2;
    while (i <= n) {
        long long c = a + b;
        a = b;
        b = c;
        i++;
    }
    return b;
}

int main(int argc, char* argv[]) {
    long long n = 70;
    long long iterations = 100000;
    if (argc > 1) {
        iterations = std::atoll(argv[1]);
    }

    auto t0 = std::chrono::high_resolution_clock::now();
    long long ans = 0;
    for (long long k = 0; k < iterations; ++k) {
        ans = fib_iter(n);
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
