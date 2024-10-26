cmake_minimum_required(VERSION 3.10)
project(Marboris)

find_package(nlohmann_json 3.11.3 REQUIRED)

add_executable(main main.cc)
target_link_libraries(main nlohmann_json::nlohmann_json)
