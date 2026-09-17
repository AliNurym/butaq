#include <iostream>
#include <vector>
#include <chrono>

int main() {
    int N = 100000;
    std::vector<double> arr(N);
    for (int i = 0; i < N; ++i) {
        arr[i] = (i * 3 + 7) % 1000;
    }

    auto t0 = std::chrono::high_resolution_clock::now();
    double total = 0.0;
    for (int pass = 0; pass < 50; ++pass) {
        for (int i = 0; i < N; ++i) {
            total += arr[i];
        }
    }
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << (long long)total << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
