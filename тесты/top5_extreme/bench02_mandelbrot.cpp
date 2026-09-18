#include <iostream>
#include <chrono>

long long mandelbrot(long long w, long long h, long long max_iter) {
    long long total = 0;
    long long y = 0;
    while (y < h) {
        double cy = -1.5 + (3.0 * y) / (double)h;
        long long x = 0;
        while (x < w) {
            double cx = -2.0 + (3.0 * x) / (double)w;
            double zx = 0.0;
            double zy = 0.0;
            long long i = 0;
            while (i < max_iter) {
                double zx2 = zx * zx;
                double zy2 = zy * zy;
                if (zx2 + zy2 > 4.0) {
                    break;
                }
                zy = 2.0 * zx * zy + cy;
                zx = zx2 - zy2 + cx;
                i++;
            }
            total += i;
            x++;
        }
        y++;
    }
    return total;
}

int main() {
    auto t0 = std::chrono::high_resolution_clock::now();
    long long ans = mandelbrot(500, 500, 500);
    auto t1 = std::chrono::high_resolution_clock::now();
    double ms = std::chrono::duration<double, std::milli>(t1 - t0).count();
    std::cout << "RESULT:" << ans << std::endl;
    std::cout << "TIME_MS:" << ms << std::endl;
    return 0;
}
