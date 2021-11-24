#include <hello/hello.hpp>
#include <print/print.hpp>

int main() {
    hello::hello();
    print::print(hello::getHello());
    return 0;
}
