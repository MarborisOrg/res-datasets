#include <iostream>
#include <filesystem>
#include <string>
#include <cstdlib>
#include <stdexcept>
#include <sstream>
#include <unistd.h>  // برای getcwd
#include <pwd.h>     // برای دسترسی به دایرکتوری خانگی کاربر

#include "util/df.cxx"

namespace fs = std::filesystem;

void copyDir(const fs::path& src, const fs::path& dst) {
    // کپی کردن دایرکتوری به دایرکتوری دیگر
    try {
        fs::copy(src, dst, fs::copy_options::recursive | fs::copy_options::overwrite_existing);
    } catch (const fs::filesystem_error& e) {
        throw std::runtime_error("Error copying directory: " + std::string(e.what()));
    }
}

int main() {
    try {
        // دریافت دایرکتوری خانگی کاربر
        const char* homeDir = getenv("HOME");
        if (!homeDir) {
            throw std::runtime_error("Error getting user home directory");
        }
        
        fs::path targetDir = fs::path(homeDir) / ".marboris";

        // حذف دایرکتوری قدیمی در صورت وجود
        if (fs::exists(targetDir)) {
            fs::remove_all(targetDir);
        }

        // ایجاد دایرکتوری جدید
        fs::create_directory(targetDir);

        // دریافت دایرکتوری کاری فعلی
        char cwd[1024];
        if (getcwd(cwd, sizeof(cwd)) != nullptr) {
            fs::path resDir = fs::path(cwd) / "res";

            // کپی کردن دایرکتوری
            copyDir(resDir, targetDir);
            
            // اجرای تابع RunChecker از ماژول util
            RunChecker(); // فرض بر این است که این تابع در همان فایل یا یک فایل دیگر تعریف شده است
            
            std::cout << "Datasets saved successfully! We are ready to go" << std::endl;
        } else {
            throw std::runtime_error("Error getting current directory");
        }
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}
