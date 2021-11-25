#include <cmath>
#include <thread>
#include <my_math/math.hpp>
#include <iostream>

void popo() {
    std::cout << "popo" << std::endl;
}

int my_pow(int base, int exp) {
    std::thread t(popo);
    t.join();
    return (int)floorf(pow(base, exp));
}
