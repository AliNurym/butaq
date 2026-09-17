#include <iostream>
#include <chrono>

double golden_ratio_reciprocal(long long n) {
    double x = 1.0;
    for (long long i = 0; i < n; ++i) {
        x = 1.0 + 1.0 / x;
    }
    return x;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    double phi = golden_ratio_reciprocal(10000000);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << phi << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
