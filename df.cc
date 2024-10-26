#include <iostream>
#include <fstream>
#include <vector>
#include <string>
#include <map>
#include <nlohmann/json.hpp>
#include <ctime>
#include <iomanip>

using json = nlohmann::json;

struct Dataset
{
    std::string tag;
    std::vector<std::string> patterns;
    std::vector<std::string> responses;
    std::string context;
};

std::vector<Dataset> loadDataset(const std::string &filename)
{
    std::ifstream file(filename);
    if (!file.is_open())
    {
        throw std::runtime_error("Error opening file: " + filename);
    }

    json j;
    file >> j;

    std::vector<Dataset> dataset;
    for (const auto &item : j)
    {
        Dataset data;
        data.tag = item["tag"];
        data.patterns = item["patterns"].get<std::vector<std::string>>();
        data.responses = item["responses"].get<std::vector<std::string>>();
        data.context = item["context"];
        dataset.push_back(data);
    }

    return dataset;
}

std::map<std::string, std::vector<std::string>> findDuplicates(const std::vector<Dataset> &dataset)
{
    std::map<std::string, std::vector<std::string>> duplicatePatterns;

    for (const auto &data : dataset)
    {
        std::map<std::string, bool> patternMap;
        for (const auto &pattern : data.patterns)
        {
            if (patternMap[pattern])
            {
                duplicatePatterns[data.tag].push_back(pattern);
            }
            else
            {
                patternMap[pattern] = true;
            }
        }
    }

    return duplicatePatterns;
}

void logDuplicates(const std::map<std::string, std::vector<std::string>> &duplicates, const std::string &logFile)
{
    std::ofstream file(logFile, std::ios::app);
    if (!file.is_open())
    {
        throw std::runtime_error("Error opening log file: " + logFile);
    }

    for (const auto &[tag, patterns] : duplicates)
    {
        if (!patterns.empty())
        {
            file << "Tag: " << tag << ", Duplicates: ";
            for (const auto &pattern : patterns)
            {
                file << pattern << " ";
            }
            file << std::endl;
        }
    }
}

void RunChecker()
{
    try
    {
        auto dataset = loadDataset("./res/locales/en/intents.json");
        auto duplicates = findDuplicates(dataset);

        for (const auto &[tag, patterns] : duplicates)
        {
            if (!patterns.empty())
            {
                std::cout << "Tag: " << tag << ", Duplicates: ";
                for (const auto &pattern : patterns)
                {
                    std::cout << pattern << " ";
                }
                std::cout << std::endl;
            }
        }

        auto t = std::time(nullptr);
        auto tm = *std::localtime(&t);
        std::ostringstream logFileName;
        logFileName << "./log/duplicate_report_"
                    << std::put_time(&tm, "%Y-%m-%d_%H-%M-%S") << ".log";

        if (!duplicates.empty())
        {
            logDuplicates(duplicates, logFileName.str());
            std::cout << "Duplicate patterns logged to '" << logFileName.str() << "'" << std::endl;
        }
        else
        {
            std::cout << "No duplicates found, no log file created." << std::endl;
        }
    }
    catch (const std::exception &e)
    {
        std::cerr << "Error: " << e.what() << std::endl;
    }
}

int test()
{
    RunChecker();
    return 0;
}
