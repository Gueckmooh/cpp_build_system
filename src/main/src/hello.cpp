#include <hello/hello.hpp>
#include <utils/print/print.hpp>

int main() {
    hello::hello();
    print::print(hello::getHello());
    return 0;
}
