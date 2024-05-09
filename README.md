Netcat - which is usually abbreviated to nc is a command line networking utility for reading and writing to network connections with TCP or UDP.

This repo is attempt to solve John Crickett's (https://codingchallenges.fyi/challenges/intro/?utm_source=substack&utm_medium=email) Coding Challenges in Golang

Steps to run - 

TCP Server - 
  1. Navtigate to root folder
  2. Run `go main . tcp` - it starts a tcp server on port 8080
  3. Connect to server with nc - `nc localhost 8080`
  4. One way communication is established from client to server

Turn process into server - 
  1. Navigate to root folder
  2. Run `go main . -e /bin/bash` - it starts a bash process and receives input and sends output to the client
