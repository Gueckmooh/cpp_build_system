#include <hello/hello.hpp>
#include <utils/print/print.hpp>

#include <string>

namespace hello {

void hello() {
    print::print("Hello");
}

std::string getHello() {
    return "Hello";
}

}
