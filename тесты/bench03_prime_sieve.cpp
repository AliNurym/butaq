#include <iostream>
#include <vector>
#include <chrono>

int count_primes(int limit) {
    std::vector<bool> is_prime(limit + 1, true);
    is_prime[0] = is_prime[1] = false;
    for (int p = 2; p * p <= limit; p++) {
        if (is_prime[p]) {
            for (int i = p * p; i <= limit; i += p) {
                is_prime[i] = false;
            }
        }
    }
    int count = 0;
    for (int p = 2; p <= limit; p++) {
        if (is_prime[p]) count++;
    }
    return count;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    int ans = count_primes(200000);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
