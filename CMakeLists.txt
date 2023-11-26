cmake_minimum_required(VERSION 3.15)
project(GoRemoteFFI)

set(TARGET remote_server_go)
file(GLOB SRC ./*.go)
set(GOPATH ${CMAKE_CURRENT_BINARY_DIR})
set(LIB remote.a)
find_program(GOLANG go REQUIRED)

add_custom_command(OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/${LIB}
  DEPENDS ${SRC}
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  COMMAND env GOPATH=${GOPATH} ${GOLANG} build -buildmode=c-archive -o "${CMAKE_CURRENT_BINARY_DIR}/${LIB}" .
  COMMENT "Building ira/goremote library..."
)
add_custom_target(${TARGET} DEPENDS ${LIB} ${HEADER})
add_library(goremote STATIC IMPORTED GLOBAL)
add_dependencies(goremote ${TARGET})
set_target_properties(goremote
  PROPERTIES
  IMPORTED_LOCATION ${CMAKE_CURRENT_BINARY_DIR}/${LIB}
  INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR})
