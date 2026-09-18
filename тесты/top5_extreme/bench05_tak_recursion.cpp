#include <iostream>
#include <chrono>

long long tak(long long x, long long y, long long z) {
    if (y < x) {
        return tak(tak(x - 1, y, z), tak(y - 1, z, x), tak(z - 1, x, y));
    }
    return z;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long ans = tak(33, 16, 8);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
