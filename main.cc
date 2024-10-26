#include <iostream>
#include <filesystem>
#include <string>
#include <cstdlib>
#include <stdexcept>
#include <sstream>
#include <unistd.h>
#include <pwd.h>

#include "df.cc"

namespace fs = std::filesystem;

void copyDir(const fs::path &src, const fs::path &dst)
{
    try
    {
        fs::copy(src, dst, fs::copy_options::recursive | fs::copy_options::overwrite_existing);
    }
    catch (const fs::filesystem_error &e)
    {
        throw std::runtime_error("Error copying directory: " + std::string(e.what()));
    }
}

int main()
{
    try
    {
        const char *homeDir = getenv("HOME");
        if (!homeDir)
        {
            throw std::runtime_error("Error getting user home directory");
        }

        fs::path targetDir = fs::path(homeDir) / ".marboris";

        if (fs::exists(targetDir))
        {
            fs::remove_all(targetDir);
        }

        fs::create_directory(targetDir);

        char cwd[1024];
        if (getcwd(cwd, sizeof(cwd)) != nullptr)
        {
            fs::path resDir = fs::path(cwd) / "res";

            copyDir(resDir, targetDir);

            RunChecker();

            std::cout << "Datasets saved successfully! We are ready to go" << std::endl;
        }
        else
        {
            throw std::runtime_error("Error getting current directory");
        }
    }
    catch (const std::exception &e)
    {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }

    return 0;
}
