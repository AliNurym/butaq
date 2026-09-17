#include <iostream>
#include <chrono>

long long collatz_len(long long n) {
    long long steps = 0;
    while (n > 1) {
        if (n % 2 == 0) {
            n /= 2;
        } else {
            n = 3 * n + 1;
        }
        steps++;
    }
    return steps;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long max_steps = 0;
    for (long long i = 1; i <= 100000; ++i) {
        long long s = collatz_len(i);
        if (s > max_steps) max_steps = s;
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << max_steps << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
