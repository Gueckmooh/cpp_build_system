#include <gtest/gtest.h>
#include <testable/testable.hpp>

TEST(Add, Add1) {
    ASSERT_EQ(testable::add(40, 2), 42);
}

int main(int argc, char** argv) {
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}
