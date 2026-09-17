#include <iostream>
#include <chrono>
#include <cstdlib>

long long fib(long long n) {
    if (n <= 1) return n;
    return fib(n - 1) + fib(n - 2);
}

int main(int argc, char* argv[]) {
    long long n = 25;
    if (argc > 1) {
        n = std::atoll(argv[1]);
    }
    auto t0 = std::chrono::high_resolution_clock::now();
    long long ans = fib(n);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
