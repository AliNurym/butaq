#include <iostream>
#include <iomanip>
#include <chrono>

double calc_pi(long long n) {
    double total = 0.0;
    double sign = 1.0;
    long long k = 0;
    while (k < n) {
        total += sign / (2.0 * (double)k + 1.0);
        sign = -sign;
        k++;
    }
    return total * 4.0;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    double ans = calc_pi(50000000);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << std::fixed << std::setprecision(5);
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
