#include <gtest/gtest.h>
#include <my_math/math.hpp>

TEST(tPow, pow1) {
    int v = my_pow(3, 2);
    ASSERT_EQ(v, 9);
}

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
