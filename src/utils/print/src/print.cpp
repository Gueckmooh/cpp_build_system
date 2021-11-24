#include <print/print.hpp>

#include <string_view>
#include <iostream>

namespace print {
void print(const std::string_view& s) {
    std::cout << s << std::endl;
}
}
