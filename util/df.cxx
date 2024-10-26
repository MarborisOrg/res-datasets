#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <filesystem>
#include <stdexcept>
#include <sstream>
#include <pwd.h>
#include <unistd.h>  // برای getcwd

namespace fs = std::filesystem;

struct Dataset {
    std::string tag;
    std::vector<std::string> patterns;
    std::vector<std::string> responses;
    std::string context;
};

std::vector<Dataset> loadDataset(const std::string& filename) {
    std::ifstream file(filename);
    if (!file.is_open()) {
        throw std::runtime_error("Error opening file: " + filename);
    }

    std::vector<Dataset> dataset;
    std::string line;

    // فرض بر این است که هر خط یک شیء JSON است
    while (std::getline(file, line)) {
        if (line.find('{') != std::string::npos) {
            Dataset data;
            // به صورت دستی داده‌ها را پر می‌کنیم
            data.tag = "example_tag";  // مقداردهی نمونه
            data.patterns.push_back("example_pattern");  // اضافه کردن الگو
            data.context = "example_context";  // مقداردهی نمونه
            dataset.push_back(data);
        }
    }

    return dataset;
}

void RunChecker() {
    try {
        auto dataset = loadDataset("./res/locales/en/intents.json");
        // ادامه پردازش...
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
    }
}
