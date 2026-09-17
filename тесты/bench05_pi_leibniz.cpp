#include <iostream>
#include <chrono>

double compute_pi(long long n) {
    double sum = 0.0;
    double sign = 1.0;
    for (long long k = 0; k < n; ++k) {
        sum += sign / (2.0 * k + 1.0);
        sign = -sign;
    }
    return sum * 4.0;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    double pi = compute_pi(5000000);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << pi << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
