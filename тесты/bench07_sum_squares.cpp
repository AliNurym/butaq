#include <iostream>
#include <chrono>

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long sum = 0;
    long long mod = 1000000007;
    for (long long i = 1; i <= 10000000; ++i) {
        sum = (sum + i * i) % mod;
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << sum << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
