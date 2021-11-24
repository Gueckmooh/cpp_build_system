#include <hello/hello.hpp>
#include <print/print.hpp>

#include <string>

namespace hello {

void hello() {
    print::print("Hello");
}

std::string getHello() {
    return "Hello";
}

}
