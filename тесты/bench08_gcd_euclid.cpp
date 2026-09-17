#include <iostream>
#include <chrono>

long long gcd(long long a, long long b) {
    while (b != 0) {
        long long t = b;
        b = a % b;
        a = t;
    }
    return a;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long total = 0;
    for (long long i = 1; i <= 1000000; ++i) {
        total += gcd(i, i + 17);
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << total << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
