#include <gtest/gtest.h>

TEST(Sub, Sub1) {
    ASSERT_EQ(46 - 4, 42);
}

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
