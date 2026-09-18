#include <iostream>
#include <chrono>

long long count_primes(long long n) {
    long long count = 0;
    long long i = 2;
    while (i <= n) {
        int is_prime = 1;
        long long d = 2;
        while (d * d <= i) {
            long long q = i / d;
            if (i - q * d == 0) {
                is_prime = 0;
                break;
            }
            d++;
        }
        if (is_prime == 1) {
            count++;
        }
        i++;
    }
    return count;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long ans = count_primes(1000000);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
